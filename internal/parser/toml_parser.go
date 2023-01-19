package parser

import (
	"fmt"
	"os"
	"strings"

	"github.com/pelletier/go-toml/v2"
)

// LoadTomlFiles loads toml files from given paths.
func LoadTomlFiles(paths []string) (map[string]map[string]interface{}, error) {
	tomlFiles := make(map[string]map[string]interface{})
	for _, path := range paths {
		entries, err := os.ReadDir(path)
		if err != nil {
			return nil, err
		}

		for _, entry := range entries {
			if entry.IsDir() {
				continue
			}

			if entry.Name()[len(entry.Name())-5:] == ".toml" {
				tag := strings.Split(entry.Name(), ".")[1]
				if _, ok := tomlFiles[tag]; !ok {
					tomlFiles[tag] = make(map[string]interface{})
				}

				tmp := make(map[string]interface{})

				data, err := os.ReadFile(path + "/" + entry.Name())
				if err != nil {
					return nil, err
				}

				err = toml.Unmarshal(data, &tmp)
				if err != nil {
					return nil, err
				}

				for k, v := range tmp {
					if _, ok := tomlFiles[tag][k]; ok {
						return nil, fmt.Errorf("duplicate key: %s for tag %s in file %s", k, tag, entry.Name())
					}
					tomlFiles[tag][k] = v
				}
			}
		}
	}
	return tomlFiles, nil
}
