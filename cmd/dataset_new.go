package cmd

import (
	"log"
	"path/filepath"
	"strings"

	"github.com/cybera/ccds/internal/languages"
	"github.com/cybera/ccds/internal/paths"
	"github.com/cybera/ccds/internal/templates"
	"github.com/cybera/ccds/internal/types"
	"github.com/cybera/ccds/internal/utils"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
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

		fileName := args[0]
		ext := filepath.Ext(fileName)
		name := strings.TrimSuffix(fileName, ext)
		if ext == "" {
			log.Fatal("please include the file extension")
		}

		datasets := viper.GetStringMap("datasets")

		for k := range datasets {
			if k == name {
				log.Fatal("dataset already exists")
			}
		}

		for _, dep := range dependencies {
			found := false

			for k := range datasets {
				if k == dep {
					found = true
				}
			}

			if found == false {
				log.Fatalf("dependency %s does not exist", dep)
			}
		}

		datasets[name] = types.Dataset{File: fileName, Generated: generated, Dependencies: dependencies}

		var datasetPath, src string
		lang := viper.GetString("PrimaryLanguage")
		root := viper.GetString("ProjectRoot")

		if generated {
			datasetPath = filepath.Join(root, paths.ProcessedDatasets(), fileName)
		} else {
			datasetPath = filepath.Join(root, paths.RawDatasets(), fileName)
		}

		src = "datasets/load" + languages.Extensions[lang]
		dest := filepath.Join(root, paths.DatasetsCode(), name+languages.Extensions[lang])
		// Relative path from import code directory to dataset file
		relPath, _ := filepath.Rel(filepath.Join(root, paths.DatasetsCode()), datasetPath)

		data := struct {
			Name, RelPath string
		}{
			Name:    name,
			RelPath: relPath,
		}

		log.Printf("Writing dataset import code to %s...", dest)
		if err := templates.WriteFile(src, dest, data); err != nil {
			log.Fatal("failed to generate dataset import code: ", err)
		}

		if err := utils.WriteConfig(); err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
	datasetCmd.AddCommand(newDatasetCmd)
	newDatasetCmd.Flags().StringSliceVarP(&dependencies, "depends-on", "d", []string{}, "List of dataset dependencies")
	newDatasetCmd.Flags().BoolVarP(&generated, "generated", "g", false, "Whether this dataset is raw or generated")
}
