package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "todocli",
	Short: "todocli is a utility to help search through upcoming Todo items by a few filters",
}

// Verbose outputs more data to help with debugging
var Verbose bool

// Tag filters the output by a specific tag
var Tag string

// File filters the output by a specific file
var File string

func init() {
	rootCmd.PersistentFlags().BoolVarP(&Verbose, "verbose", "v", false, "verbose output")
	rootCmd.PersistentFlags().StringVarP(&Tag, "tag", "t", "", "tag to filter for")
	rootCmd.PersistentFlags().StringVarP(&File, "file", "f", "", "file to filter for")
}

// Execute executes the root command
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
