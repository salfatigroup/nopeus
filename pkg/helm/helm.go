package helm

import (
	"context"
	"time"

	helmclient "github.com/mittwald/go-helm-client"
	"github.com/salfatigroup/nopeus/kubernetes"
	"github.com/salfatigroup/nopeus/logger"
	"helm.sh/helm/v3/pkg/release"
	helmrepo "helm.sh/helm/v3/pkg/repo"
)

// define the helm client for nopeus
type HelmClient struct {
	Client helmclient.Client
}

// create a new helm client
func NewHelmClient(namespace, context string) (*HelmClient, error) {
	// get kubeconfig
	kubeconfig, err := kubernetes.GetKubeconfigAsBytes()
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
				// fmt.Printf(format, v...)
				// fmt.Printf("\n")
				logger.Debugf(format, v...)
			},
		},
		KubeContext: context,
		KubeConfig:  kubeconfig,
	}

	// initialize helm client with the correct kube context
	helmClient, err := helmclient.NewClientFromKubeConf(opt, helmclient.Burst(100), helmclient.Timeout(time.Duration(time.Minute*15)))
	if err != nil {
		return nil, err
	}

	return &HelmClient{Client: helmClient}, nil
}

// install a helm chart
func (h *HelmClient) InstallChart(releaseName, chartName, namespace, values string, dryRun bool) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(time.Minute*15))
	defer cancel()

	_, err := h.Client.InstallOrUpgradeChart(ctx, &helmclient.ChartSpec{
		ReleaseName:     releaseName,
		ChartName:       chartName,
		Namespace:       namespace,
		CreateNamespace: namespace != "",
		ValuesYaml:      values,
		DryRun:          dryRun,
		Wait:            true,
		Timeout:         time.Duration(time.Minute * 15),
	}, nil)

	return err
}

// delete a helm chart
func (h *HelmClient) UninstallChart(releaseName string) error {
	return h.Client.UninstallReleaseByName(releaseName)
}

// add a chart to helm
func AddChartRepo(repo helmrepo.Entry) error {
	// create new client without any configs
	helmClient, err := helmclient.New(nil)
	if err != nil {
		return err
	}

	return helmClient.AddOrUpdateChartRepo(repo)
}

// get chart by name
func (h *HelmClient) GetChartByName(name string) (*release.Release, error) {
	return h.Client.GetRelease(name)
}
