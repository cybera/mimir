package cmd

import (
	"log"
	"os"

	"github.com/cybera/ccds/internal/commands"
	"github.com/spf13/cobra"
)

var stopCmd = &cobra.Command{
	Use:   "stop",
	Short: "Stops the Jupyter Notebook server",
	Args:  cobra.ExactArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		stop := commands.DockerCompose("down")
		stop.Stdout = os.Stdout
		stop.Stderr = os.Stderr
		err := stop.Run()
		if err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
	jupyterCmd.AddCommand(stopCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// stopCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// stopCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
