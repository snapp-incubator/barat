package cli

import (
	"fmt"

	"github.com/fatih/color"

	"github.com/snapp-incubator/barat/internal/parser"
)

type checkerCmd struct {
	Paths []string `arg:"" name:"path" help:"Paths to load toml files." type:"path"`
}

func (c *checkerCmd) Run() error {
	tomlFiles, err := parser.LoadTomlFiles(c.Paths)
	if err != nil {
		return err
	}

	errs := parser.Validation(tomlFiles)
	if errs != nil {
		for _, err := range errs {
			color.Red(">> " + err.Error())
		}
		return fmt.Errorf("validation failed: %d errors", len(errs))
	}
	return nil
}
