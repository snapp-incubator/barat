package parser

import (
	"fmt"
	"regexp"

	"github.com/snapp-incubator/barat/internal/config"
)

func TomlValidation(tomlFiles map[string]map[string]interface{}) (errs []error) {
	checkedKeys := map[string]struct{}{}

	for lang, tomls := range tomlFiles { // lang = en, ru, etc. tomls = map[string]interface{}
		for key, value := range tomls { // key: [keyInTomlFile] value: description, other, etc.
			// check if key is valid for all lang or not
			_, ok := checkedKeys[key]
			if !ok {
				for lan, _ := range tomlFiles {
					if !isExcluded(key, config.C.Exclude.RegexKey) {
						if _, ok := tomlFiles[lan][key]; !ok {
							errs = append(errs,
								fmt.Errorf("key \"%s\" not found in tag \"%s\"", key, lan))
						}
					}
				}
				checkedKeys[key] = struct{}{}
			}

			// check if value is valid for all tags or not
			if config.C.Options.Enable.DescriptionCheck {
				if d, ok := value.(map[string]interface{})["description"]; ok {
					if d == "" {
						errs = append(errs,
							fmt.Errorf("description key is empty: key \"%s\", tag \"%s\"", key, lang))
					}
				} else {
					errs = append(errs,
						fmt.Errorf("description key not found: key \"%s\", tag \"%s\"", key, lang))
				}
			}

			if config.C.Options.Enable.OtherKeyCheck {
				if o, ok := value.(map[string]interface{})["other"]; ok {
					if o == "" {
						errs = append(errs,
							fmt.Errorf("other key is empty: key \"%s\", tag \"%s\" ", key, lang))
					}
				} else {
					errs = append(errs,
						fmt.Errorf("other key not found: key \"%s\", tag \"%s\"", key, lang))
				}
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
