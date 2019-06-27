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

var newDatasetCmd = &cobra.Command{
	Use:   "new [dataset]",
	Short: "Generates boilerplate code for a dataset",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if len(dependencies) > 0 {
			generated = true
		}

		for i, d := range dependencies {
			ext := filepath.Ext(d)
			dependencies[i] = strings.TrimSuffix(d, ext)
		}

		if err := datasets.New(args[0], "local", generated, dependencies); err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
	datasetCmd.AddCommand(newDatasetCmd)
	newDatasetCmd.Flags().StringSliceVarP(&dependencies, "depends-on", "d", []string{}, "List of dataset dependencies")
	newDatasetCmd.Flags().BoolVarP(&generated, "generated", "g", false, "Whether this dataset is raw or generated")
}
