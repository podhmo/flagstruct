package main

import (
	"fmt"
	"log"

	"github.com/spf13/cobra"
)

func main() {
	if err := Execute(); err != nil {
		log.Fatalf("!! %+v", err)
	}
}

var rootCmd = &cobra.Command{
	Use:   "hello",
	Short: "hello world",
	Long: `hello world
	hello world`,

	Run: func(cmd *cobra.Command, args []string) {
		name, err := cmd.Flags().GetString("name")
		fmt.Println("hello", name, err, "-", args)
	},
}

func init() {
	rootCmd.Flags().StringP("name", "", "", "name for the user")
	rootCmd.Flags().BoolP("debug", "", false, "debug option")
	rootCmd.MarkFlagRequired("name") // required flag(s) "name" not set
}

func Execute() error {
	return rootCmd.Execute()
}
