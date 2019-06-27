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
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 2 {
			src := datasets.Source{Name: source, Target: args[1]}
			dataset := datasets.Dataset{File: args[0], Source: src, Generated: false, Dependencies: nil}

			if err := fetch(dataset); err != nil {
				log.Fatal(err)
			}

			if err := datasets.New(dataset.File, src, false, nil); err != nil {
				log.Fatal(err)
			}
		} else if len(args) == 0 {
			if err := fetchAll(); err != nil {
				log.Fatal(err)
			}
		} else {
			log.Fatal("unexpected number of arguments")
		}
	},
}

func fetchAll() error {
	var datasets map[string]datasets.Dataset

	viper.UnmarshalKey("Datasets", &datasets)

	for _, dataset := range datasets {
		if err := fetch(dataset); err != nil {
			log.Println(err)
		}
	}

	return nil
}

func fetch(dataset datasets.Dataset) error {
	name := dataset.File
	target := dataset.Source.Target

	if dataset.Source.Name == "local" {
		return nil
	}

	fetcher, err := fetchers.NewFetcher(source, target, nil)
	if err != nil {
		return err
	}

	bytes, err := fetcher.Fetch()
	if err != nil {
		return err
	}

	root := viper.GetString("ProjectRoot")
	path := filepath.Join(root, paths.RawDatasets(), name)

	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	if _, err := file.Write(bytes); err != nil {
		return err
	}

	return nil
}

func init() {
	datasetCmd.AddCommand(fetchDatasetCmd)
	fetchDatasetCmd.Flags().StringVarP(&source, "source", "s", "swift", "How to access the dataset")
}
