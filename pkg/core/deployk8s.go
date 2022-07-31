package core

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"time"

	helmclient "github.com/mittwald/go-helm-client"
	"github.com/salfatigroup/nopeus/config"
)

// creates a helm client based on the kubeconfig and required namespace
func newHelmClient(namespace string, context string) (helmclient.Client, error) {
    // get kubeconfig
    kubeconfig, err := getKubeconfig()
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

        // TODO: implement kube context
        KubeContext: "arn:aws:eks:us-west-1:117583274941:cluster/salfatigroup-nopeus-prod",
        KubeConfig: kubeconfig,
    }

    // initialize helm client with the correct kube context
    helmClient, err := helmclient.NewClientFromKubeConf(opt, helmclient.Burst(100), helmclient.Timeout(time.Duration(time.Minute*15)))
    if err != nil {
        return nil, err
    }

    return helmClient, nil
}

// run and deploy k8s and helm files per environment
func runK8s(cfg *config.NopeusConfig) error {


    // apply the helm charts for each environment
    for _, env := range cfg.Runtime.Environments {
        // run manual helm commands before the generic setup
        // this should be avoid as much as possible
        if err := manualHelmCommands(cfg, env); err != nil {
            return err
        }

        if err := applyK8sHelmCharts(cfg, env); err != nil {
            return err
        }
    }

    return nil
}

func getKubeconfig() ([]byte, error) {
    // get the kubeconfig
    home, err := os.UserHomeDir()
    if err != nil {
        return nil, err
    }

    kubeconfigPath := home + "/.kube/config"
    return ioutil.ReadFile(kubeconfigPath)
}

func manualHelmCommands(cfg *config.NopeusConfig, env string) error {
    // install cert-manager manually
    fmt.Println("Installing cert-manager...")

    // get client pointing to cert-manager namespace
    helmClient, err := newHelmClient("cert-manager", "")
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
        DryRun: false, // cfg.Runtime.DryRun,
        Wait: true,
        Timeout: time.Duration(time.Minute*15),
    }, nil)

    if err != nil {
        return err
    }

    return nil
}

func applyK8sHelmCharts(cfg *config.NopeusConfig, env string) (err error) {
    // apply the helm charts for the environment
    for _, service := range cfg.Runtime.HelmRuntime.ServiceTemplateData {
        if err = applyHelmChart(cfg, service); err != nil {
            return
        }
    }

    return nil
}

// apply a helm chart for the environment
func applyHelmChart(cfg *config.NopeusConfig, service config.ServiceTemplateData) error {
    fmt.Printf("Applying helm chart for service %s\n", service.GetName())
    // get chart specifications
    chartSpec, err := service.GetChartSpec()
    if err != nil {
        return err
    }

    chartSpec.DryRun = false

    // get the helm client
    // get client pointing to cert-manager namespace
    helmClient, err := newHelmClient(chartSpec.Namespace, "")
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
