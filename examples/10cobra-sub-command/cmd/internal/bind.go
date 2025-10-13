package internal

import (
	"github.com/podhmo/flagstruct"
	"github.com/spf13/cobra"
)

func BindFlags(cmd *cobra.Command, options interface{}) {
	binder := &flagstruct.Binder{Config: flagstruct.DefaultConfig()}
	assignByEnvVars := binder.Bind(cmd.Flags(), options)
	cmd.PreRunE = func(cmd *cobra.Command, args []string) error {
		if err := assignByEnvVars(cmd.Flags()); err != nil {
			return err
		}
		return binder.ValidateRequiredFlags(cmd.Flags())
	}
}
