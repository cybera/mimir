package cmd

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/cybera/ccds/internal/languages"
)

func TestInit(t *testing.T) {
	author := "John Doe"
	license := "MIT"

	for _, language = range languages.Supported {
		t.Run(language, func(t *testing.T) {
			testDir, _ := filepath.Abs("_test")
			if err := os.Mkdir(testDir, os.ModePerm); err != nil {
				t.Fatalf("error creating test directory: %v", err)
			}
			defer os.RemoveAll(testDir)
			if err := os.Chdir(testDir); err != nil {
				t.Fatalf("error changing to test directory: %v", err)
			}
			defer os.Chdir("../")

			cmd := exec.Command("go", "run", "../../main.go", "init", "-n", "-f", "--author", author, "--license", license, "--language", language)
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr

			if err := cmd.Run(); err != nil {
				t.Fatalf("process exited with err: %v", err)
			}
		})
	}
}
