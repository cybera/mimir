package datasets

import (
	"log"
	"os"
	"path/filepath"
	"strings"

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

func (d Dataset) FetchAndWrite() error {
	bytes, err := d.Fetch()
	if err != nil {
		return err
	}

	file, err := os.Create(d.AbsPath())
	if err != nil {
		return err
	}
	defer file.Close()

	if _, err := file.Write(bytes); err != nil {
		return err
	}

	return nil
}

func (d Dataset) Fetch() ([]byte, error) {
	if d.Generated {
		return nil, errors.New("can't fetch a generated dataset")
	}

	fetcher, err := fetchers.NewFetcher(d.FetcherConfig)
	if err != nil {
		return nil, err
	}

	bytes, err := fetcher.Fetch()
	if err != nil {
		return nil, err
	}

	return bytes, nil
}

func (d Dataset) GenerateCode() error {
	ext := filepath.Ext(d.File)
	name := strings.TrimSuffix(d.File, ext)

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
