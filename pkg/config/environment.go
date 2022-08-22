package config

import (
	"path/filepath"
	"strings"

	"github.com/hashicorp/terraform-exec/tfexec"
	"github.com/joho/godotenv"
	"github.com/salfatigroup/nopeus/helm"
	"github.com/salfatigroup/nopeus/logger"
)

// define the environment configs
type EnvironmentConfig struct {
	EnvFileLocation string `yaml:"env_file"`
	kubeContext     string
	checksumMap     map[string]string
	outputs         map[string]tfexec.OutputMeta
}

func NewEnvironmentConfig() *EnvironmentConfig {
	return &EnvironmentConfig{}
}

// set the terraform outputs to the infrastructure config
func (i *EnvironmentConfig) SetOutputs(outputs map[string]tfexec.OutputMeta) {
	i.outputs = outputs
}

// returns the terraform outputs
func (i *EnvironmentConfig) GetOutputs() map[string]tfexec.OutputMeta {
	return i.outputs
}

// returns the env file to load for the environment
func (i *EnvironmentConfig) GetEnvFileLocation() string {
	return i.EnvFileLocation
}

// load the environment file only if env_file location was provided
func (i *EnvironmentConfig) LoadEnvironmentFile(basepath string) error {
	if i.EnvFileLocation == "" {
		return nil
	}

	// define the file abs location
	file := filepath.Join(basepath, i.EnvFileLocation)
	logger.Debugf("Loading environment file %s", file)

	// load the dotenv files if they exists
	return godotenv.Load(file)
}

// Set the KubeContext to the envData
func (i *EnvironmentConfig) SetKubeContext(kubeContext string) {
	i.kubeContext = kubeContext
}

// return the kubecontext of the environment
func (i *EnvironmentConfig) GetKubeContext() string {
	return i.kubeContext
}

// load checksum map for the environment
func (i *EnvironmentConfig) LoadChecksumMap() error {
	helmClient, err := helm.NewHelmClient("nopeus", i.GetKubeContext())
	if err != nil {
		return err
	}
	logger.Debugf("Loading checksum map for environment %s", i.GetKubeContext())

	// get the checksum chart from helm by name
	checksumChart, err := helmClient.GetChartByName("checksum")
	if err != nil && strings.Contains(err.Error(), "not found") {
		logger.Debugf("No checksum chart found")
		return nil
	} else if err != nil {
		return err
	}
	logger.Debugf("Checksum chart found: %+v", checksumChart)

	// convert map[string]interface{} to map[string]string
	checksumMap := make(map[string]string)
	for k, v := range checksumChart.Config["checksum"].(map[string]interface{}) {
		checksumMap[k] = v.(string)
	}
	logger.Debugf("Checksum map: %+v", checksumMap)

	i.checksumMap = checksumMap
	return nil
}

// get service checksum by name
func (i *EnvironmentConfig) GetChecksum(name string) string {
	hash := i.checksumMap[name]
	logger.Debugf("Service %s checksum is %s", name, hash)
	return hash
}
