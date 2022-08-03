package config

// define the service data structure
// each service represent a microservice that will be
// deployed to the final cluster
type Service struct {
    // the docker image to be used
    Image string `yaml:"image"`

    // the image version to deploy
    Version string `yaml:"version"`

    // custom environment variables
    Environment map[string]string `yaml:"environment"`

    // helath check url
    HealthCheckURL string `yaml:"health_url"`

    // amount of replicas
    Replicas int `yaml:"replicas"`

    // custom ingress definitions
    // if not a single ingress is presented the cluster will
    // stay private
    Ingress *Ingress `yaml:"ingress"`

    // extend the final k8s configs with whatever you want
    Extend map[string]interface{} `yaml:"extend"`
}
