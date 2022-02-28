package main

import (
	"fmt"
	"os"

	"github.com/podhmo/flagstruct"
)

type Options struct {
	Name    string `flag:"name" help:"name of greeting"`
	Verbose bool   `flag:"verbose" short:"v"`
}

func main() {
	options := &Options{Name: "foo"} // default value

	b := flagstruct.NewBuilder()
	b.Name = "hello"
	b.EnvPrefix = "X_"

	fs := b.Build(options)
	fs.Parse(os.Args[1:])

	fmt.Printf("parsed: %#+v\n", options)

}
