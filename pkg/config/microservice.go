package config

import (
	"bytes"
	"context"
	"crypto/md5"
	"encoding/gob"
	"fmt"
	"os"
	"time"

	helmclient "github.com/mittwald/go-helm-client"
	"github.com/salfatigroup/nopeus/cli/util"
	"github.com/salfatigroup/nopeus/helm"
	"github.com/salfatigroup/nopeus/logger"
)

// define the runtime default microservice template data
type NopeusDefaultMicroservice struct {
	Name           string              `yaml:"name"`
	HelmPackage    string              `yaml:"helm_package"`
	ValuesTemplate string              `yaml:"values_template"`
	ValuesPath     string              `yaml:"values_path"`
	Values         *HelmRendererValues `yaml:"values"`
	Namespace      string              `yaml:"namespace"`
	dryRun         bool                `yaml:"-"`
}

// return the name of the service
func (m *NopeusDefaultMicroservice) GetName() string {
	return m.Name
}

// return the helm package to install
func (m *NopeusDefaultMicroservice) GetHelmPackage() string {
	return m.HelmPackage
}

// return the helm values template
func (m *NopeusDefaultMicroservice) GetHelmValuesTemplate() string {
	return m.ValuesTemplate
}

// return the helm value path location
func (m *NopeusDefaultMicroservice) GetHelmValuesFile() string {
	return m.ValuesPath
}

// return the helm values to use for rendering the helm values file
func (m *NopeusDefaultMicroservice) GetHelmValues() *HelmRendererValues {
	return m.Values
}

// return the chart specification required for the helm chart deployment
func (m *NopeusDefaultMicroservice) GetChartSpec() (*helmclient.ChartSpec, error) {
	buf, err := os.ReadFile(m.ValuesPath)
	if err != nil {
		return nil, err
	}

	chartSpec := helmclient.ChartSpec{
		ReleaseName:      m.Name,
		ChartName:        m.HelmPackage,
		Version:          "0.1.0",
		ValuesYaml:       string(buf),
		DryRun:           m.dryRun,
		Wait:             false, // true, TODO: wait for jobs
		DependencyUpdate: true,
		Timeout:          time.Duration(time.Minute * 15),
	}

	if m.Namespace != "" {
		chartSpec.Namespace = m.Namespace
		chartSpec.CreateNamespace = true
	}

	return &chartSpec, nil
}

// apply the given chart to the cluster
func (m *NopeusDefaultMicroservice) ApplyHelmChart(kubeContext string) error {
	fmt.Println(util.GrayText("Applying helm chart for service " + m.GetName()))
	// get chart specifications
	chartSpec, err := m.GetChartSpec()
	if err != nil {
		return err
	}

	// get the helm client
	// get client pointing to cert-manager namespace
	helmClient, err := helm.NewHelmClient(chartSpec.Namespace, kubeContext)
	if err != nil {
		return err
	}

	// install the chart
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(time.Minute*15))
	defer cancel()
	if _, err := helmClient.Client.InstallOrUpgradeChart(ctx, chartSpec, nil); err != nil {
		return err
	}

	return nil
}

// delete the given chart from the cluster
func (m *NopeusDefaultMicroservice) DeleteHelmChart(kubeContext string) error {
	logger.Debugf("removing helm chart for service %s", m.GetName())

	// get chart specifications
	chartSpec, err := m.GetChartSpec()
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

// return the checksum of the service
func (m *NopeusDefaultMicroservice) GetChecksum() (string, error) {
	var b bytes.Buffer
	if err := gob.NewEncoder(&b).Encode(m); err != nil {
		return "", err
	}
	md5sum := md5.Sum(b.Bytes())
	return fmt.Sprintf("%x", md5sum), nil
}
