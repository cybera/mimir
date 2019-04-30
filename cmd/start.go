package cmd

import (
	"log"
	"os"

	"github.com/cybera/ccds/internal/commands"
	"github.com/spf13/cobra"
)

var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Starts the Jupyter Notebook server",
	Args:  cobra.ExactArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		start := commands.DockerCompose("up", "-d")
		start.Stdout = os.Stdout
		start.Stderr = os.Stderr
		err := start.Run()
		if err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
	jupyterCmd.AddCommand(startCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// startCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// startCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
