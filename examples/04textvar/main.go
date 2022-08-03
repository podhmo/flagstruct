package main

import (
	"encoding/json"
	"fmt"
	"net"
	"os"

	"github.com/podhmo/flagstruct"
)

type Options struct {
	IP  *net.IP `flag:"ip"`
	IP2 net.IP  `flag:"ip2"`
}

func main() {
	options := &Options{IP2: net.IPv4(0, 0, 0, 1)} // default value

	b := flagstruct.NewBuilder()
	b.Name = "textvar"
	b.EnvPrefix = "X_"

	fs := b.Build(options)
	fs.Parse(os.Args[1:])

	fmt.Printf("parsed: %#+v\n", options)
	fmt.Println("----------------------------------------")
	fmt.Println("json:")
	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	enc.Encode(options)
}
