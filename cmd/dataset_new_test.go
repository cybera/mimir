package cmd

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/cybera/ccds/internal/languages"
	"github.com/cybera/ccds/internal/paths"
	"github.com/spf13/viper"
)

func TestDatasetNew(t *testing.T) {
	testDir, _ := filepath.Abs("_test")
	if err := os.Mkdir(testDir, os.ModePerm); err != nil {
		t.Fatalf("error creating test directory: %v", err)
	}
	defer os.RemoveAll(testDir)

	if err := os.Chdir(testDir); err != nil {
		t.Fatalf("error changing to test directory: %v", err)
	}
	defer os.Chdir("../")

	for _, d := range []string{".ccds", paths.DatasetsCode()} {
		if err := os.MkdirAll(d, os.ModePerm); err != nil {
			t.Fatalf("error creating directory: %v", err)
		}
	}

	for _, language := range languages.Supported {
		t.Run(language, func(t *testing.T) {
			viper.Set("PrimaryLanguage", language)
			viper.WriteConfigAs(paths.Config())

			cmd := exec.Command("go", "run", "../../main.go", "dataset", "new", "titanic.csv")
			var b strings.Builder
			cmd.Stdout = &b
			cmd.Stderr = &b
			if err := cmd.Run(); err != nil {
				t.Errorf("process exited with err: %v", err)
			}
			t.Log("output:\n", b.String())

			cmd = exec.Command("go", "run", "../../main.go", "dataset", "new", "titanic_clean.csv", "-d=titanic")
			b = strings.Builder{}
			cmd.Stdout = &b
			cmd.Stderr = &b
			if err := cmd.Run(); err != nil {
				t.Errorf("process exited with err: %v", err)
			}
			t.Log("output:\n", b.String())
		})
	}
}
