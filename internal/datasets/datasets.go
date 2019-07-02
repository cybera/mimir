package datasets

import (
	"fmt"
	"log"
	"os"
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

func (d Dataset) Exists() (bool, error) {
	if _, err := os.Stat(d.AbsPath()); os.IsNotExist(err) {
		return false, nil
	} else {
		return true, err
	}
}

func Get(name string) (Dataset, error) {
	cleanedName := strings.TrimSuffix(name, filepath.Ext(name))

	var dataset Dataset

	if err := viper.UnmarshalKey("datasets."+cleanedName, &dataset); err != nil {
		return Dataset{}, err
	}

	if dataset.File == "" {
		return dataset, errors.New("dataset does not exist")
	}

	return dataset, nil
}

func GetAll() (map[string]Dataset, error) {
	var datasets map[string]Dataset

	err := viper.UnmarshalKey("datasets", &datasets)

	return datasets, err
}

func (d Dataset) GenerateCode() error {
	ext := filepath.Ext(d.File)
	name := strings.TrimSuffix(d.File, ext)

	if ext == "" {
		return errors.New("missing file extension")
	}

	datasets := viper.GetStringMap("datasets")

	if _, err := Get(name); err == nil {
		return errors.New("dataset already exists")
	}

	for _, dep := range d.Dependencies {
		if _, err := Get(dep); err != nil {
			return fmt.Errorf("dependency %s does not exist", dep)
		}
	}

	datasets[name] = d

	lang := viper.GetString("PrimaryLanguage")
	root := viper.GetString("ProjectRoot")

	src := "datasets/load" + languages.Extensions[lang]
	dest := filepath.Join(root, paths.DatasetsCode(), name+languages.Extensions[lang])
	// Relative path from import code directory to dataset file
	relPath, _ := filepath.Rel(filepath.Join(root, paths.DatasetsCode()), d.AbsPath())

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
