/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/podhmo/structflag"
	"github.com/spf13/cobra"
)

var lsCmdOptions struct {
	Recursive bool `flag:"recursive"`
}

// lsCmd represents the ls command
var lsCmd = &cobra.Command{
	Use:   "ls",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("s3 ls called")
		json.NewEncoder(os.Stdout).Encode(lsCmdOptions)
	},
}

func init() {
	b := &structflag.Binder{Config: structflag.DefaultConfig()}
	assinByEnvVars := b.Bind(lsCmd.Flags(), &lsCmdOptions)
	lsCmd.PreRunE = func(cmd *cobra.Command, args []string) error {
		return assinByEnvVars(cmd.Flags())
	}

	s3Cmd.AddCommand(lsCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// lsCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// lsCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
