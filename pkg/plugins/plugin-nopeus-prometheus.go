package plugins

import (
	"github.com/salfatigroup/nopeus/config"
	"github.com/salfatigroup/nopeus/helm"
	helmrepo "helm.sh/helm/v3/pkg/repo"
)

// define the prometheus plugin
type PrometheusPlugin struct{}

// define the plugin name
func (p *PrometheusPlugin) Name() string {
	return "prometheus"
}

// install the chart repo on init
func (p *PrometheusPlugin) RunOnInit(cfg *config.NopeusConfig) error {
	repo := helmrepo.Entry{
		Name: "prometheus-community",
		URL:  "https://prometheus-community.github.io/helm-charts",
	}

	if err := helm.AddChartRepo(repo); err != nil {
		return err
	}

	return nil
}

// define the plugin logic before generate
func (p *PrometheusPlugin) RunBeforeGenerate(cfg *config.NopeusConfig, envName string, envData *config.EnvironmentConfig) error {
	service := &config.NopeusDefaultMicroservice{
		Name:        "prometheus",
		HelmPackage: "prometheus-community/kube-prometheus-stack",
		Namespace:   "nopeus",
	}

	cfg.Runtime.HelmRuntime.ServiceTemplateData = append(cfg.Runtime.HelmRuntime.ServiceTemplateData, service)
	return nil
}

func (p *PrometheusPlugin) RunAfterGenerate(cfg *config.NopeusConfig, envName string, envData *config.EnvironmentConfig) error {
	return nil
}

func (p *PrometheusPlugin) RunBeforeDeploy(cfg *config.NopeusConfig, envName string, envData *config.EnvironmentConfig) error {
	return nil
}

func (p *PrometheusPlugin) RunAfterDeploy(cfg *config.NopeusConfig, envName string, envData *config.EnvironmentConfig) error {
	return nil
}

func (p *PrometheusPlugin) RunOnFinish(cfg *config.NopeusConfig) error {
	return nil
}
