package cmd

import (
	"log"
	"os"
	"path/filepath"

	"github.com/cybera/ccds/internal/datasets"
	"github.com/cybera/ccds/internal/fetchers"
	"github.com/cybera/ccds/internal/paths"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var source string

var fetchDatasetCmd = &cobra.Command{
	Use:   "fetch [name] [target]",
	Short: "Downloads a dataset from a remote source and generates boilerplate",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		name := args[0]
		target := args[1]

		fetcher, err := fetchers.NewFetcher(source, target, nil)
		if err != nil {
			log.Fatal(err)
		}

		bytes, err := fetcher.Fetch()
		if err != nil {
			log.Fatal(err)
		}

		root := viper.GetString("ProjectRoot")
		path := filepath.Join(root, paths.RawDatasets(), name)

		file, err := os.Create(path)
		if err != nil {
			log.Fatal(err)
		}
		defer file.Close()

		if _, err := file.Write(bytes); err != nil {
			log.Fatal(err)
		}

		if err := datasets.New(name, source, false, nil); err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
	datasetCmd.AddCommand(fetchDatasetCmd)
	fetchDatasetCmd.Flags().StringVarP(&source, "source", "s", "swift", "How to access the dataset")
}
