package config

// define the cloud application layer config
// which includes the applications that are needed to be deployed
// to the cloud
type CloudApplicationLayerConfig struct {
    // nopeus config version
    ConfigVersion string `yaml:"version"`

    // the cloud vendor the applications will be deployed to
    CloudVendor string `yaml:"vendor"`

    // supported hosts for TLS and DNS
    Hosts []string `yaml:"hosts"`

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
