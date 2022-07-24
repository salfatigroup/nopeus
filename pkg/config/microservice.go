package config

import (
	"os/exec"
)

// define the runtime default microservice template data
type NopeusDefaultMicroservice struct {
    Name string `yaml:"name"`
    HelmPackage string `yaml:"helm_package"`
    ValuesTemplate string `yaml:"values_template"`
    ValuesPath string `yaml:"values_path"`
    Values *HelmRendererValues `yaml:"values"`
}

// return the name of the service
func (m *NopeusDefaultMicroservice) GetName() string {
    return m.Name
}

// return the helm package to install
func (m *NopeusDefaultMicroservice) GetHelmPackage() string {
    return m.HelmPackage
}

// return the helm values template
func (m *NopeusDefaultMicroservice) GetHelmValuesTemplate() string {
    return m.ValuesTemplate
}

// return the helm value path location
func (m *NopeusDefaultMicroservice) GetHelmValuesFile() string {
    return m.ValuesPath
}

// return the helm values to use for rendering the helm values file
func (m *NopeusDefaultMicroservice) GetHelmValues() *HelmRendererValues {
    return m.Values
}

// return the helm command to use for the installation
func (m *NopeusDefaultMicroservice) GetHelmCommand() (cmd *exec.Cmd) {
    return exec.Command("helm", "upgrade", "--install", m.Name, m.HelmPackage, "--values", m.ValuesPath)
}

// return the uninstall command
func (m *NopeusDefaultMicroservice) GetUninstallCommand() (cmd *exec.Cmd) {
    return exec.Command("helm", "uninstall", m.Name)
}
