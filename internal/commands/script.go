package commands

import (
	"os/exec"
	"path/filepath"

	"github.com/cybera/ccds/internal/paths"
)

func Script(name string, args ...string) *exec.Cmd {
	path := filepath.Join(paths.Scripts(), name)
	return exec.Command(path, args...)
}
