package main

import (
	"fmt"

	"github.com/podhmo/flagstruct"
)

type Options struct {
	Name        string `json:"name"`
	Verbose     bool   `json:"verbose" short:"v"`
	Ignored     bool   `json:"ignored" flag:"-"`
	AnotherName string `json:"anotherName" flag:"another-name"`
}

func main() {
	options := &Options{Name: "foo"} // default value

	flagstruct.Parse(options, func(b *flagstruct.Builder) {
		b.Name = "hello"
		b.EnvPrefix = "X_"
		b.FlagnameTags = append(b.FlagnameTags, "json")
	})
	fmt.Printf("parsed: %#+v\n", options)
}
