package main

import (
	"github.com/alecthomas/kong"
)

var cli struct {
	Builtin    builtinCmd `cmd:"" help:"generate documentation for the builtin funcers of a given category"`
	Preprocess struct {
		Supports   supportsCmd   `cmd:"" help:"called by mdBook to check support"`
		Preprocess preprocessCmd `cmd:"" help:"do the actual preprocessing" default:"1"`
	} `cmd:"" help:"preprocess a mdBook .md file, formatting Bach code examples"`
}

func main() {
	ctx := kong.Parse(&cli)
	err := ctx.Run()
	ctx.FatalIfErrorf(err)
}
