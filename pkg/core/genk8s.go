package core

import (
	"strconv"

	"github.com/salfatigroup/nopeus/config"
	"github.com/salfatigroup/nopeus/templates"
)

// generateK8sHelmCharts generates the k8s/helm charts and manifests
// based on the provided configurations
func generateK8sHelmCharts(cfg *config.NopeusConfig) error {
    for _, env := range cfg.Runtime.Environments {
        // map the data from the config to the runtime services for helm rendering
        // for each service
        if err := updateHelmRuntime(cfg, env); err != nil {
            return err
        }
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
func updateHelmRuntime(cfg *config.NopeusConfig, env string) error {
    // store all the ingresss from each service
    ingressList := []*config.Ingress{}

    // map the data from the config to the runtime services for helm rendering
    // for each database service
    for _, db := range cfg.CAL.Storage.Database {
        // create a new service template data
        serviceTemplateData := config.NewDatabaseServiceTemplateData(cfg, db, env)

        // add the service template data to the helm runtime
        cfg.Runtime.HelmRuntime.ServiceTemplateData = append(cfg.Runtime.HelmRuntime.ServiceTemplateData, serviceTemplateData)
    }

    // map the data from the config to the runtime services for helm rendering
    // for each service
    for serviceName, service := range cfg.CAL.Services {
        // create a new service template data
        serviceTemplateData := config.NewServiceTemplateData(cfg, serviceName, service, env)

        // add the service template data to the helm runtime
        cfg.Runtime.HelmRuntime.ServiceTemplateData = append(cfg.Runtime.HelmRuntime.ServiceTemplateData, serviceTemplateData)

        // add the ingress to the list if it exists
        if service.Ingress != nil {
            // TODO: consider a better way to add the service name
            // to the root CAL config
            service.Ingress.ServiceName = serviceName
            service.Ingress.Namespace = cfg.Runtime.DefaultNamespace
            if service.Environment["PORT"] != "" {
                port, err := strconv.Atoi(service.Environment["PORT"])
                if err != nil {
                    return err
                }

                service.Ingress.Port = port
            }
            ingressList = append(ingressList, service.Ingress)
        }

        // add the DATABASE_URL to each service for each storage
        for _, db := range cfg.CAL.Storage.Database {
            service.Environment["STORAGE_DATABASE_URL"] = db.Name
        }
    }

    // if ingress exists create a proxy and add it to the helm runtime
    if len(ingressList) > 0 {
        // create a new ingress template data
        ingressTemplateData := config.NewIngressTemplateData(cfg, ingressList, env)

        // add the ingress template data to the helm runtime
        cfg.Runtime.HelmRuntime.ServiceTemplateData = append(cfg.Runtime.HelmRuntime.ServiceTemplateData, ingressTemplateData)
    }

    return nil
}
