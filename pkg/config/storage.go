package config

import (
	"context"
	"fmt"
	"os"
	"time"

	helmclient "github.com/mittwald/go-helm-client"
	"github.com/salfatigroup/nopeus/cli/util"
	"github.com/salfatigroup/nopeus/helm"
	"github.com/salfatigroup/nopeus/logger"
)

// define the storage config for the cluster
type Storage struct {
	// database configs
	Database []*DatabaseStorage `yaml:"database"`
}

// define the database storage data structure
type DatabaseStorage struct {
	// the service name
	Name string `yaml:"name"`

	// the database type
	Type string `yaml:"type"`

	// the database version
	Version string `yaml:"version"`
}

// return one of the supported defalt database storage types
func GetDbImage(dbType string) (string, error) {
	switch dbType {
	case "postgres":
		return "bitnami/postgresql-ha", nil
	default:
		return "", fmt.Errorf("unsupported database type: %s", dbType)
	}
}

// define the nopeus storage service template
// this is required because during the installation of storage
// services we require custom actions
type NopeusStorageMicroservice struct {
	Name           string              `yaml:"name"`
	HelmPackage    string              `yaml:"helm_package"`
	ValuesTemplate string              `yaml:"values_template"`
	ValuesPath     string              `yaml:"values_path"`
	Values         *HelmRendererValues `yaml:"values"`
	Namespace      string              `yaml:"namespace"`
	dryRun         bool                `yaml:"dry_run"`
}

// implement the getname function
func (n *NopeusStorageMicroservice) GetName() string {
	return n.Name
}

// implement the get helm package function
func (n *NopeusStorageMicroservice) GetHelmPackage() string {
	return n.HelmPackage
}

// implement the get values template function
func (n *NopeusStorageMicroservice) GetHelmValuesTemplate() string {
	return n.ValuesTemplate
}

// implement the get values path function
func (n *NopeusStorageMicroservice) GetHelmValuesFile() string {
	return n.ValuesPath
}

// implement the get helm values function
func (n *NopeusStorageMicroservice) GetHelmValues() *HelmRendererValues {
	return n.Values
}

// implement the get chart spec function
func (n *NopeusStorageMicroservice) GetChartSpec() (*helmclient.ChartSpec, error) {
	buf, err := os.ReadFile(n.ValuesPath)
	if err != nil {
		return nil, err
	}

	chartSpec := helmclient.ChartSpec{
		ReleaseName:      n.Name,
		ChartName:        n.HelmPackage,
		Version:          "0.1.0",
		ValuesYaml:       string(buf),
		DryRun:           n.dryRun,
		Wait:             false,
		DependencyUpdate: false,
		Timeout:          time.Duration(time.Minute * 15),
	}

	if n.Namespace != "" {
		chartSpec.Namespace = n.Namespace
		chartSpec.CreateNamespace = true
	}

	return &chartSpec, nil
}

// apply the given chart to the cluster
func (n *NopeusStorageMicroservice) ApplyHelmChart(kubeContext string) error {
	// get chart specifications
	chartSpec, err := n.GetChartSpec()
	if err != nil {
		return err
	}

	// get the helm client
	// get client pointing to cert-manager namespace
	helmClient, err := helm.NewHelmClient(chartSpec.Namespace, kubeContext)
	if err != nil {
		return err
	}

	// apply db only once
	if release, _ := helmClient.GetChartByName(n.GetName()); release != nil {
		logger.Debugf("service %s already applied. skipping to avoid password changes", n.GetName())
		return nil
	}

	fmt.Println(util.GrayText("Applying helm chart for service " + n.GetName()))

	// install the chart
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(time.Minute*15))
	defer cancel()
	if _, err := helmClient.Client.InstallOrUpgradeChart(ctx, chartSpec, nil); err != nil {
		return err
	}

	return nil
}

// delete the given chart from the cluster
func (n *NopeusStorageMicroservice) DeleteHelmChart(kubeContext string) error {
	logger.Debugf("removing helm chart for service %s", n.GetName())

	// get chart specifications
	chartSpec, err := n.GetChartSpec()
	if err != nil {
		return err
	}

	// get the helm Client
	helmClient, err := helm.NewHelmClient(chartSpec.Namespace, kubeContext)
	if err != nil {
		return err
	}

	// delete the chart
	return helmClient.Client.UninstallRelease(chartSpec)
}
