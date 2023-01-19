package parser

import (
	"fmt"
	"regexp"
)

func Validation(tomlFiles map[string]map[string]interface{}, excludeRegex []string) (errs []error) {
	for keyTag, valueTag := range tomlFiles {
		for key, value := range valueTag {
			// check if key is valid for all tags or not
			for tag, _ := range tomlFiles {
				if !isExcluded(key, excludeRegex) {
					if _, ok := tomlFiles[tag][key]; !ok {
						errs = append(errs,
							fmt.Errorf("key \"%s\" not found in tag \"%s\"", key, tag))
					}
				}
			}

			// check if value is valid for all tags or not
			if d, ok := value.(map[string]interface{})["description"]; ok {
				if d == "" {
					errs = append(errs,
						fmt.Errorf("description of key \"%s\" in tag \"%s\" is empty", key, keyTag))
				}
			} else {
				errs = append(errs,
					fmt.Errorf("description of key \"%s\" in tag \"%s\" not found", key, keyTag))
			}

			if o, ok := value.(map[string]interface{})["other"]; ok {
				if o == "" {
					errs = append(errs,
						fmt.Errorf("other of key \"%s\" in tag \"%s\" is empty", key, keyTag))
				}
			} else {
				errs = append(errs,
					fmt.Errorf("other of key \"%s\" in tag \"%s\" not found", key, keyTag))
			}
		}
	}
	return errs
}

func isExcluded(str string, regexes []string) bool {
	for _, regex := range regexes {
		re := regexp.MustCompile(regex)
		rs := re.FindStringSubmatch(str)
		if len(rs) > 0 {
			return true
		}
	}
	return false
}
