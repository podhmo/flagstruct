package main

import (
	"fmt"

	"github.com/podhmo/flagstruct"
)

type Options struct {
	Name    string `flag:"name" help:"name of greeting"`
	Verbose bool   `flag:"verbose" short:"v"`
}

func main() {
	options := &Options{Name: "foo"} // default value

	flagstruct.Parse(options, func(b *flagstruct.Builder) {
		b.Name = "hello"
		b.EnvPrefix = "X_"
	})
	fmt.Printf("parsed: %#+v\n", options)
}
