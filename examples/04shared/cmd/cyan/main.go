package main

import (
	"github.com/fatih/color"
	"github.com/podhmo/structflag"
	"github.com/podhmo/structflag/examples/04shared/internal"
)

type Options struct {
	internal.BaseConfig
	Output internal.OutputConfig `flag:"output"`
}

func main() {
	options := &Options{}
	if err := structflag.Parse(options); err != nil {
		panic(err)
	}

	internal.SetDebug(options.Debug)
	echo := options.Output.New(color.FgCyan).Echo

	echo("hello")
	echo(color.RedString("warning"))
}
