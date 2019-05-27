package utils

import (
	"os"
	"strings"

	"github.com/cybera/ccds/internal/paths"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v2"
)

func Chomp(s string) string {
	return strings.Trim(s, " \r\n")
}

func WriteConfig() error {
	settings := viper.AllSettings()
	delete(settings, "projectroot")

	bytes, err := yaml.Marshal(settings)
	if err != nil {
		return errors.Wrapf(err, "error writing config")
	}

	file, err := os.Create(paths.Config())
	if err != nil {
		return errors.Wrapf(err, "error writing config")
	}
	defer file.Close()

	if _, err := file.Write(bytes); err != nil {
		return errors.Wrapf(err, "error writing config")
	}

	return nil
}
