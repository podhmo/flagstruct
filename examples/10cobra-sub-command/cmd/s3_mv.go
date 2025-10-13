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

type MvOptions struct {
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

var mvCmdOptions = &MvOptions{}

func init() {
	internal.BindFlags(mvCmd, mvCmdOptions)
	s3Cmd.AddCommand(mvCmd)
}
