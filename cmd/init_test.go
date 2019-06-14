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
		})
	}
}
