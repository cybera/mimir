package fetchers

import (
	"github.com/pkg/errors"
)

var Factories = map[string]func(string, interface{}) (Fetcher, error){}

type Fetcher interface {
	Fetch() ([]byte, error)
}

func NewFetcher(source, target string, args interface{}) (Fetcher, error) {
	if factory, ok := Factories[source]; ok {
		return factory(target, args)
	}

	return nil, errors.New("no supported fetcher for " + source)
}
