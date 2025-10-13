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

type GetCurrentIdentityOptions struct {
	CliInputJson        string `flag:"cli-input-json"`
	GenerateCliSkeleton string `flag:"generate-cli-skeleton"`
}

// getCurrentIdentityCmd represents the getCurrentIdentity command
var getCurrentIdentityCmd = &cobra.Command{
	Use:   "get-current-identity",
	Short: "get-current-identity for sts",
	Long:  `awscli's sts get-current-identity like command`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("sts get-current-identity called")
		json.NewEncoder(os.Stdout).Encode(getCurrentIdentityCmdOptions)
	},
}

var getCurrentIdentityCmdOptions = &GetCurrentIdentityOptions{}

func init() {
	internal.BindFlags(getCurrentIdentityCmd, getCurrentIdentityCmdOptions)
	stsCmd.AddCommand(getCurrentIdentityCmd)
}
