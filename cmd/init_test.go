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

	for _, language = range languages.Supported {
		t.Run(language, func(t *testing.T) {
			testDir, err := test.CreateTestDir()
			if err != nil {
				t.Fatal(err)
			}
			defer os.RemoveAll(testDir)
			defer os.Chdir("../")

			output, err := test.RunCommand("init", "-n", "-f", "--author", author, "--license", license, "--language", language)
			if err != nil {
				t.Errorf("process exited with err: %v", err)
			}
			t.Log("output:\n", output)
		})
	}
}
