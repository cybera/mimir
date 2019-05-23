package paths

func ContainerRoot() string {
	return "/project/"
}

func DockerCompose() string {
	return "docker-compose.yml"
}

func Dockerfile() string {
	return "Dockerfile"
}

func Scripts() string {
	return "src/scripts/"
}

func DatasetsCode() string {
	return "src/datasets/"
}

func RawDatasets() string {
	return "data/raw/"
}

func ProcessedDatasets() string {
	return "data/processed/"
}

func Config() string {
	return ".ccds/config.yaml"
}
