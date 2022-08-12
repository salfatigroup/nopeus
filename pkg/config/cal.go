package config

import "fmt"

// define the cloud application layer config
// which includes the applications that are needed to be deployed
// to the cloud
type CloudApplicationLayerConfig struct {
    // the stack deployment name
    Name string `yaml:"name"`

    // nopeus config version
    ConfigVersion string `yaml:"version"`

    // the cloud vendor the applications will be deployed to
    CloudVendor string `yaml:"vendor"`

    // the environment that should be setup (prod/stage/dev)
    Environments map[string]*EnvironmentConfig `yaml:"environments"`

    // the applications that will be deployed to the cloud
    Services map[string]*Service `yaml:"services"`

    // define the storage configs
    Storage *Storage `yaml:"storage"`
}

// create a new instance of the cloud application layer config
// with all the required default configs
func NewCloudApplicationLayerConfig() *CloudApplicationLayerConfig {
    return &CloudApplicationLayerConfig{}
}

// return the name or the default nopeus name
func (c *CloudApplicationLayerConfig) GetName() string {
    if c.Name == "" {
        return "nopeus"
    }

    return c.Name
}

// return the config version or the default config version
func (c *CloudApplicationLayerConfig) GetConfigVersion() string {
    if c.ConfigVersion == "" {
        return "latest"
    }

    return c.ConfigVersion
}

// return the cloud vendor or an error if the cloud vendor is not set
func (c *CloudApplicationLayerConfig) GetCloudVendor() (string, error) {
    if c.CloudVendor == "" {
        return "", fmt.Errorf("cloud vendor is not set")
    }

    return c.CloudVendor, nil
}

// return the environments or the default prod environment
func (c *CloudApplicationLayerConfig) GetEnvironments() map[string]*EnvironmentConfig {
    if len(c.Environments) == 0 {
        return map[string]*EnvironmentConfig{
            "prod": NewEnvironmentConfig(),
        }
    }

    return c.Environments
}

// return the services or an error if the services are not set
func (c *CloudApplicationLayerConfig) GetServices() (map[string]*Service, error) {
    if len(c.Services) == 0 {
        return nil, fmt.Errorf("services are not set")
    }

    return c.Services, nil
}

// return the storage config
func (c *CloudApplicationLayerConfig) GetStorage() *Storage {
    return c.Storage
}
