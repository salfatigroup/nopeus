package plugins

import (
	"github.com/salfatigroup/nopeus/config"
	"github.com/salfatigroup/nopeus/logger"
)

type Plugin interface {
	Name() string
	RunOnInit(cfg *config.NopeusConfig) error
	RunBeforeGenerate(cfg *config.NopeusConfig, envName string, envData *config.EnvironmentConfig) error
	RunAfterGenerate(cfg *config.NopeusConfig, envName string, envData *config.EnvironmentConfig) error
	RunBeforeDeploy(cfg *config.NopeusConfig, envName string, envData *config.EnvironmentConfig) error
	RunAfterDeploy(cfg *config.NopeusConfig, envName string, envData *config.EnvironmentConfig) error
	RunOnFinish(cfg *config.NopeusConfig) error
}

var PluginManagerInstance = &PluginManager{
	plugins: []Plugin{},
}

// define the plugin manager
type PluginManager struct {
	plugins []Plugin
}

// register a new plugin
func (p *PluginManager) Register(plugin Plugin) {
	logger.Debugf("registering plugin %s", plugin.Name())
	p.plugins = append(p.plugins, plugin)
}

// register all plugins
func (p *PluginManager) RegisterAll() {
	p.Register(&CertManagerPlugin{})
	p.Register(&PrometheusPlugin{})

	// NOTE: the checksum plugin must be registered last
	p.Register(&ChecksumPlugin{})
}

// register all plugins
func RegisterPlugins() {
	logger.Debug("registering plugins")
	PluginManagerInstance.RegisterAll()
}

// run all on init functions for all plugins
func RunOnInit(cfg *config.NopeusConfig) error {
	for _, plugin := range PluginManagerInstance.plugins {
		if err := plugin.RunOnInit(cfg); err != nil {
			return err
		}
	}
	return nil
}

// run all before generate functions for all plugins
func RunBeforeGenerate(cfg *config.NopeusConfig, envName string, envData *config.EnvironmentConfig) error {
	for _, plugin := range PluginManagerInstance.plugins {
		if err := plugin.RunBeforeGenerate(cfg, envName, envData); err != nil {
			return err
		}
	}
	return nil
}

// run all after generate functions for all plugins
func RunAfterGenerate(cfg *config.NopeusConfig, envName string, envData *config.EnvironmentConfig) error {
	for _, plugin := range PluginManagerInstance.plugins {
		if err := plugin.RunAfterGenerate(cfg, envName, envData); err != nil {
			return err
		}
	}
	return nil
}

// run all before deploy functions for all plugins
func RunBeforeDeploy(cfg *config.NopeusConfig, envName string, envData *config.EnvironmentConfig) error {
	for _, plugin := range PluginManagerInstance.plugins {
		if err := plugin.RunBeforeDeploy(cfg, envName, envData); err != nil {
			return err
		}
	}
	return nil
}

// run all after deploy functions for all plugins
func RunAfterDeploy(cfg *config.NopeusConfig, envName string, envData *config.EnvironmentConfig) error {
	for _, plugin := range PluginManagerInstance.plugins {
		if err := plugin.RunAfterDeploy(cfg, envName, envData); err != nil {
			return err
		}
	}
	return nil
}

// run all on finish functions for all plugins
func RunOnFinish(cfg *config.NopeusConfig) error {
	for _, plugin := range PluginManagerInstance.plugins {
		if err := plugin.RunOnFinish(cfg); err != nil {
			return err
		}
	}
	return nil
}
