package core

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/salfatigroup/gologsnag"
	"github.com/salfatigroup/nopeus/cache"
	"github.com/salfatigroup/nopeus/cli/util"
	"github.com/salfatigroup/nopeus/config"
	"github.com/salfatigroup/nopeus/logger"
	"github.com/salfatigroup/nopeus/remote"
)

// Deploy the application to the cloud based on
// the provided configurations
func Deploy(cfg *config.NopeusConfig) error {
	// in parallel deploy all the environments
	for envName, envData := range cfg.CAL.GetEnvironments() {
		// load the environment variables for this deployment
		if err := envData.LoadEnvironmentFile(filepath.Dir(cfg.Runtime.ConfigPath)); err != nil {
			return err
		}

		logger.Debugf("Deploying environment %s", envName)
		logger.Publish(&gologsnag.PublishOptions{Event: "deploy", Description: "Deploying environment " + envName})
		// NOTE: unable to parallelize this loop because I'm using
		// the default kubectl to connect with the k8s cluster
		// and the kubectl is not thread safe
		// Ideally, I should use the k8s client library to connect.
		if err := deployEnvironment(envName, envData, cfg); err != nil {
			return err
		}

		logger.Insight(&gologsnag.InsightOptions{Title: "deployments-by-nopeus", Value: 1, Icon: "🛰️"})
		logger.Insight(&gologsnag.InsightOptions{Title: "deployed-apps", Value: len(cfg.CAL.Services), Icon: "🚀"})
	}

	return nil
}

// deploy a single environment to the cloud
func deployEnvironment(envName string, envData *config.EnvironmentConfig, cfg *config.NopeusConfig) error {
	// notify the user
	fmt.Println(util.GrayText("Launching ") + util.GrayText(envName) + util.GrayText(" environment to the cloud"))

	// generate files
	// in parallel, generate the terraform files and the k8s/helm charts and manifests
	// and deploy the application to the cloud
	var wg sync.WaitGroup
	var err1 error
	var err2 error

	// generate the terraform files
	wg.Add(1)
	go func() {
		defer wg.Done()
		logger.Debug("Generating terraform files")
		logger.Publish(&gologsnag.PublishOptions{Event: "generating-terraform-files", Tags: &gologsnag.Tags{"environment": envName}})
		if err1 = generateTerraformFiles(envName, envData, cfg); err1 != nil {
			return
		}
	}()

	// generate the k8s/helm charts and manifests
	wg.Add(1)
	go func() {
		defer wg.Done()
		logger.Debug("Generating k8s/helm charts and manifests")
		logger.Publish(&gologsnag.PublishOptions{Event: "generating-k8s-helm-charts-and-manifests", Tags: &gologsnag.Tags{"environment": envName}})
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

	// unfold nopeus.state files
	_, err := unfoldNopeusState(envName, envData, cfg)
	if err != nil {
		return err
	}

	// deploy the application to the cloud
	if err := deployToCloud(envName, envData, cfg); err != nil {
		return err
	}

	// generate nopeus.state file
	newstate, err := generateNopeusState(envName, envData, cfg)
	if err != nil {
		return err
	}

	// remote caching to nopeus cloud
	if err := setRemoteCache(cfg, newstate); err != nil {
		return err
	}

	return nil
}

// generate the terraform files and write them to the root nopeus directory
func generateNopeusState(envName string, envData *config.EnvironmentConfig, cfg *config.NopeusConfig) (*cache.NopeusState, error) {
	logger.Debug("Generating nopeus state")
	logger.Publish(&gologsnag.PublishOptions{Event: "generating-nopeus-state", Tags: &gologsnag.Tags{"environment": envName}})
	// create the nopeus state
	state, err := cache.NewNopeusState(envName, envData, cfg)
	if err != nil {
		return nil, err
	}

	// write the nopeus state to the root nopeus directory
	nopeusStateLocation := filepath.Join(cfg.Runtime.RootNopeusDir, "state", envName+".nopeus.state")
	if err := state.WriteNopeusState(nopeusStateLocation); err != nil {
		return nil, err
	}
	return state, nil
}

// unfold the nopeus state files to the correct folders as a caching mechanism
func unfoldNopeusState(envName string, envData *config.EnvironmentConfig, cfg *config.NopeusConfig) (*cache.NopeusState, error) {
	logger.Publish(&gologsnag.PublishOptions{Event: "unfolding-nopeus-state", Tags: &gologsnag.Tags{"environment": envName}})
	// get the nopeus state
	nopeusStateLocation := filepath.Join(cfg.Runtime.RootNopeusDir, "state", envName+".nopeus.state")
	// check if file exists
	if _, err := os.Stat(nopeusStateLocation); os.IsNotExist(err) {
		return nil, nil
	}

	// read the nopeus state and parse
	state, err := cache.ReadNopeusState(nopeusStateLocation)
	if err != nil {
		return nil, err
	}
	// unfold the nopeus state
	if err := state.UnfoldNopeusState(cfg); err != nil {
		return nil, err
	}
	return state, nil
}

// // generate the terraform and helm files in parallel
// func generateFiles(envName string, envData *config.EnvironmentConfig, cfg *config.NopeusConfig) error {
//     // in parallel, generate the terraform files and the k8s/helm charts and manifests
//     // and deploy the application to the cloud
//     var wg sync.WaitGroup
//     var err1 error
//     var err2 error

//     // generate the terraform files
//     wg.Add(1)
//     go func() {
//         defer wg.Done()
//         if err1 = generateTerraformFiles(envName, envData, cfg); err1 != nil {
//             return
//         }
//     }()

//     // generate the k8s/helm charts and manifests
//     wg.Add(1)
//     go func() {
//         defer wg.Done()
//         if err2 = generateK8sHelmCharts(envName, envData, cfg); err2 != nil {
//             return
//         }
//     }()
//     // wait for all the goroutines to finish
//     wg.Wait()

//     // validate errors from the goroutines
//     if err1 != nil {
//         return err1
//     } else if err2 != nil {
//         return err2
//     }

//     return nil
// }

// deployToCloud deploys the application to the cloud
// based on the provided configurations, terraform files,
// k8s/helm charts and manifests
func deployToCloud(envName string, envData *config.EnvironmentConfig, cfg *config.NopeusConfig) error {
	// deploy the terraform files
	logger.Debug("Deploying terraform files")
	logger.Publish(&gologsnag.PublishOptions{Event: "deploy-terraform-files", Tags: &gologsnag.Tags{"environment": envName}})
	if err := runTerraform(envName, envData, cfg); err != nil {
		return err
	}

	fmt.Println(
		"🚀",
		util.GradientText("[NOPEUS::MAX-Q::"+strings.ToUpper(envName)+"]", "#db2777", "#f9a8d4"),
		"- applying the cloud configurations",
	)

	// deploy the k8s/helm charts and manifests
	logger.Debug("Deploying k8s/helm charts and manifests")
	logger.Publish(&gologsnag.PublishOptions{Event: "deploy-k8s-helm-charts", Tags: &gologsnag.Tags{"environment": envName}})
	if err := runK8s(envName, envData, cfg); err != nil {
		return err
	}

	return nil
}

func setRemoteCache(cfg *config.NopeusConfig, state *cache.NopeusState) error {
	// create remote session if token is provided
	if cfg.Runtime.NopeusCloudToken != "" {
		logger.Debug("Setting remote cache")
		logger.Publish(&gologsnag.PublishOptions{Event: "set-remote-cache", Tags: &gologsnag.Tags{"environment": state.EnvironmentName}})
		session, err := remote.NewRemoteSession(cfg.Runtime.NopeusCloudToken)
		if err != nil {
			return err
		}

		// remote caching
		if err := session.SetRemoteCache(cfg, state); err != nil {
			return err
		}
	}

	return nil
}

func getRemoteCache(envName string, envData *config.EnvironmentConfig, cfg *config.NopeusConfig) error {
	// create remote session if token is provided
	if cfg.Runtime.NopeusCloudToken != "" {
		logger.Debug("Getting remote cache")
		logger.Publish(&gologsnag.PublishOptions{Event: "get-remote-cache", Tags: &gologsnag.Tags{"environment": envName}})
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
