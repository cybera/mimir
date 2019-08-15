package datasets

import (
	"log"
	"os"
	"path/filepath"

	"github.com/cybera/mimir/internal/fetchers"
	"github.com/cybera/mimir/internal/languages"
	"github.com/cybera/mimir/internal/paths"
	"github.com/cybera/mimir/internal/templates"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
)

type Dataset struct {
	File          string
	FetcherConfig fetchers.FetcherConfig
	Generated     bool
	Dependencies  []string
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

func (d Dataset) Fetch() error {
	if d.Generated {
		return errors.New("can't fetch a generated dataset")
	}

	fetcher, err := fetchers.NewFetcher(d.FetcherConfig)
	if err != nil {
		return err
	}

	file, err := os.Create(d.AbsPath())
	if err != nil {
		return err
	}
	defer file.Close()

	if err := fetcher.Fetch(file); err != nil {
		return err
	}

	return nil
}

func (d Dataset) GenerateCode() error {
	name := canonicalName(d.File)

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

	return nil
}
