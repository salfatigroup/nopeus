package config

import (
	"fmt"

	helmclient "github.com/mittwald/go-helm-client"
)

func (c *NopeusConfig) loadHelmRepos() error {
    helmClient, err := helmclient.New(nil)
    if err != nil {
        return err
    }

    for _, repo := range c.Runtime.HelmRepos {
        fmt.Printf("Loading helm repo %s\n", repo.Name)
        if err := helmClient.AddOrUpdateChartRepo(*repo); err != nil {
            return err
        }
    }

    return nil
}
