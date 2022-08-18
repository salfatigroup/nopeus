package config

// generate the chcecksum map for all the given services
func generateChecksumMap(services []ServiceTemplateData) (map[string]string, error) {
	checksumMap := make(map[string]string)
	for _, service := range services {
		checksum, err := service.GetChecksum()
		if err != nil {
			return nil, err
		}
		checksumMap[service.GetName()] = checksum
	}
	return checksumMap, nil
}
