package fetchers

import (
	"io"

	"github.com/pkg/errors"
)

var Factories = map[string]func(FetcherConfig) (Fetcher, error){}

type Fetcher interface {
	Fetch(io.Writer) error
}

type FetcherConfig struct {
	Name string
	From string
	Args interface{}
}

func NewFetcher(config FetcherConfig) (Fetcher, error) {
	if factory, ok := Factories[config.Name]; ok {
		return factory(config)
	}

	return nil, errors.New("no supported fetcher for " + config.Name)
}
