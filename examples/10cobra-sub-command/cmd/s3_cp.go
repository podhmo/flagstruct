/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/podhmo/flagstruct/examples/10cobra-sub-command/cmd/internal"
	"github.com/podhmo/flagstruct/examples/10cobra-sub-command/s3"
	"github.com/spf13/cobra"
)

// cpCmd represents the cp command
var cpCmd = &cobra.Command{
	Use:   "cp",
	Short: "cp for s3",
	Long:  `awscli's s3 cp like command`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("s3 cp called")
		json.NewEncoder(os.Stdout).Encode(cpCmdOptions)
	},
}

var cpCmdOptions = &s3.CpOptions{}

func init() {
	internal.BindFlags(cpCmd, cpCmdOptions)
	s3Cmd.AddCommand(cpCmd)
}
