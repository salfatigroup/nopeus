package config

// define the ingress data for the service
type Ingress struct {
    // the hostname of the ingress
    Host string `yaml:"host"`

    // supported paths
    Paths []IngressPath `yaml:"paths"`
}

// define a single ingress path
type IngressPath struct {
    // the path of the ingress
    Path string `yaml:"path"`

    // stripe the ingress path
    Strip bool `yaml:"strip"`
}
