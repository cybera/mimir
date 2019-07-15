package fetchers

import (
	"github.com/pkg/errors"
)

var Factories = map[string]func(string, interface{}) (Fetcher, error){}

type Fetcher interface {
	Fetch() ([]byte, error)
}

func NewFetcher(name, from string, args interface{}) (Fetcher, error) {
	if factory, ok := Factories[name]; ok {
		return factory(from, args)
	}

	return nil, errors.New("no supported fetcher for " + name)
}
