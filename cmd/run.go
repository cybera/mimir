package cmd

import (
	"log"
	"os"
	"os/exec"

	"github.com/cybera/ccds/internal/commands"
	"github.com/spf13/cobra"
)

var runCmd = &cobra.Command{
	Use:   "run [script]",
	Short: "Runs the specified script from the src/scripts directory",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		var script *exec.Cmd

		if len(args) > 2 {
			script = commands.Script(args[0], args[1:]...)
		} else {
			script = commands.Script(args[0])
		}

		script.Stdout = os.Stdout
		script.Stderr = os.Stderr

		err := script.Run()
		if err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(runCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// runCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// runCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
