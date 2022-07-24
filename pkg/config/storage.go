package config

import "fmt"

// define the storage config for the cluster
type Storage struct {
    // database configs
    Database []*DatabaseStorage `yaml:"database"`
}

// define the database storage data structure
type DatabaseStorage struct {
    // the service name
    Name string `yaml:"name"`

    // the database type
    Type string `yaml:"type"`

    // the database version
    Version string `yaml:"version"`
}

// return one of the supported defalt database storage types
func GetDbImage(dbType string) (string, error) {
    switch dbType {
        case "postgres":
            return "bitnami/postgresql-ha", nil
        default:
            return "", fmt.Errorf("unsupported database type: %s", dbType)
    }
}
