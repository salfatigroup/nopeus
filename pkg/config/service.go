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
    EnvironmentVariables map[string]string `yaml:"environment"`

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

// return the image name
func (s *Service) GetImage() string {
    return s.Image
}

// return the image version if exists or latest
func (s *Service) GetVersion() string {
    if s.Version == "" {
        return "latest"
    }

    return s.Version
}

// return the environment variables
func (s *Service) GetEnvironmentVariables() map[string]string {
    return s.EnvironmentVariables
}

// return the replicas
func (s *Service) GetReplicas() int {
    return s.Replicas
}

// return the health check url or /status
func (s *Service) GetHealthCheckURL() string {
    if s.HealthCheckURL == "" {
        return "/status"
    }

    return s.HealthCheckURL
}
