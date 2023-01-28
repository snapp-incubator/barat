// Package cli is the commands collector for the CLI application.
package cli

import (
	"fmt"
	"strings"

	"github.com/fatih/color"

	"github.com/snapp-incubator/barat/internal/config"
	"github.com/snapp-incubator/barat/internal/parser"
)

// lintCmd is the struct for the checker command.
type lintCmd struct {
	ConfigPath              string         `help:"Path to config file."`
	TomlPaths               []string       `name:"toml-paths" help:"paths to load toml files." type:"existingdir"`
	ExcludeRegexKeys        []string       `short:"e" help:"exclude keys that match the given regex."`
	MapFunctionNamesToArgNo map[string]int `help:"it's map of the function's name that returns the message by i18n To number of MessageID in arguments."`
	ProjectPath             string         `help:"paths to project for check all files." type:"existingdir"`
	ExcludeFolders          []string       `help:"list of exclude folders for check localization."`
}

// Run runs the checker command.
func (c *lintCmd) Run() error {
	if len(c.ExcludeRegexKeys) > 0 {
		var tmp []string
		for _, regex := range c.ExcludeRegexKeys {
			regex = strings.Replace(regex, "*", "(.*?)", -1)
			tmp = append(tmp, regex)
		}
		c.ExcludeRegexKeys = tmp
	}

	if c.ConfigPath != "" {
		err := config.LoadConfig(c.ConfigPath)
		if err != nil {
			err = fmt.Errorf("error in loading config file: %w", err)
			return err
		}
	} else {
		codeCheck := false
		if c.ProjectPath != "" {
			codeCheck = true
		}
		config.C = &config.Config{
			TomlPaths:    c.TomlPaths,
			ProjectPath:  c.ProjectPath,
			Exclude:      config.Exclude{Folders: c.ExcludeFolders, RegexKeys: c.ExcludeRegexKeys},
			MessageFuncs: config.ToMessageFuncs(c.MapFunctionNamesToArgNo),
			Options: config.Opts{
				Enable: config.Enable{
					TomlCheck:        true,
					DescriptionCheck: true,
					OtherKeyCheck:    true,
					CodeCheck:        codeCheck,
				},
			},
		}
	}

	mapLangToToml, err := parser.LoadTomlFiles()
	if err != nil {
		err = fmt.Errorf("error in loading toml files: %w", err)
		return err
	}

	// check all toml files for duplicate keys and missing keys
	var errs []error
	var cntErr int
	if config.C.Options.Enable.TomlCheck {
		errs = parser.TomlValidation(mapLangToToml)
		cntErr += len(errs)
		if errs != nil {
			errs = append(errs, fmt.Errorf("toml validation failed: %d errors", len(errs)))
			printErrors(errs)
		}
	}

	// check code for localization functions and find keys that are not available in toml files
	if config.C.Options.Enable.CodeCheck {
		errs = parser.CheckCodeForLocalizationFunctions(mapLangToToml, config.C.ProjectPath)
		cntErr += len(errs)
		if errs != nil {
			errs = append(errs, fmt.Errorf("code localization validation failed: %d errors", len(errs)))
			printErrors(errs)
		}
	}

	if len(errs) > 0 {
		return fmt.Errorf("validation failed: %d errors", cntErr)
	}
	return nil
}

// printErrors prints errors in red color.
func printErrors(errs []error) {
	for _, err := range errs {
		color.Red(">> " + err.Error())
	}
}
