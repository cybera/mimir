package commands

import (
	"os/exec"

	"github.com/cybera/ccds/internal/paths"
)

func DockerCompose(args ...string) *exec.Cmd {
	args = append([]string{"-f", paths.DockerCompose()}, args...)
	return exec.Command("docker-compose", args...)
}
