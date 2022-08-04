package config

import (
	"fmt"
	"os/exec"
	"path/filepath"

	helmclient "github.com/mittwald/go-helm-client"
)

// define the helm values that will
// be used to render the helm charts
type HelmRendererValues struct {
    // the service name
    Name string `yaml:"name"`

    // the docker image
    Image string `yaml:"image"`

    // the docker image version
    Version string `yaml:"version"`

    // the environment variables
    Environment map[string]string `yaml:"environment"`

    // any other values
    Custom map[string]interface{} `yaml:"-"`
}

// the interface that will be used by each
// service type to implement their own rendering and parsing
type ServiceTemplateData interface {
    // the name of the service
    GetName() string

    // return the helm package that should
    // be applied to the cluster
    GetHelmPackage() string

    // get the helm values template
    GetHelmValuesTemplate() string

    // return the location of the helm values file
    GetHelmValuesFile() string

    // return the values that will be used to render the helm values file
    GetHelmValues() *HelmRendererValues

    // return the chart specification required for the helm chart deployment
    GetChartSpec() (*helmclient.ChartSpec, error)

    // returns the helm command to use for the installation
    GetHelmCommand() (cmd *exec.Cmd)

    // return the uninstall command
    GetUninstallCommand() (cmd *exec.Cmd)
}


// return the default nopeus microservice template data
func NewServiceTemplateData(cfg *NopeusConfig, name string, service *Service, env string) ServiceTemplateData {
    workingDir := filepath.Join(cfg.Runtime.TmpFileLocation, cfg.CAL.CloudVendor, env)

    return &NopeusDefaultMicroservice{
        Name: name,
        HelmPackage: "salfatigroup/default-microservice",
        ValuesTemplate: "service.values.yaml",
        ValuesPath: fmt.Sprintf("%s/%s.values.yaml", workingDir, name),
        Namespace: cfg.Runtime.DefaultNamespace,
        dryRun: cfg.Runtime.DryRun,
        Values: &HelmRendererValues{
            Name: name,
            Image: service.Image,
            Version: service.Version,
            Environment: service.Environment,
            Custom: map[string]interface{}{
                "ImagePullSecret": "dockerconfig",
                "Replicas": service.Replicas,
                "HealthCheckURL": service.HealthCheckURL,
            },
        },
    }
}

// return the database service tempalte data
func NewDatabaseServiceTemplateData(cfg *NopeusConfig, db *DatabaseStorage, env string) ServiceTemplateData {
    workingDir := filepath.Join(cfg.Runtime.TmpFileLocation, cfg.CAL.CloudVendor, env)
    dbImage, err := GetDbImage(db.Type)
    if err != nil {
        // TODO: propegate the error
        panic(err)
    }

    return &NopeusDefaultMicroservice{
        Name: db.Name,
        HelmPackage: "salfatigroup/database",
        ValuesTemplate: "storage.values.yaml",
        ValuesPath: fmt.Sprintf("%s/%s.values.yaml", workingDir, db.Name),
        Namespace: cfg.Runtime.DefaultNamespace,
        dryRun: cfg.Runtime.DryRun,
        Values: &HelmRendererValues{
            Name: db.Name,
            Image: dbImage,
            Version: db.Version,
        },
    }
}