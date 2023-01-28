package parser

import (
	"fmt"
	"regexp"

	"github.com/snapp-incubator/barat/internal/config"
)

func TomlValidation(mapLangToToml map[Language]TomlFile) (errs []error) {
	checkedKeys := map[MessageID]struct{}{}

	for language, tomlFiles := range mapLangToToml { // lang = en, ru, etc. tomls = map[string]interface{}
		for messageID, tomlArgs := range tomlFiles { // key: [keyInTomlFile] value: description, other, etc.
			// check if key is valid for all lang or not
			_, ok := checkedKeys[messageID]
			if !ok {
				for lang, _ := range mapLangToToml {
					if !isExcluded(messageID, config.C.Exclude.RegexKeys) {
						if _, ok := mapLangToToml[lang][messageID]; !ok {
							errs = append(errs,
								fmt.Errorf("MessageID \"%s\" not found in language \"%s\"", messageID, lang))
						}
					}
				}
				checkedKeys[messageID] = struct{}{}
			}

			// check if value is valid for all tags or not
			if config.C.Options.Enable.DescriptionCheck {
				if d, ok := tomlArgs["description"]; ok {
					if d == "" {
						errs = append(errs,
							fmt.Errorf("description key is empty: MessageID \"%s\", language \"%s\"",
								messageID, language))
					}
				} else {
					errs = append(errs,
						fmt.Errorf("description key not found: MessageID \"%s\", language \"%s\"",
							messageID, language))
				}
			}

			if config.C.Options.Enable.OtherKeyCheck {
				if o, ok := tomlArgs["other"]; ok {
					if o == "" {
						errs = append(errs,
							fmt.Errorf("other key is empty: MessageID \"%s\", language \"%s\" ",
								messageID, language))
					}
				} else {
					errs = append(errs,
						fmt.Errorf("other key not found: MessageID \"%s\", language \"%s\"",
							messageID, language))
				}
			}
		}
	}
	return errs
}

func isExcluded(messageID MessageID, regexes []string) bool {
	for _, regex := range regexes {
		re := regexp.MustCompile(regex)
		rs := re.FindStringSubmatch(string(messageID))
		if len(rs) > 0 {
			return true
		}
	}
	return false
}
