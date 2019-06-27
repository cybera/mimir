package paths

import (
	"errors"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
)

func ProjectRoot() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
		return "", err
	}

	for {
		rel, err := filepath.Rel(FSRoot, dir)
		if err != nil {
			return "", err
		}
		if rel == "." {
			return "", errors.New("not under a valid project directory")
		}

		files, err := ioutil.ReadDir(dir)
		if err != nil {
			return "", err
		}

		for _, f := range files {
			if f.Name() == ".ccds" && f.IsDir() {
				return dir, nil
			}
		}

		dir = filepath.Join(dir, "../")
	}
}
