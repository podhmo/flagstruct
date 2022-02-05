package internal

import "log"

type BaseConfig struct {
	Debug bool `flag:"debug"`
}

var debug bool

func SetDebug(value bool) {
	debug = value
}
func Debug(msg string, args ...interface{}) {
	if debug {
		log.Printf(msg, args...)
	}
}
