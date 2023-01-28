package parser

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/snapp-incubator/barat/internal/config"

	"github.com/pelletier/go-toml/v2"
)

type Language string
type TomlFile map[MessageID]TomlArgs
type MessageID string
type TomlArgs map[string]interface{}

// LoadTomlFiles loads toml files from given paths.
func LoadTomlFiles() (map[Language]TomlFile, error) {
	mapLangToToml := make(map[Language]TomlFile)
	for _, path := range config.C.TomlPaths {
		entries, err := os.ReadDir(path)
		if err != nil {
			return nil, err
		}

		for _, entry := range entries {
			if entry.IsDir() {
				continue
			}

			if filepath.Ext(entry.Name()) == ".toml" {
				nameList := strings.Split(entry.Name(), ".")
				if len(nameList) < 3 {
					return nil, fmt.Errorf("invalid toml file name: %s", entry.Name())
				}

				lang := Language(nameList[1])
				if _, ok := mapLangToToml[lang]; !ok {
					mapLangToToml[lang] = make(TomlFile)
				}

				unmarshalledData := make(TomlFile)
				data, err := os.ReadFile(filepath.Join(path, entry.Name()))
				if err != nil {
					return nil, err
				}

				err = toml.Unmarshal(data, &unmarshalledData)
				if err != nil {
					return nil, err
				}

				for messageID, tomlArgs := range unmarshalledData {
					if _, ok := mapLangToToml[lang][messageID]; ok {
						return nil, fmt.Errorf("duplicate MessageID: %s for language %s in file %s",
							messageID, lang, entry.Name())
					}
					mapLangToToml[lang][messageID] = tomlArgs
				}
			}
		}
	}
	return mapLangToToml, nil
}
