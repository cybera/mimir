package paths

import (
	"errors"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
)

func ProjectRoot() string {
	root, err := ProjectRootSafe()
	if err != nil {
		log.Fatal(err)
	}

	return root
}

func ProjectRootSafe() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	for {
		rel, err := filepath.Rel("/", dir)
		if err != nil {
			log.Fatal(err)
		}
		if rel == "." {
			return "", errors.New("Not under a valid project directory")
		}

		files, err := ioutil.ReadDir(dir)
		if err != nil {
			log.Fatal(err)
		}

		for _, f := range files {
			if f.Name() == ".ccds" && f.IsDir() {
				return dir, nil
			}
		}

		dir = filepath.Join(dir, "../")
	}
}
