package config

import (
	"io/ioutil"
	"time"

	helmclient "github.com/mittwald/go-helm-client"
)

// define the runtime default microservice template data
type NopeusDefaultMicroservice struct {
    Name string `yaml:"name"`
    HelmPackage string `yaml:"helm_package"`
    ValuesTemplate string `yaml:"values_template"`
    ValuesPath string `yaml:"values_path"`
    Values *HelmRendererValues `yaml:"values"`
    Namespace string `yaml:"namespace"`
    dryRun bool `yaml:"dry_run"`
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
    buf, err := ioutil.ReadFile(m.ValuesPath)
    if err != nil {
        return nil, err
    }

    chartSpec := helmclient.ChartSpec{
        ReleaseName: m.Name,
        ChartName: m.HelmPackage,
        Version: "0.1.0",
        ValuesYaml: string(buf),
        DryRun: m.dryRun,
        Wait: false, // true, TODO: wait for jobs
        DependencyUpdate: true,
        Timeout: time.Duration(time.Minute*15),
    }

    if m.Namespace != "" {
        chartSpec.Namespace = m.Namespace
        chartSpec.CreateNamespace = true
    }

    return &chartSpec, nil
}
