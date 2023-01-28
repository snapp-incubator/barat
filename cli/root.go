// Package cli is the commands collector for the CLI application.
package cli

// CLI is the main struct for the CLI application.
var CLI struct {
	Checker checkerCmd `cmd:"" help:"Check all toml file for validation and find message keys in code."`
}
