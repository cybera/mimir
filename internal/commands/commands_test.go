package commands

import (
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/cybera/mimir/internal/paths"
	"github.com/spf13/viper"
)

func TestMain(m *testing.M) {
	viper.Set("ProjectRoot", "/my/project/root")
	os.Exit(m.Run())
}

func TestDockerCompose(t *testing.T) {
	composeFile := filepath.Join(viper.GetString("ProjectRoot"), paths.DockerCompose())

	want := exec.Command("docker-compose", "-f", composeFile, "up", "-d")
	got := DockerCompose("up", "-d")

	if !reflect.DeepEqual(got, want) {
		t.Errorf("\nGot : %v\nWant: %v", got, want)
	}
}

func TestScript(t *testing.T) {
	composeFile := filepath.Join(viper.GetString("ProjectRoot"), paths.DockerCompose())
	scriptsDir := filepath.Join(paths.ContainerRoot(), paths.Scripts())
	scriptFile := filepath.Join(scriptsDir, "script.sh")

	want := exec.Command("docker-compose", "-f", composeFile, "run", "--rm", "jupyter", scriptFile, "--flag")
	got := Script("script.sh", "--flag")

	if !reflect.DeepEqual(got, want) {
		t.Errorf("\nGot : %v\nWant: %v", got, want)
	}
}
