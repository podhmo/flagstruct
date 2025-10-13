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

// mvCmd represents the mv command
var mvCmd = &cobra.Command{
	Use:   "mv",
	Short: "mv for s3",
	Long:  `awscli's s3 mv like command`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("s3 mv called")
		json.NewEncoder(os.Stdout).Encode(mvCmdOptions)
	},
}

var mvCmdOptions = &s3.MvOptions{}

func init() {
	internal.BindFlags(mvCmd, mvCmdOptions)
	s3Cmd.AddCommand(mvCmd)
}
