package cmd

import (
	"log"

	"github.com/cybera/mimir/internal/datasets"
	"github.com/cybera/mimir/internal/fetchers"
	"github.com/spf13/cobra"
)

var dependencies []string
var generated bool
var from string
var source string

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

		fetcherConfig := fetchers.FetcherConfig{Name: source, From: from}
		_, err := fetchers.NewFetcher(fetcherConfig)
		if err != nil {
			log.Fatal(err)
		}

		dataset, err := datasets.Create(args[0], fetcherConfig, generated, dependencies)
		if err != nil {
			log.Fatal(err)
		}

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
	addDatasetCmd.Flags().StringVarP(&source, "source", "s", "local", "Remote origin hosting the specified dataset")
}
