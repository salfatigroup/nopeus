package core

import (
	"sync"

	"github.com/salfatigroup/nopeus/config"
)

// Deploy the application to the cloud based on
// the provided configurations
func Deploy(cfg *config.NopeusConfig) (err error) {
    // in parallel, generate the terraform files and the k8s/helm charts and manifests
    // and deploy the application to the cloud
    var wg sync.WaitGroup

    // generate the terraform files
    wg.Add(1)
    go func() {
        defer wg.Done()
        if err = generateTerraformFiles(cfg); err != nil {
            return
        }
    }()

    // generate the k8s/helm charts and manifests
    wg.Add(1)
    go func() {
        defer wg.Done()
        if err = generateK8sHelmCharts(cfg); err != nil {
            return
        }
    }()
    // wait for all the goroutines to finish
    wg.Wait()

    // validate errors from the goroutines
    if err != nil {
        return err
    }

    // deploy the application to the cloud
    if err = deployToCloud(cfg); err != nil {
        return err
    }

    return nil
}

// deployToCloud deploys the application to the cloud
// based on the provided configurations, terraform files,
// k8s/helm charts and manifests
func deployToCloud(cfg *config.NopeusConfig) error {
    // deploy the terraform files
    if err := runTerraform(cfg); err != nil {
        return err
    }

    // deploy the k8s/helm charts and manifests
    if err := runK8s(cfg); err != nil {
        return err
    }

    return nil
}
