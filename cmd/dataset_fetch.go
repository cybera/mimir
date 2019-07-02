package cmd

import (
	"log"
	"os"

	"github.com/cybera/ccds/internal/datasets"
	"github.com/cybera/ccds/internal/fetchers"
	"github.com/spf13/cobra"
)

var source string

var fetchDatasetCmd = &cobra.Command{
	Use:   "fetch [name] [target]",
	Short: "Downloads a dataset from a remote source and generates boilerplate",
	Run: func(cmd *cobra.Command, args []string) {
		switch len(args) {
		case 2:
			src := datasets.Source{Name: source, Target: args[1]}
			dataset := datasets.Dataset{File: args[0], Source: src, Generated: false, Dependencies: nil}

			if err := fetch(dataset); err != nil {
				log.Fatal(err)
			}

			if err := dataset.GenerateCode(); err != nil {
				log.Fatal(err)
			}
		case 1:
			dataset, err := datasets.Get(args[0])
			if err != nil {
				log.Fatal(err)
			}

			if err := fetch(dataset); err != nil {
				log.Fatal(err)
			}
		case 0:
			if err := fetchAll(); err != nil {
				log.Fatal(err)
			}
		default:
			log.Fatal("unexpected number of arguments")
		}
	},
}

func fetchAll() error {
	datasets, err := datasets.GetAll()
	if err != nil {
		return err
	}

	for _, dataset := range datasets {
		if exists, err := dataset.Exists(); err != nil {
			log.Fatal(err)
		} else if exists {
			continue
		}

		if err := fetch(dataset); err != nil {
			log.Println(err)
		}
	}

	return nil
}

func fetch(dataset datasets.Dataset) error {
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

	file, err := os.Create(dataset.AbsPath())
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
