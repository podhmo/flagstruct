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

type LsOptions struct {
	Recursive     bool   `flag:"recursive"`
	PageSize      int    `flag:"page-size"`
	HumanReadable bool   `flag:"human-redable"`
	Summarize     bool   `flag:"summarize"`
	RequestPlayer string `flag:"request-player"`
}

// lsCmd represents the ls command
var lsCmd = &cobra.Command{
	Use:   "ls",
	Short: "ls for s3",
	Long:  `awscli's s3 ls like command`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("s3 ls called")
		json.NewEncoder(os.Stdout).Encode(lsCmdOptions)
	},
}

var lsCmdOptions = &LsOptions{PageSize: 10}

func init() {
	internal.BindFlags(lsCmd, lsCmdOptions)
	s3Cmd.AddCommand(lsCmd)
}
