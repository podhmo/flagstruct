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

// lsCmd represents the ls command
var lsCmd = &cobra.Command{
	Use:   "ls",
	Short: "ls for s3",
	Long: `awscli's s3 ls like command`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("s3 ls called")
		json.NewEncoder(os.Stdout).Encode(lsCmdOptions)
	},
}

var lsCmdOptions = &s3.LsOptions{PageSize: 10}

func init() {
	internal.BindFlags(lsCmd, lsCmdOptions)
	s3Cmd.AddCommand(lsCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// lsCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// lsCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
