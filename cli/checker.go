package cli

import (
	"fmt"
	"strings"

	"github.com/fatih/color"

	"github.com/snapp-incubator/barat/internal/parser"
)

type checkerCmd struct {
	Paths                   []string       `arg:"" name:"path" help:"paths to load toml files." type:"existingdir"`
	ExcludeKeyRegex         []string       `short:"e" help:"exclude keys that match the given regex."`
	MapFunctionNamesToArgNo map[string]int `help:"it's map of the function's name that returns the message by i18n To number of MessageID in arguments."`
	ProjectPath             string         `help:"paths to project for check all files." type:"existingdir"`
	ExcludeFolders          []string       `help:"list of exclude folders for check localization."`
}

func (c *checkerCmd) Run() error {
	if c.ExcludeKeyRegex != nil {
		var tmp []string
		for _, regex := range c.ExcludeKeyRegex {
			regex = strings.Replace(regex, "*", "(.*?)", -1)
			tmp = append(tmp, regex)
		}
		c.ExcludeKeyRegex = tmp
	}

	tomlFiles, err := parser.LoadTomlFiles(c.Paths)
	if err != nil {
		return err
	}

	// check all toml files for duplicate keys and missing keys
	errs := parser.TomlValidation(tomlFiles, c.ExcludeKeyRegex)
	if errs != nil {
		errorsShower(errs)
		return fmt.Errorf("toml validation failed: %d errors", len(errs))
	}

	// check code for localization functions and find keys that are not available in toml files
	if c.ProjectPath != "" {
		errs = parser.CheckCodeForLocalizationFunctions(
			tomlFiles, c.ExcludeKeyRegex, c.ExcludeFolders, c.MapFunctionNamesToArgNo, c.ProjectPath,
		)
		if errs != nil {
			errorsShower(errs)
			return fmt.Errorf("code localization validation failed: %d errors", len(errs))
		}
	}
	return nil
}

func errorsShower(errs []error) {
	for _, err := range errs {
		color.Red(">> " + err.Error())
	}
}
