package cmd

import (
	"os"
	"testing"

	"github.com/cybera/ccds/internal/languages"
	"github.com/cybera/ccds/internal/paths"
	"github.com/cybera/ccds/internal/test"
	"github.com/cybera/ccds/internal/utils"
	"github.com/spf13/viper"
)

func TestDatasetNew(t *testing.T) {
	testDir, err := test.CreateTestDir()
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(testDir)
	defer os.Chdir("../")

	for _, d := range []string{".ccds", paths.DatasetsCode()} {
		if err := os.MkdirAll(d, os.ModePerm); err != nil {
			t.Fatalf("error creating directory: %v", err)
		}
	}

	for _, language := range languages.Supported {
		t.Run(language, func(t *testing.T) {
			viper.Set("PrimaryLanguage", language)
			utils.WriteConfig()

			output, err := test.RunCommand("dataset", "new", "titanic.csv")
			if err != nil {
				t.Errorf("process exited with err: %v", err)
			}
			t.Log("output:\n", output)

			output, err = test.RunCommand("dataset", "new", "titanic_clean.csv", "-d=titanic")
			if err != nil {
				t.Errorf("process exited with err: %v", err)
			}
			t.Log("output:\n", output)
		})
	}
}
