package config

type HelmRepo struct {
    Name string `yaml:"name"`
    URL string `yaml:"url"`
}

// add all the helm repos
func (c *NopeusConfig) loadHelmRepos() error {
    return nil
}
