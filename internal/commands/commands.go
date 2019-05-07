package commands

import (
	"os/exec"
	"path/filepath"

	"github.com/cybera/ccds/internal/paths"
	"github.com/spf13/viper"
)

// DockerCompose generates a docker-compose command
func DockerCompose(subCommand string, args ...string) *exec.Cmd {
	projectRoot := viper.GetString("ProjectRoot")
	args = append([]string{"-f", paths.DockerCompose(projectRoot), subCommand}, args...)
	return exec.Command("docker-compose", args...)
}

// Script generates a docker-compose command to run the specified script
func Script(name string, args ...string) *exec.Cmd {
	containerRoot := viper.GetString("ContainerRoot")
	path := filepath.Join(paths.Scripts(containerRoot), name)
	args = append([]string{"--rm", "jupyter", path}, args...)
	return DockerCompose("run", args...)
}
