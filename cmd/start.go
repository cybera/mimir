package cmd

import (
	"bufio"
	"bytes"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/cybera/mimir/internal/commands"
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

		var out bytes.Buffer
		logs := commands.DockerCompose("logs", "jupyter")
		logs.Stdout = &out

		time.Sleep(3 * time.Second)

		if err = logs.Run(); err != nil {
			log.Fatal("Error retrieving container logs:", err)
		}

		scanner := bufio.NewScanner(&out)
		var token string

		for scanner.Scan() {
			line := scanner.Text()

			// Don't break the loop, we want the latest token in the logs
			if idx := strings.Index(line, "?token="); idx != -1 {
				idx += 7
				if len(line) >= idx {
					token = line[idx:]
				}
			}
		}

		notebooksURL := "http://localhost:8888/?token=" + token
		labURL := "http://localhost:8888/lab?token=" + token
		fmt.Println("\nNotebooks URL:", notebooksURL)
		fmt.Println("Lab URL:", labURL)
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
