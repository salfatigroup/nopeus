package core

import (
	"strconv"

	"github.com/salfatigroup/nopeus/config"
	"github.com/salfatigroup/nopeus/templates"
)

// generateK8sHelmCharts generates the k8s/helm charts and manifests
// based on the provided configurations
func generateK8sHelmCharts(envName string, envData *config.EnvironmentConfig, cfg *config.NopeusConfig) error {
	// map the data from the config to the runtime services for helm rendering
	// for each service
	if err := updateHelmRuntime(cfg, envName, envData); err != nil {
		return err
	}

	// render the helm charts fo each service in the runtime
	for _, serviceTemplateData := range cfg.Runtime.HelmRuntime.ServiceTemplateData {
		// render the helm values file
		if err := templates.RenderHelmTemplateFile(serviceTemplateData); err != nil {
			return err
		}
	}

	return nil
}

// updateHelmRuntime updates the helm runtime data
// based on the provided CAL configurations
func updateHelmRuntime(cfg *config.NopeusConfig, envName string, envData *config.EnvironmentConfig) error {
	// store all the ingresss from each service
	ingressList := []*config.Ingress{}

	// apply cert-manager to the cluster (salfatigroup/cert-manager)
	certManagerTemplateData, err := config.NewCertManagerTemplateData(cfg, envName)
	if err != nil {
		return err
	}

	// add the service template data to the helm runtime
	cfg.Runtime.HelmRuntime.ServiceTemplateData = append(cfg.Runtime.HelmRuntime.ServiceTemplateData, certManagerTemplateData)

	// map the data from the config to the runtime services for helm rendering
	// for each database service
	for _, db := range cfg.CAL.GetStorage().Database {
		// create a new service template data
		serviceTemplateData, err := config.NewDatabaseServiceTemplateData(cfg, db, envName)
		if err != nil {
			return err
		}

		// add the service template data to the helm runtime
		cfg.Runtime.HelmRuntime.ServiceTemplateData = append(cfg.Runtime.HelmRuntime.ServiceTemplateData, serviceTemplateData)
	}

	// map the data from the config to the runtime services for helm rendering
	// for each service
	services, err := cfg.CAL.GetServices()
	if err != nil {
		return err
	}
	for serviceName, service := range services {
		// create a new service template data
		serviceTemplateData, err := config.NewServiceTemplateData(cfg, serviceName, service, envName, envData)
		if err != nil {
			return err
		}

		// add the service template data to the helm runtime
		cfg.Runtime.HelmRuntime.ServiceTemplateData = append(cfg.Runtime.HelmRuntime.ServiceTemplateData, serviceTemplateData)

		// add the ingress to the list if it exists
		if service.Ingress != nil {
			// TODO: consider a better way to add the service name
			// to the root CAL config
			service.Ingress.ServiceName = serviceName
			service.Ingress.Namespace = cfg.Runtime.DefaultNamespace
			if service.GetEnvironmentVariables(envName)["PORT"] != "" {
				port, err := strconv.Atoi(service.EnvironmentVariables["PORT"])
				if err != nil {
					return err
				}

				service.Ingress.Port = port
			}
			ingressList = append(ingressList, service.Ingress)
		}

		// add the DATABASE_URL to each service for each storage
		for _, db := range cfg.CAL.GetStorage().Database {
			service.AddEnvironmentVariable("STORAGE_DATABASE_URL", db.Name, envName)
		}
	}

	// if ingress exists create a proxy and add it to the helm runtime
	if len(ingressList) > 0 {
		// create a new ingress template data
		ingressTemplateData, err := config.NewIngressTemplateData(cfg, ingressList, envName)
		if err != nil {
			return err
		}

		// add the ingress template data to the helm runtime
		cfg.Runtime.HelmRuntime.ServiceTemplateData = append(cfg.Runtime.HelmRuntime.ServiceTemplateData, ingressTemplateData)
	}

	// NOTE: this must come last MUST come LAST!!!
	checksumMap, err := config.NewChecksumTemplateData(cfg, envName)
	if err != nil {
		return err
	}
	cfg.Runtime.HelmRuntime.ServiceTemplateData = append(cfg.Runtime.HelmRuntime.ServiceTemplateData, checksumMap)

	return nil
}
