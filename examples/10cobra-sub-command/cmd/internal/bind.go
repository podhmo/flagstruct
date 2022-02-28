package internal

import (
	"github.com/podhmo/flagstruct"
	"github.com/spf13/cobra"
)

var binder = &flagstruct.Binder{Config: flagstruct.DefaultConfig()}

func BindFlags(cmd *cobra.Command, options interface{}) {
	assinByEnvVars := binder.Bind(cmd.Flags(), options)
	cmd.PreRunE = func(cmd *cobra.Command, args []string) error {
		return assinByEnvVars(cmd.Flags())
	}
}
