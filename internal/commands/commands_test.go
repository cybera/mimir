package commands

import (
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/cybera/ccds/internal/paths"
	"github.com/spf13/viper"
)

func TestMain(m *testing.M) {
	viper.Set("ProjectRoot", "/my/project/root")
	viper.Set("ContainerRoot", "/project")
	os.Exit(m.Run())
}

func TestDockerCompose(t *testing.T) {
	projectRoot := viper.GetString("ProjectRoot")
	composeFile := paths.DockerCompose(projectRoot)

	want := exec.Command("docker-compose", "-f", composeFile, "up", "-d")
	got := DockerCompose("up", "-d")

	if !reflect.DeepEqual(got, want) {
		t.Errorf("\nGot : %v\nWant: %v", got, want)
	}
}

func TestScript(t *testing.T) {
	projectRoot := viper.GetString("ProjectRoot")
	containerRoot := viper.GetString("ContainerRoot")
	composeFile := paths.DockerCompose(projectRoot)
	scriptsDir := paths.Scripts(containerRoot)
	scriptFile := filepath.Join(scriptsDir, "script.sh")

	want := exec.Command("docker-compose", "-f", composeFile, "run", "--rm", "jupyter", scriptFile, "--flag")
	got := Script("script.sh", "--flag")

	if !reflect.DeepEqual(got, want) {
		t.Errorf("\nGot : %v\nWant: %v", got, want)
	}
}
