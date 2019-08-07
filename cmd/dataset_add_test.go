package cmd

import (
	"os"
	"testing"

	"github.com/cybera/mimir/internal/languages"
	"github.com/cybera/mimir/internal/paths"
	"github.com/cybera/mimir/internal/test"
	"github.com/cybera/mimir/internal/utils"
	"github.com/spf13/viper"
)

func TestDataseAdd(t *testing.T) {
	err := test.CreateTestDir("_test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll("_test")
	defer os.Chdir("../")

	for _, d := range []string{".mimir", paths.DatasetsCode()} {
		if err := os.MkdirAll(d, os.ModePerm); err != nil {
			t.Fatalf("error creating directory: %v", err)
		}
	}

	for _, language := range languages.Supported {
		t.Run(language, func(t *testing.T) {
			viper.Set("PrimaryLanguage", language)
			utils.WriteConfig()

			output, err := test.GoRun("dataset", "add", "titanic.csv", "-f", "iceberg.csv")
			if err != nil {
				t.Errorf("process exited with err: %v", err)
			}
			t.Log("output:\n", output)

			// Add generated dataset
			output, err = test.GoRun("dataset", "add", "titanic_clean.csv", "-d=titanic")
			if err != nil {
				t.Errorf("process exited with err: %v", err)
			}
			t.Log("output:\n", output)

			// Add remote dataset
			output, err = test.GoRun("dataset", "add", "titanic_results.csv", "-s", "swift", "-f", "newfoundland/iceberg.csv")
			if err != nil {
				t.Errorf("process exited with err: %v", err)
			}
			t.Log("output:\n", output)
		})
	}
}
