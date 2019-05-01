package paths

import (
	"log"
	"path/filepath"
)

func ContainerPath(localPath string) string {
	root := ProjectRoot()
	rel, err := filepath.Rel(root, localPath)
	if err != nil {
		log.Fatal(err)
	}

	return filepath.Join("/project", rel)
}
