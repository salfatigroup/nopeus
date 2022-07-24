package config

import (
	"fmt"
	"path/filepath"
)

// define the ingress data for the service
type Ingress struct {
    // the hostname of the ingress
    Host string `yaml:"host"`

    // supported paths
    Paths []IngressPath `yaml:"paths"`

    // refer to the backend service
    ServiceName string `yaml:"service_name"`

    // the port of the backend service
    Port int `yaml:"port"`
}

// define a single ingress path
type IngressPath struct {
    // the path of the ingress
    Path string `yaml:"path"`

    // stripe the ingress path
    Strip bool `yaml:"strip"`
}

// create ingress service data
func NewIngressTemplateData(cfg *NopeusConfig, ingressList []*Ingress, env string) ServiceTemplateData {
    workingDir := filepath.Join(cfg.Runtime.TmpFileLocation, cfg.CAL.CloudVendor, env)

    return &NopeusDefaultMicroservice{
        Name: "api-gateway",
        HelmPackage: "nopeus/proxy",
        ValuesTemplate: "proxy.values.yaml",
        ValuesPath: fmt.Sprintf("%s/api-gateway.values.yaml", workingDir),
        Values: &HelmRendererValues{
            Name: "api-gateway",
            Image: "kong",
            Version: "latest",
            Custom: map[string]interface{}{
                "Ingress": ingressList,
            },
        },
    }
}
