/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/podhmo/flagstruct/examples/10cobra-sub-command/cmd/internal"
	"github.com/podhmo/flagstruct/examples/10cobra-sub-command/sts"
	"github.com/spf13/cobra"
)

// assumeroleCmd represents the assumerole command
var assumeroleCmd = &cobra.Command{
	Use:   "assume-role",
	Short: "assume-role for sts",
	Long:  `awscli's sts assume-role like command`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("sts assume-role called")
		json.NewEncoder(os.Stdout).Encode(assumeRoleCmdOptions)
	},
}

var assumeRoleCmdOptions = &sts.AssumeRoleOptions{}

func init() {
	internal.BindFlags(assumeroleCmd, assumeRoleCmdOptions)
	stsCmd.AddCommand(assumeroleCmd)
}
