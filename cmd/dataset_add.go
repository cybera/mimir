package cmd

import (
	"log"
	"path/filepath"
	"strings"

	"github.com/cybera/ccds/internal/datasets"
	"github.com/spf13/cobra"
)

var dependencies []string
var generated bool
var from string
var sourceName string

var addDatasetCmd = &cobra.Command{
	Use:   "add [name]",
	Short: "Generates boilerplate code for a dataset",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if len(dependencies) > 0 {
			generated = true
		}

		if !generated && from == "" {
			log.Fatal("Non-generated datasets require --from to be defined")
		}

		for i, d := range dependencies {
			ext := filepath.Ext(d)
			dependencies[i] = strings.TrimSuffix(d, ext)
		}

		source := datasets.Source{Name: sourceName, Target: from, Args: nil}
		dataset := datasets.Dataset{File: args[0], Source: source, Generated: generated, Dependencies: dependencies}

		if err := dataset.GenerateCode(); err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
	datasetCmd.AddCommand(addDatasetCmd)
	addDatasetCmd.Flags().StringSliceVarP(&dependencies, "depends-on", "d", []string{}, "List of dataset dependencies")
	addDatasetCmd.Flags().BoolVarP(&generated, "generated", "g", false, "Whether this dataset is raw or generated")
	addDatasetCmd.Flags().StringVarP(&from, "from", "f", "", "Path to locally stored dataset or remote path")
	addDatasetCmd.Flags().StringVarP(&sourceName, "source", "s", "local", "Remote origin hosting the specified dataset")
}
