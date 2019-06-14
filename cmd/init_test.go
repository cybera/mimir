package cmd

import (
	"os"
	"testing"

	"github.com/cybera/ccds/internal/languages"
	"github.com/cybera/ccds/internal/test"
)

func TestInit(t *testing.T) {
	author := "John Doe"
	license := "MIT"

	// Make sure cached templates won't fail tests
	os.Chdir("../")
	if output, err := test.Run("packr2", "clean"); err != nil {
		t.Logf("failed to clear packr2 cache: %s", output)
	}
	os.Chdir("cmd")

	for _, language = range languages.Supported {
		t.Run(language, func(t *testing.T) {
			testDir, err := test.CreateTestDir()
			if err != nil {
				t.Fatal(err)
			}
			defer os.RemoveAll(testDir)
			defer os.Chdir("../")

			output, err := test.GoRun("init", "-n", "-f", "--author", author, "--license", license, "--language", language)
			if err != nil {
				t.Errorf("process exited with err: %v", err)
			}
			t.Log("output:\n", output)

			t.Run("build docker image", func(t *testing.T) {
				output, err := test.Run("docker-compose", "build")

				if err != nil {
					t.Errorf("process exited with err: %v", err)
				}

				t.Log("output:\n", output)
			})

			t.Run("load settings", func(t *testing.T) {
				if output, err := test.Run("cp", "project-settings.toml.example", "project-settings.toml"); err != nil {
					t.Fatalf("failed to copy project-settings.toml: %s", output)
				}

				args := append([]string{"--log-level", "ERROR", "run", "--rm", "jupyter"}, settingsTests[language].Command...)
				args = append(args, settingsTests[language].Script)
				output, err := test.Run("docker-compose", args...)

				if err != nil {
					t.Errorf("process exited with err: %v", err)
				} else if output != settingsTests[language].Output {
					t.Errorf("expected: %s got: %s", settingsTests[language].Output, output)
				} else {
					t.Log("output:\n", output)
				}
			})
		})
	}
}

type settingsTest struct {
	Script, Output string
	Command        []string
}

var settingsTests = map[string]settingsTest{
	"python": settingsTest{
		Command: []string{"python", "-B", "-c"},
		Script: `from src import settings
print(settings.settings["downsample"])`,
		Output: "True\n",
	},
	"r": settingsTest{
		Command: []string{"Rscript", "-e"},
		Script: `source("/project/src/settings.R")
print(settings["downsample"])`,
		Output: "$downsample\n[1] TRUE\n\n",
	},
}
