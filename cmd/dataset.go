package cmd

import (
	"github.com/spf13/cobra"
)

var datasetCmd = &cobra.Command{
	Use:   "dataset",
	Short: "Commands for managing datasets",
}

func init() {
	rootCmd.AddCommand(datasetCmd)
}
