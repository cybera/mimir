package utils

import (
	"os"
	"strings"

	"github.com/BurntSushi/toml"
	"github.com/cybera/ccds/internal/paths"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
)

func Chomp(s string) string {
	return strings.Trim(s, " \r\n")
}

func Contains(slice []string, item string) bool {
	for _, x := range slice {
		if x == item {
			return true
		}
	}

	return false
}

func WriteConfig() error {
	settings := viper.AllSettings()
	delete(settings, "projectroot")

	file, err := os.Create(paths.ProjectMetadata())
	if err != nil {
		return errors.Wrapf(err, "error writing config")
	}
	defer file.Close()

	enc := toml.NewEncoder(file)
	if err := enc.Encode(settings); err != nil {
		return errors.Wrapf(err, "error writing config")
	}

	return nil
}
