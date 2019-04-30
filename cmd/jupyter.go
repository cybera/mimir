package cmd

import (
	"github.com/spf13/cobra"
)

var jupyterCmd = &cobra.Command{
	Use:   "jupyter",
	Short: "Commands for managing the Jupyter Notebook server",
}

func init() {
	rootCmd.AddCommand(jupyterCmd)
}
