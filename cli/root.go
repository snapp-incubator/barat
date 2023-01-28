// Package cli is the root command for the CLI application.
package cli

var CLI struct {
	Checker checkerCmd `cmd:"" help:"Check all toml file for validation and find message keys in code."`
}
