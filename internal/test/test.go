package test

import (
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/cybera/mimir/internal/utils"
	"github.com/pkg/errors"
)

func CreateTestDir(dir string) error {
	path, _ := filepath.Abs(dir)

	if err := os.Mkdir(path, os.ModePerm); err != nil {
		return errors.Wrap(err, "error creating test dir")
	}

	if err := os.Chdir(path); err != nil {
		return errors.Wrap(err, "error changing to test dir")
	}

	return nil
}

func InitProject(dir, author, license, language string) (string, error) {
	err := CreateTestDir(dir)
	if err != nil {
		return "", err
	}

	return GoRun("init", "-n", "-f", "--author", license, "--license", license, "--language", language)
}

func GoRun(subcommand string, args ...string) (string, error) {
	fullargs := append([]string{"run", "../../main.go", subcommand}, args...)
	return Run("go", fullargs...)
}

func Run(command string, args ...string) (string, error) {
	cmd := exec.Command(command, args...)

	var b strings.Builder
	cmd.Stdout = &b
	cmd.Stderr = &b

	err := cmd.Run()

	return b.String(), err
}

func FileContentsEquals(path, expected string) (bool, error) {
	contents, err := ioutil.ReadFile(path)
	if err != nil {
		return false, errors.Wrap(err, "couldn't read file")
	}

	return utils.Chomp(string(contents)) == utils.Chomp(expected), nil
}
