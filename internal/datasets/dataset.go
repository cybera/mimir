package datasets

import (
	"log"
	"path/filepath"
	"strings"

	"github.com/cybera/ccds/internal/languages"
	"github.com/cybera/ccds/internal/paths"
	"github.com/cybera/ccds/internal/templates"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
)

type Dataset struct {
	File          string
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
