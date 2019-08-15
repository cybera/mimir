package fetchers

import (
	"io"
	"os"

	"github.com/pkg/errors"
)

type LocalFetcher struct {
	path string
}

func NewLocalFetcher(config FetcherConfig) (Fetcher, error) {
	return LocalFetcher{path: config.From}, nil
}

func (f LocalFetcher) Fetch(writer io.Writer) error {
	file, err := os.Open(f.path)
	if err != nil {
		return errors.Wrapf(err, "error opening %s", f.path)
	}

	_, err = io.Copy(writer, file)
	if err != nil {
		return errors.Wrapf(err, "error copying %s", f.path)
	}

	return nil
}

func init() {
	Factories["local"] = NewLocalFetcher
}
