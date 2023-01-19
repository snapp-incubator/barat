package main

import (
	"github.com/alecthomas/kong"
	"github.com/snapp-incubator/barat/cli"
)

func main() {
	ctx := kong.Parse(&cli.CLI)
	err := ctx.Run()
	ctx.FatalIfErrorf(err)
}
