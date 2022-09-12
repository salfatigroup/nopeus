package config

import (
	"fmt"
	"path/filepath"
)

// define the ingress data for the service
type Ingress struct {
	// supported paths
	Paths []IngressPath `yaml:"paths"`

	// refer to the backend service
	ServiceName string `yaml:"service_name"`

	// the port of the backend service
	Port int `yaml:"port"`

	// the namespace the upstream resides in
	Namespace string `yaml:"namespace"`
}

// define a single ingress path
type IngressPath struct {
	// the path of the ingress
	Path string `yaml:"path"`

	// stripe the ingress path
	Strip bool `yaml:"strip"`

	// the host domains
	Hosts []string `yaml:"hosts"`
}

// create ingress service data
func NewIngressTemplateData(cfg *NopeusConfig, ingressList []*Ingress, env string) (ServiceTemplateData, error) {
	cloudVendor, err := cfg.CAL.GetCloudVendor()
	if err != nil {
		return &NopeusDefaultMicroservice{}, err
	}

	custom := map[string]interface{}{
		"Ingress":    ingressList,
		"HostPrefix": "",
	}

	if env != "prod" {
		custom["HostPrefix"] = env
	}

	workingDir := filepath.Join(cfg.Runtime.TmpFileLocation, cloudVendor, env)

	return &NopeusDefaultMicroservice{
		Name:           "api-gateway",
		HelmPackage:    "salfatigroup/proxy",
		ValuesTemplate: "proxy.values.yaml",
		ValuesPath:     fmt.Sprintf("%s/api-gateway.values.yaml", workingDir),
		Namespace:      "apigw",
		DryRun:         cfg.Runtime.DryRun,
		Values: &HelmRendererValues{
			Name:    "api-gateway",
			Image:   "kong",
			Version: "latest",
			Custom:  custom,
		},
	}, nil
}
