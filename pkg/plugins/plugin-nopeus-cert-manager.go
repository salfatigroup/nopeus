package plugins

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/salfatigroup/nopeus/cli/util"
	"github.com/salfatigroup/nopeus/config"
	"github.com/salfatigroup/nopeus/helm"
	"github.com/salfatigroup/nopeus/logger"
	helmrepo "helm.sh/helm/v3/pkg/repo"
)

type CertManagerPlugin struct{}

// define the name of the plugin
func (p *CertManagerPlugin) Name() string {
	return "nopeus-cert-manager"
}

func manualHelmCommands(cfg *config.NopeusConfig, envName string, envData *config.EnvironmentConfig, kubeContext string) error {
	// check if cert manager exists before installing again
	if helmClient, err := helm.NewHelmClient("cert-manager", kubeContext); err != nil {
		return err
	} else {
		if release, err := helmClient.GetChartByName("cert-manager"); err == nil && release != nil {
			logger.Debugf("Found existing cert-manager release, skipping install")
			return nil
		}
	}

	// install cert-manager manually
	fmt.Println(util.GrayText("Installing cert-manager..."))

	// get client pointing to cert-manager namespace
	helmClient, err := helm.NewHelmClient("cert-manager", kubeContext)
	if err != nil {
		return err
	}

	// install cert-manager
	return helmClient.InstallChart(
		"cert-manager",
		"jetstack/cert-manager",
		"cert-manager",
		"installCRDs: true",
		cfg.Runtime.DryRun,
	)
}

func (p *CertManagerPlugin) RunOnInit(cfg *config.NopeusConfig) error {
	// append the chart repo
	repo := helmrepo.Entry{
		Name: "jetstack",
		URL:  "https://charts.jetstack.io",
	}

	if err := helm.AddChartRepo(repo); err != nil {
		return err
	}

	return nil
}

func (p *CertManagerPlugin) RunBeforeGenerate(cfg *config.NopeusConfig, envName string, envData *config.EnvironmentConfig) error {
	cloudVendor, err := cfg.CAL.GetCloudVendor()
	if err != nil {
		return err
	}

	workingDir := filepath.Join(cfg.Runtime.TmpFileLocation, cloudVendor, envName)
	service := &config.NopeusDefaultMicroservice{
		Name:           "cert-manager-nopeus",
		HelmPackage:    "salfatigroup/cert-manager",
		ValuesTemplate: "cert-manager.values.yaml",
		ValuesPath:     fmt.Sprintf("%s/cert-manager.values.yaml", workingDir),
		Namespace:      cfg.Runtime.DefaultNamespace,
		DryRun:         cfg.Runtime.DryRun,
		Values: &config.HelmRendererValues{
			Name: "cert-manager",
			Custom: map[string]interface{}{
				"Email":   "certificates@salfati.group",
				"Staging": strings.Contains(envName, "staging"),
			},
		},
	}

	cfg.Runtime.HelmRuntime.ServiceTemplateData = append(cfg.Runtime.HelmRuntime.ServiceTemplateData, service)

	return nil
}

func (p *CertManagerPlugin) RunAfterGenerate(cfg *config.NopeusConfig, envName string, envData *config.EnvironmentConfig) error {
	return nil
}

// define the plugin logic
func (p *CertManagerPlugin) RunBeforeDeploy(cfg *config.NopeusConfig, envName string, envData *config.EnvironmentConfig) error {
	// manual cert manager setup
	if err := manualHelmCommands(cfg, envName, envData, envData.GetKubeContext()); err != nil {
		return err
	}

	return nil
}

func (p *CertManagerPlugin) RunAfterDeploy(cfg *config.NopeusConfig, envName string, envData *config.EnvironmentConfig) error {
	return nil
}

func (p *CertManagerPlugin) RunOnFinish(cfg *config.NopeusConfig) error {
	return nil
}
