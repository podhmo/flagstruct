package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/podhmo/structflag"
)

type DBConfig struct {
	URI   string `flag:"uri"`
	Debug bool   `flag:"debug"`
}

type Options struct {
	DB        DBConfig `flag:"db"`         // add --db.uri, --db.debug
	AnotherDB DBConfig `flag:"another-db"` // add --another-db.uri, --another-db.debug
}

func main() {
	options := &Options{}
	b := structflag.NewBuilder()
	fs := b.Build(options)
	if err := fs.Parse(os.Args[1:]); err != nil {
		panic(err)
	}
	fmt.Println("parsed")
	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	enc.Encode(options)
}
