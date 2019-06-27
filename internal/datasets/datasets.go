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

type Source struct {
	Name   string
	Target string
	Args   interface{}
}

type Dataset struct {
	File         string
	Source       Source
	Generated    bool
	Dependencies []string
}

func (d Dataset) AbsPath() string {
	root := viper.GetString("ProjectRoot")

	var dir string

	if d.Generated {
		dir = paths.ProcessedDatasets()
	} else {
		dir = paths.RawDatasets()
	}

	return filepath.Join(root, dir, d.File)
}

func New(filename string, source Source, generated bool, dependencies []string) error {
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

	dataset := Dataset{File: filename, Source: source, Generated: generated, Dependencies: dependencies}
	datasets[name] = dataset

	lang := viper.GetString("PrimaryLanguage")
	root := viper.GetString("ProjectRoot")

	src := "datasets/load" + languages.Extensions[lang]
	dest := filepath.Join(root, paths.DatasetsCode(), name+languages.Extensions[lang])
	// Relative path from import code directory to dataset file
	relPath, _ := filepath.Rel(filepath.Join(root, paths.DatasetsCode()), dataset.AbsPath())

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
