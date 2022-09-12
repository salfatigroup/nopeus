package plugins

import (
	"fmt"
	"path/filepath"

	"github.com/salfatigroup/nopeus/config"
)

type ChecksumPlugin struct{}

// define the name of the plugin
func (p *ChecksumPlugin) Name() string {
	return "nopeus-checksum"
}

// define the plugin logic
func (p *ChecksumPlugin) RunBeforeGenerate(cfg *config.NopeusConfig, envName string, envData *config.EnvironmentConfig) error {
	cloudVendor, err := cfg.CAL.GetCloudVendor()
	if err != nil {
		return err
	}

	checksumMap, err := generateChecksumMap(cfg.Runtime.HelmRuntime.ServiceTemplateData)
	if err != nil {
		return err
	}

	workingDir := filepath.Join(cfg.Runtime.TmpFileLocation, cloudVendor, envName)
	service := &config.NopeusDefaultMicroservice{
		Name:           "checksum",
		HelmPackage:    "salfatigroup/checksum",
		ValuesTemplate: "checksum.values.yaml",
		ValuesPath:     fmt.Sprintf("%s/checksum.values.yaml", workingDir),
		Namespace:      "nopeus",
		DryRun:         cfg.Runtime.DryRun,
		Values: &config.HelmRendererValues{
			Custom: map[string]interface{}{
				"Checksum": checksumMap,
			},
		},
	}

	cfg.Runtime.HelmRuntime.ServiceTemplateData = append(cfg.Runtime.HelmRuntime.ServiceTemplateData, service)

	return nil
}

// generate the chcecksum map for all the given services
func generateChecksumMap(services []config.ServiceTemplateData) (map[string]string, error) {
	checksumMap := make(map[string]string)
	for _, service := range services {
		checksum, err := service.GetChecksum()
		if err != nil {
			return nil, err
		}
		checksumMap[service.GetName()] = checksum
	}
	return checksumMap, nil
}

func (p *ChecksumPlugin) RunOnInit(cfg *config.NopeusConfig) error {
	return nil
}

func (p *ChecksumPlugin) RunAfterGenerate(cfg *config.NopeusConfig, envName string, envData *config.EnvironmentConfig) error {
	return nil
}

func (p *ChecksumPlugin) RunBeforeDeploy(cfg *config.NopeusConfig, envName string, envData *config.EnvironmentConfig) error {
	return nil
}

func (p *ChecksumPlugin) RunAfterDeploy(cfg *config.NopeusConfig, envName string, envData *config.EnvironmentConfig) error {
	return nil
}

func (p *ChecksumPlugin) RunOnFinish(cfg *config.NopeusConfig) error {
	return nil
}
