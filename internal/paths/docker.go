package paths

import "path/filepath"

func DockerCompose() string {
	return filepath.Join(ProjectRoot(), "docker-compose.yml")
}

func Dockerfile() string {
	return filepath.Join(ProjectRoot(), "Dockerfile")
}
