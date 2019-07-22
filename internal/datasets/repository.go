package datasets

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/cybera/ccds/internal/utils"

	"github.com/pkg/errors"
	"github.com/spf13/viper"
)

func Create(file string, generated bool, dependencies []string) (Dataset, error) {
	ext := filepath.Ext(file)
	if ext == "" {
		return Dataset{}, errors.New("missing file extension")
	}

	name := strings.TrimSuffix(file, ext)
	if _, err := Get(name); err == nil {
		return Dataset{}, errors.New("dataset already exists")
	}

	for _, dep := range dependencies {
		if _, err := Get(dep); err != nil {
			return Dataset{}, fmt.Errorf("dependency %s does not exist", dep)
		}
	}

	dataset := Dataset{File: file, Generated: generated, Dependencies: dependencies}

	viper.Set("datasets."+name, dataset)
	if err := utils.WriteConfig(); err != nil {
		return dataset, errors.Wrap(err, "error updating project metadate")
	}

	return dataset, nil
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
