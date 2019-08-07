package fetchers

import (
	"io/ioutil"

	"github.com/pkg/errors"
)

type LocalFetcher struct {
	path string
}

func NewLocalFetcher(config FetcherConfig) (Fetcher, error) {
	return LocalFetcher{path: config.From}, nil
}

func (f LocalFetcher) Fetch() ([]byte, error) {
	data, err := ioutil.ReadFile(f.path)
	if err != nil {
		return nil, errors.Wrap(err, "failed to fetch local file")
	}

	return data, nil
}

func init() {
	Factories["local"] = NewLocalFetcher
}
