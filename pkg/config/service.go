package config

import (
	"fmt"
	"os"
	"strings"

	"github.com/salfatigroup/nopeus/logger"
)

// define the service data structure
// each service represent a microservice that will be
// deployed to the final cluster
type Service struct {
	// the docker image to be used
	Image string `yaml:"image"`

	// the image version to deploy
	Version string `yaml:"version"`

	// raw environment variables
	EnvironmentVariables map[string]string `yaml:"environment"`
	// parsed enviroenment variables
	// envVars [environment name] [environment variable name] environment variable value
	envVars map[string]map[string]string `yaml:"-"`

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

// add an environment variable at runtime
func (s *Service) AddEnvironmentVariable(key, value, envName string) {
	if s.envVars == nil {
		s.envVars = make(map[string]map[string]string)
	}

	if s.envVars[envName] == nil {
		s.envVars[envName] = make(map[string]string)
	}

	s.envVars[envName][key] = value
}

// return the image version if exists or latest
func (s *Service) GetVersion() string {
	if s.Version == "" {
		return "latest"
	}

	return s.Version
}

func (s *Service) ParseEnvironmentVariables(envName string) error {
	if s.envVars == nil {
		s.envVars = make(map[string]map[string]string)
	}

	if s.envVars[envName] == nil {
		s.envVars[envName] = make(map[string]string)
	}

	// iterate through each environment variable
	for key, value := range s.GetRawEnvironmentVariables() {
		logger.Debugf("converting env variables - key: %s, value: %s", key, value)
		// check if the value is in the following format ${ENV_VAR}
		if strings.HasPrefix(value, "${") && strings.HasSuffix(value, "}") {
			// get the environment variable name
			envVar := value[2 : len(value)-1]
			// get the environment variable value
			envValue := os.Getenv(envVar)
			// check if the environment variable value is empty
			if envValue == "" {
				return fmt.Errorf("environment variable %s is not set", envVar)
			}
			// update the environment variable value
			s.envVars[envName][key] = envValue
		} else {
			// update the environment variable value
			s.envVars[envName][key] = value
		}
	}

	return nil
}

// return the environment variables
func (s *Service) GetRawEnvironmentVariables() map[string]string {
	return s.EnvironmentVariables
}

// return the parsed environment variables
func (s *Service) GetEnvironmentVariables(envName string) map[string]string {
	return s.envVars[envName]
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
