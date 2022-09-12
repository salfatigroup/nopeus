package core

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/salfatigroup/nopeus/cli/util"
	"github.com/salfatigroup/nopeus/config"
	sgck "github.com/salfatigroup/nopeus/kubernetes"
	"github.com/salfatigroup/nopeus/logger"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	ctrl "sigs.k8s.io/controller-runtime"
)

// return a new kubernetes client
func newKubernetesClient() (*kubernetes.Clientset, error) {
	config, err := ctrl.GetConfig()
	if err != nil {
		return nil, err
	}

	return kubernetes.NewForConfig(config)
}

// run and deploy k8s and helm files per environment
func runK8s(envName string, envData *config.EnvironmentConfig, cfg *config.NopeusConfig) error {
	// apply the helm charts for each environment
	var kubeContext string

	if cfg.Runtime.DryRun {
		kubeContext = "dryrun"
	} else {
		kubeContext, err := connectToCluster(cfg, envName, envData)
		if err != nil {
			return err
		}

		// the the kubecontext onto the environment
		envData.SetKubeContext(kubeContext)

		// load the checksum map onto the envData
		if err := envData.LoadChecksumMap(); err != nil {
			return err
		}

		// create private registry secrets from dockerconfig
		if err = createPrivateRegistrySecrets(cfg, envName, envData, kubeContext); err != nil {
			return err
		}
	}

	if err := applyK8sHelmCharts(cfg, envName, envData, kubeContext); err != nil {
		return err
	}

	return nil
}

// connect to private registries via the .dockerconfig file
// this function assumes the user executed `docker login` beforehand
func createPrivateRegistrySecrets(cfg *config.NopeusConfig, envName string, envData *config.EnvironmentConfig, kubeContext string) error {
	// get $NOPEUS_DOCKER_SERVER, $NOPEUS_DOCKER_USERNAME, $NOPEUS_DOCKER_PASSWORD, $NOPEUS_DOCKER_EMAIL from env
	dockerServer := os.Getenv("NOPEUS_DOCKER_SERVER")
	dockerUsername := os.Getenv("NOPEUS_DOCKER_USERNAME")
	dockerPassword := os.Getenv("NOPEUS_DOCKER_PASSWORD")
	dockerEmail := os.Getenv("NOPEUS_DOCKER_EMAIL")

	// check if should login to private registry
	if dockerServer == "" || dockerUsername == "" || dockerPassword == "" || dockerEmail == "" {
		return nil
	}

	// fmt.Println("Logging into private registry...")

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

		// fmt.Printf("Created namespace %s\n", cfg.Runtime.DefaultNamespace)
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

		// fmt.Println("Successfully created secret: dockerconfig")
	}

	return nil
}

func applyK8sHelmCharts(cfg *config.NopeusConfig, envName string, envData *config.EnvironmentConfig, kubeContext string) error {
	// apply the helm charts for the environment
	for _, service := range cfg.Runtime.HelmRuntime.ServiceTemplateData {
		// compare service checksum to the checksum map
		// and skip if the same
		serviceChecksum, err := service.GetChecksum()
		if err != nil {
			return err
		}

		logger.Debugf("serviceChecksum: %s", serviceChecksum)
		logger.Debugf("service values: %+v", service)
		logger.Debugf("service helm values: %+v", service.GetHelmValues())

		logger.Debugf("envData checksum: %s", envData.GetChecksum(service.GetName()))
		if serviceChecksum == envData.GetChecksum(service.GetName()) {
			fmt.Println(util.GrayText("Skipping service " + service.GetName() + " because it is up to date"))
			continue
		}

		if err := service.ApplyHelmChart(kubeContext); err != nil {
			return err
		}
	}

	return nil
}

// connect to relevant k8s cluster
func connectToCluster(cfg *config.NopeusConfig, envName string, envData *config.EnvironmentConfig) (string, error) {
	cloudVendor, err := cfg.CAL.GetCloudVendor()
	if err != nil {
		return "", err
	}

	switch cloudVendor {
	case "aws":
		// connect to aws
		return connectToEks(cfg, envName, envData)
	default:
		return "", fmt.Errorf("cloud vendor %s not supported at the moment", cloudVendor)
	}
}

func connectToEks(cfg *config.NopeusConfig, envName string, envData *config.EnvironmentConfig) (string, error) {
	// get terraform outputs
	tfOutputs := envData.GetOutputs()

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
	_, err := exec.Command("aws", cmdArray...).Output()
	// fmt.Println(string(out))
	if err != nil {
		return "", err
	}

	return getActiveKubeContext()
}

// return the active kube context
func getActiveKubeContext() (string, error) {
	kubeconfig, err := sgck.LoadKubeconfig()
	if err != nil {
		return "", err
	}

	return kubeconfig.CurrentContext, nil
}
