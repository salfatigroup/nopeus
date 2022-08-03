package core

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
	"time"

	helmclient "github.com/mittwald/go-helm-client"
	"github.com/salfatigroup/nopeus/config"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/tools/clientcmd/api"
	ctrl "sigs.k8s.io/controller-runtime"
)

// creates a helm client based on the kubeconfig and required namespace
func newHelmClient(namespace string, context string) (helmclient.Client, error) {
    // get kubeconfig
    kubeconfig, err := getKubeconfigAsBytes()
    if err != nil {
        return nil, err
    }

    opt := &helmclient.KubeConfClientOptions{
        Options: &helmclient.Options{
            Namespace:        namespace, // Change this to the namespace you wish to install the chart in.
            RepositoryCache:  "/tmp/.helmcache",
            RepositoryConfig: "/tmp/.helmrepo",
            Debug:            true,
            Linting:          true, // Change this to false if you don't want linting.
            DebugLog: func(format string, v ...interface{}) {
                // Change this to your own logger. Default is 'log.Printf(format, v...)'.
                fmt.Printf(format, v...)
                fmt.Printf("\n")
            },
        },
        KubeContext: context,
        KubeConfig: kubeconfig,
    }

    // initialize helm client with the correct kube context
    helmClient, err := helmclient.NewClientFromKubeConf(opt, helmclient.Burst(100), helmclient.Timeout(time.Duration(time.Minute*15)))
    if err != nil {
        return nil, err
    }

    return helmClient, nil
}

// return a new kubernetes client
func newKubernetesClient() (*kubernetes.Clientset, error) {
    config, err := ctrl.GetConfig()
    if err != nil {
        return nil, err
    }

    return kubernetes.NewForConfig(config)
}

// run and deploy k8s and helm files per environment
func runK8s(cfg *config.NopeusConfig) error {
    // apply the helm charts for each environment
    for _, env := range cfg.Runtime.Environments {
        var kubeContext string

        if cfg.Runtime.DryRun {
            kubeContext = "dryrun"
        } else {
            kubeContext, err := connectToCluster(cfg, env)
            if err != nil {
                return err
            }

            // create private registry secrets from dockerconfig
            if err = createPrivateRegistrySecrets(cfg, env, kubeContext); err != nil {
                return err
            }
        }

        // run manual helm commands before the generic setup
        // this should be avoid as much as possible
        if err := manualHelmCommands(cfg, env, kubeContext); err != nil {
            return err
        }

        if err := applyK8sHelmCharts(cfg, env, kubeContext); err != nil {
            return err
        }
    }

    return nil
}

// connect to private registries via the .dockerconfig file
// this function assumes the user executed `docker login` beforehand
func createPrivateRegistrySecrets(cfg *config.NopeusConfig, env string, kubeContext string) error {
    // get $NOPEUS_DOCKER_SERVER, $NOPEUS_DOCKER_USERNAME, $NOPEUS_DOCKER_PASSWORD, $NOPEUS_DOCKER_EMAIL from env
    dockerServer := os.Getenv("NOPEUS_DOCKER_SERVER")
    dockerUsername := os.Getenv("NOPEUS_DOCKER_USERNAME")
    dockerPassword := os.Getenv("NOPEUS_DOCKER_PASSWORD")
    dockerEmail := os.Getenv("NOPEUS_DOCKER_EMAIL")

    // check if should login to private registry
    if dockerServer == "" || dockerUsername == "" || dockerPassword == "" || dockerEmail == "" {
        return nil
    }

    fmt.Println("Logging into private registry...")

    // create namespace
    kubeClient, err := newKubernetesClient()
    if err != nil {
        return err
    }

    // create the Runtime.DefaultNamespace if it doesn't exist
    if _, err := kubeClient.CoreV1().Namespaces().Get(context.Background(), cfg.Runtime.DefaultNamespace, metav1.GetOptions{}); err != nil {
        _, err := kubeClient.CoreV1().Namespaces().Create(context.Background(), &v1.Namespace{
            ObjectMeta: metav1.ObjectMeta{
                Name: cfg.Runtime.DefaultNamespace,
            },
        }, metav1.CreateOptions{})

        if err != nil {
            return err
        }

        fmt.Printf("Created namespace %s\n", cfg.Runtime.DefaultNamespace)
    }

    // create secret from the following environment variables $NOPEUS_DOCKER_SERVER, $NOPEUS_DOCKER_USERNAME, $NOPEUS_DOCKER_PASSWORD
    if _, err := kubeClient.CoreV1().Secrets(cfg.Runtime.DefaultNamespace).Get(context.Background(), "dockerconfig", metav1.GetOptions{}); err != nil {
        _, err := kubeClient.CoreV1().Secrets(cfg.Runtime.DefaultNamespace).Create(context.Background(), &v1.Secret{
            ObjectMeta: metav1.ObjectMeta{
                Name: "dockerconfig",
            },
            Type: "kubernetes.io/dockerconfigjson",
            Data: map[string][]byte{
                // auth is dockerUsername:dockerPassword encoded in base64
                ".dockerconfigjson": []byte(fmt.Sprintf(`{"auths":{"%s":{"username":"%s","password":"%s","email":"%s","auth":"%s"}}}`, dockerServer, dockerUsername, dockerPassword, dockerEmail, base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s:%s", dockerUsername, dockerPassword))))),
            },
        }, metav1.CreateOptions{})

        if err != nil {
            return err
        }

        fmt.Println("Successfully created secret: dockerconfig")
    }

    return nil
}

func getKubeconfigAsBytes() ([]byte, error) {
    // get the kubeconfig
    kubeconfigPath, err := findKubeConfig()
    if err != nil {
        return nil, err
    }

    return ioutil.ReadFile(kubeconfigPath)
}

func manualHelmCommands(cfg *config.NopeusConfig, env string, kubeContext string) error {
    // install cert-manager manually
    fmt.Println("Installing cert-manager...")

    // get client pointing to cert-manager namespace
    helmClient, err := newHelmClient("cert-manager", kubeContext)
    if err != nil {
        return err
    }

    // install cert-manager
    ctx, cancel := context.WithTimeout(context.Background(), time.Duration(time.Minute*15))
    defer cancel()

    _, err = helmClient.InstallOrUpgradeChart(ctx, &helmclient.ChartSpec{
        ReleaseName: "cert-manager",
        ChartName: "jetstack/cert-manager",
        Namespace: "cert-manager",
        CreateNamespace: true,
        ValuesYaml: "installCRDs: true",
        DryRun: cfg.Runtime.DryRun,
        Wait: true,
        Timeout: time.Duration(time.Minute*15),
    }, nil)

    if err != nil {
        return err
    }

    return nil
}

func applyK8sHelmCharts(cfg *config.NopeusConfig, env string, kubeContext string) (err error) {
    // apply the helm charts for the environment
    for _, service := range cfg.Runtime.HelmRuntime.ServiceTemplateData {
        if err = applyHelmChart(cfg, service, kubeContext); err != nil {
            return
        }
    }

    return nil
}

// apply a helm chart for the environment
func applyHelmChart(cfg *config.NopeusConfig, service config.ServiceTemplateData, kubeContext string) error {
    fmt.Printf("Applying helm chart for service %s\n", service.GetName())
    // get chart specifications
    chartSpec, err := service.GetChartSpec()
    if err != nil {
        return err
    }

    // get the helm client
    // get client pointing to cert-manager namespace
    helmClient, err := newHelmClient(chartSpec.Namespace, kubeContext)
    if err != nil {
        return err
    }

    // install the chart
    ctx, cancel := context.WithTimeout(context.Background(), time.Duration(time.Minute*15))
    defer cancel()
    if _, err := helmClient.InstallOrUpgradeChart(ctx, chartSpec, nil); err != nil {
        return err
    }

    return nil
}

// connect to relevant k8s cluster
func connectToCluster(cfg *config.NopeusConfig, env string) (string, error) {
    switch cfg.CAL.CloudVendor {
    case "aws":
        // connect to aws
        return connectToEks(cfg, env)
    default:
        return "", fmt.Errorf("cloud vendor %s not supported at the moment", cfg.CAL.CloudVendor)
    }
}

func connectToEks(cfg *config.NopeusConfig, env string) (string, error) {
    // get terraform outputs
    tfOutputs := cfg.Runtime.Infrastructure[env].GetOutputs()

    // get the region and cluster name values from the terraform outputs
    var region string
    var clusterName string

    if err := json.Unmarshal(tfOutputs["region"].Value, &region); err != nil {
        return "", err
    }
    if err := json.Unmarshal(tfOutputs["name"].Value, &clusterName); err != nil {
        return "", err
    }

    cmdString := fmt.Sprintf("eks update-kubeconfig --region %s --name %s", region, clusterName)
    cmdArray := strings.Split(cmdString, " ")
    out, err := exec.Command("aws", cmdArray...).Output()
    fmt.Println(string(out))
    if err != nil {
        return "", err
    }

    return getActiveKubeContext()
}

// get the active kube context
func getActiveKubeContext() (string, error) {
    kubeconfig, err := loadKubeConfig()
    if err != nil {
        return "", err
    }

    return kubeconfig.CurrentContext, nil
}

// findKubeConfig finds path from env:KUBECONFIG or ~/.kube/config
func findKubeConfig() (string, error) {
    env := os.Getenv("KUBECONFIG")
    if env != "" {
        return env, nil
    }

    // get the kubeconfig
    home, err := os.UserHomeDir()
    if err != nil {
        return "", err
    }

    kubeconfigPath := home + "/.kube/config"
	return kubeconfigPath, nil
}

// load the kube config yaml
func loadKubeConfig() (*api.Config, error) {
    kubeConfigPath, err := findKubeConfig()
	if err != nil {
		return nil, err
	}

    kubeConfig, err := clientcmd.LoadFromFile(kubeConfigPath)
	if err != nil {
		return nil, err
	}

    return kubeConfig, nil
}

