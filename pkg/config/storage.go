package config

// define the storage config for the cluster
type Storage struct {
    // database configs
    Database []DatabaseStorage `yaml:"database"`
}

// define the database storage data structure
type DatabaseStorage struct {
    // the database type
    Type string `yaml:"type"`

    // the database version
    Version string `yaml:"version"`
}
