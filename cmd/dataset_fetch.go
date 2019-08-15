package cmd

import (
	"bufio"
	"fmt"
	"log"
	"os"

	"github.com/cybera/mimir/internal/datasets"
	"github.com/cybera/mimir/internal/utils"
	"github.com/spf13/cobra"
)

var fetchDatasetCmd = &cobra.Command{
	Use:   "fetch [name]",
	Short: "Downloads an added dataset",
	Run: func(cmd *cobra.Command, args []string) {
		reader := bufio.NewReader(os.Stdin)

		switch len(args) {
		case 1:
			dataset, err := datasets.Get(args[0])
			if err != nil {
				log.Fatal(err)
			}

			exists, err := dataset.Exists()
			if err != nil {
				log.Fatal(err)
			}

			if exists {
				question := fmt.Sprintf("%s already exists, fetch a new copy?", args[0])
				if !yesToAll && !utils.GetYesNo(reader, question, false, nonInteractive) {
					os.Exit(0)
				}
			}

			log.Println("Fetching dataset...")

			if err := fetch(dataset); err != nil {
				log.Fatal(err)
			}
		case 0:
			question := "This will attempt to fetch all datasets, continue?"
			if !yesToAll && !utils.GetYesNo(reader, question, false, nonInteractive) {
				os.Exit(0)
			}

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

		log.Printf("Fetching dataset %s...", dataset.File)

		if err := fetch(dataset); err != nil {
			log.Println(err)
		}
	}

	return nil
}

func fetch(dataset datasets.Dataset) error {
	if dataset.Generated {
		return nil
	}

	return dataset.Fetch()
}

func init() {
	datasetCmd.AddCommand(fetchDatasetCmd)
}
