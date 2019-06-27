package datasets

import (
	"fmt"
	"log"
	"path/filepath"
	"strings"

	"github.com/cybera/ccds/internal/languages"
	"github.com/cybera/ccds/internal/paths"
	"github.com/cybera/ccds/internal/templates"
	"github.com/cybera/ccds/internal/utils"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
)

type Dataset struct {
	File         string
	Source       string
	Generated    bool
	Dependencies []string
}

func New(filename, source string, generated bool, dependencies []string) error {
	ext := filepath.Ext(filename)
	name := strings.TrimSuffix(filename, ext)

	if ext == "" {
		return errors.New("please include the file extension")
	}

	datasets := viper.GetStringMap("datasets")

	for k := range datasets {
		if k == name {
			return errors.New("dataset already exists")
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
			return fmt.Errorf("dependency %s does not exist", dep)
		}
	}

	datasets[name] = Dataset{File: filename, Source: "local", Generated: generated, Dependencies: dependencies}

	var datasetPath, src string
	lang := viper.GetString("PrimaryLanguage")
	root := viper.GetString("ProjectRoot")

	if generated {
		datasetPath = filepath.Join(root, paths.ProcessedDatasets(), filename)
	} else {
		datasetPath = filepath.Join(root, paths.RawDatasets(), filename)
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
		return errors.Wrap(err, "failed to generate dataset import code")
	}

	return utils.WriteConfig()
}
