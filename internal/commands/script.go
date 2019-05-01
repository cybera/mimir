package commands

import (
	"os/exec"
	"path/filepath"

	"github.com/cybera/ccds/internal/paths"
)

func Script(name string, args ...string) *exec.Cmd {
	path := paths.ContainerPath(filepath.Join(paths.Scripts(), name))
	args = append([]string{"run", "--rm", "jupyter", path}, args...)
	return DockerCompose(args...)
}
