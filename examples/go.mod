module github.com/podhmo/flagstruct/examples

go 1.18

replace github.com/podhmo/flagstruct => ../

require (
	github.com/podhmo/flagstruct v0.4.1
	github.com/spf13/cobra v1.5.0
	github.com/spf13/pflag v1.0.5
)

require github.com/inconshreveable/mousetrap v1.0.0 // indirect
