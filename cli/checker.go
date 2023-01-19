package cli

import (
	"fmt"
	"strings"

	"github.com/fatih/color"

	"github.com/snapp-incubator/barat/internal/parser"
)

type checkerCmd struct {
	ExcludeKeyRegex []string `short:"e" help:"Path to toml files."`
	Paths           []string `arg:"" name:"path" help:"Paths to load toml files." type:"existingdir"`
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

	errs := parser.Validation(tomlFiles, c.ExcludeKeyRegex)
	if errs != nil {
		for _, err := range errs {
			color.Red(">> " + err.Error())
		}
		return fmt.Errorf("validation failed: %d errors", len(errs))
	}
	return nil
}
