package main

import (
	"fmt"
	"os"

	"github.com/podhmo/structflag"
)

type Options struct {
	Name    string `json:"name"`
	Verbose bool   `json:"verbose" short:"v"`
}

func main() {
	options := &Options{Name: "foo"} // default value

	b := structflag.NewBuilder()
	b.Name = "hello"
	b.EnvPrefix = "X_"
	b.FlagnameTags = append(b.FlagnameTags, "json")

	fs := b.Build(options)
	fs.Parse(os.Args[1:])

	fmt.Printf("parsed: %#+v\n", options)

}
