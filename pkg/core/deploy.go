package core

import (
	"fmt"
	"sync"

	"github.com/salfatigroup/nopeus/cli/util"
	"github.com/salfatigroup/nopeus/config"
)

// Deploy the application to the cloud based on
// the provided configurations
func Deploy(cfg *config.NopeusConfig) error {
    // in parallel, generate the terraform files and the k8s/helm charts and manifests
    // and deploy the application to the cloud
    var wg sync.WaitGroup
    var err1 error
    var err2 error

    // generate the terraform files
    wg.Add(1)
    go func() {
        defer wg.Done()
        if err1 = generateTerraformFiles(cfg); err1 != nil {
            return
        }
    }()

    // generate the k8s/helm charts and manifests
    wg.Add(1)
    go func() {
        defer wg.Done()
        if err2 = generateK8sHelmCharts(cfg); err2 != nil {
            return
        }
    }()
    // wait for all the goroutines to finish
    wg.Wait()

    // validate errors from the goroutines
    if err1 != nil {
        return err1
    } else if err2 != nil {
        return err2
    }

    // deploy the application to the cloud
    if err := deployToCloud(cfg); err != nil {
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

    fmt.Println(
        "🚀 ",
        util.GradientText("[NOPEUS::MAX-Q]", "#db2777", "#f9a8d4"),
        " - applying the cloud configurations",
    )

    // deploy the k8s/helm charts and manifests
    if err := runK8s(cfg); err != nil {
        return err
    }

    return nil
}
