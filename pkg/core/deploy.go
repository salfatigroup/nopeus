package core

import (
	"fmt"
	"sync"

	"github.com/salfatigroup/nopeus/cli/util"
	"github.com/salfatigroup/nopeus/config"
	"github.com/salfatigroup/nopeus/remote"
)

// Deploy the application to the cloud based on
// the provided configurations
func Deploy(cfg *config.NopeusConfig) error {
    // in parallel deploy all the environments
    for envName, envData := range cfg.CAL.GetEnvironments() {
        // NOTE: unable to parallelize this loop because I'm using
        // the default kubectl to connect with the k8s cluster
        // and the kubectl is not thread safe
        // Ideally, I should use the k8s client library to connect.
        if err := deployEnvironment(envName, envData, cfg); err != nil {
            return err
        }
    }

    return nil
}

// deploy a single environment to the cloud
func deployEnvironment(envName string, envData *config.EnvironmentConfig, cfg *config.NopeusConfig) error {
    // notify
    fmt.Println(util.GrayText("Launching ") + util.GrayText(envName) + util.GrayText(" environment to the cloud"))
    // in parallel, generate the terraform files and the k8s/helm charts and manifests
    // and deploy the application to the cloud
    var wg sync.WaitGroup
    var err1 error
    var err2 error

    // generate the terraform files
    wg.Add(1)
    go func() {
        defer wg.Done()
        if err1 = generateTerraformFiles(envName, envData, cfg); err1 != nil {
            return
        }
    }()

    // generate the k8s/helm charts and manifests
    wg.Add(1)
    go func() {
        defer wg.Done()
        if err2 = generateK8sHelmCharts(envName, envData, cfg); err2 != nil {
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

    // get remote cache from nopeus cloud
    if err := getRemoteCache(envName, envData, cfg); err != nil {
        return err
    }

    // deploy the application to the cloud
    if err := deployToCloud(envName, envData, cfg); err != nil {
        return err
    }

    // remote caching to nopeus cloud
    if err := setRemoteCache(envName, envData, cfg); err != nil {
        return err
    }

    return nil
}

// deployToCloud deploys the application to the cloud
// based on the provided configurations, terraform files,
// k8s/helm charts and manifests
func deployToCloud(envName string, envData *config.EnvironmentConfig, cfg *config.NopeusConfig) error {
    // deploy the terraform files
    if err := runTerraform(envName, envData, cfg); err != nil {
        return err
    }

    fmt.Println(
        "🚀",
        util.GradientText("[NOPEUS::MAX-Q::" + envName +"]", "#db2777", "#f9a8d4"),
        "- applying the cloud configurations",
    )

    // deploy the k8s/helm charts and manifests
    if err := runK8s(envName, envData, cfg); err != nil {
        return err
    }

    return nil
}

func setRemoteCache(envName string, envData *config.EnvironmentConfig, cfg *config.NopeusConfig) error {
    // create remote session if token is provided
    if cfg.Runtime.NopeusCloudToken != "" {
        session, err := remote.NewRemoteSession(cfg.Runtime.NopeusCloudToken)
        if err != nil {
            return err
        }

        // remote caching
        if err := session.SetRemoteCache(cfg, envName); err != nil {
            return err
        }
    }

    return nil
}

func getRemoteCache(envName string, envData *config.EnvironmentConfig, cfg *config.NopeusConfig) error {
    // create remote session if token is provided
    if cfg.Runtime.NopeusCloudToken != "" {
        session, err := remote.NewRemoteSession(cfg.Runtime.NopeusCloudToken)
        if err != nil {
            return err
        }

        // remote caching
        if err := session.GetRemoteCache(envName); err != nil {
            return err
        }
    }

    return nil
}
