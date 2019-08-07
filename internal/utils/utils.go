package utils

import (
	"bufio"
	"fmt"
	"log"
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

func GetInput(reader *bufio.Reader, nonInteractive bool) string {
	if nonInteractive {
		log.Fatal("\nerror: input required in non-interactive mode")
	}

	input, _ := reader.ReadString('\n')
	return Chomp(input)
}

func GetYesNo(reader *bufio.Reader, question string, def, nonInteractive bool) bool {
	for {
		var suffix string

		if def {
			suffix = " [Y/n]: "
		} else {
			suffix = " [y/N]: "
		}

		fmt.Print(question, suffix)
		input := GetInput(reader, nonInteractive)

		switch input {
		case "yes":
			fallthrough
		case "y":
			return true
		case "no":
			fallthrough
		case "n":
			return false
		case "":
			return def
		}
	}
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
