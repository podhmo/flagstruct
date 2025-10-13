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

type AssumeRoleOptions struct {
	RoleArn             string   `flag:"role-arn" required:"true"`
	RoleSessionName     string   `flag:"role-session-name" required:"true"`
	Policy              string   `flag:"policy"`
	DurationSeconds     int      `flag:"duration-seconds"`
	TransitiveTagKeys   []string `flag:"transitive-tag-keys"`
	ExternalId          string   `flag:"external-id"`
	SerialNumber        string   `flag:"serial-number"`
	TokenCode           string   `flag:"token-code"`
	SourceIdentity      string   `flag:"source-identity"`
	CliInputJson        string   `flag:"cli-input-json"`
	GenerateCliSkeleton string   `flag:"generate-cli-skeleton"`
}

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

var assumeRoleCmdOptions = &AssumeRoleOptions{}

func init() {
	internal.BindFlags(assumeroleCmd, assumeRoleCmdOptions)
	stsCmd.AddCommand(assumeroleCmd)
}
