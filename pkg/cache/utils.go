package cache

import (
	"os"
)

// read the tfstate file from the given location
func readTfstate(location string) (string, error) {
	tfstate, err := os.ReadFile(location)
	if err != nil {
		return "", err
	}
	return string(tfstate), nil
}

// get all the keys from map
func getKeys[V any](m map[string]V) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}
