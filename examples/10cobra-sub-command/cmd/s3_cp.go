/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/podhmo/flagstruct/examples/10cobra-sub-command/cmd/internal"
	"github.com/spf13/cobra"
)

type CpOptions struct {
	Dryrun            bool              `flag:"dryrun"`
	Quiet             bool              `flag:"quiet"`
	Include           string            `flag:"include"`
	Exclude           string            `flag:"exclude"`
	Acl               string            `flag:"acl"`
	Metadata          map[string]string `flag:"metadata"`
	MetadataDirective string            `flag:"metadata-directive"`
	ExpectedSize      int               `flag:"expected-size"`
	Recursive         bool              `flag:"recursive"`
}

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

var cpCmdOptions = &CpOptions{}

func init() {
	internal.BindFlags(cpCmd, cpCmdOptions)
	s3Cmd.AddCommand(cpCmd)
}
