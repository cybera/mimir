package commands

import (
	"os/exec"
	"path/filepath"

	"github.com/cybera/mimir/internal/paths"
	"github.com/spf13/viper"
)

// DockerCompose generates a docker-compose command
func DockerCompose(subCommand string, args ...string) *exec.Cmd {
	projectRoot := viper.GetString("ProjectRoot")
	args = append([]string{"-f", filepath.Join(projectRoot, paths.DockerCompose()), subCommand}, args...)
	return exec.Command("docker-compose", args...)
}

// Script generates a docker-compose command to run the specified script
func Script(name string, args ...string) *exec.Cmd {
	path := filepath.Join(paths.ContainerRoot(), paths.Scripts(), name)
	args = append([]string{"--rm", "jupyter", path}, args...)
	return DockerCompose("run", args...)
}
