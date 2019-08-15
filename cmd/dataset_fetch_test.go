package cmd

import (
	"io"
	"os"
	"path/filepath"
	"strconv"
	"testing"

	"github.com/cybera/mimir/internal/datasets"
	"github.com/cybera/mimir/internal/fetchers"
	"github.com/cybera/mimir/internal/paths"
	"github.com/cybera/mimir/internal/test"
	"github.com/cybera/mimir/internal/utils"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
)

func TestFetch(t *testing.T) {
	fetchers.Factories["mock"] = NewMockFetcher

	err := test.CreateTestDir("_test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll("_test")
	defer os.Chdir("../")

	for _, d := range []string{".mimir", paths.RawDatasets(), paths.DatasetsCode()} {
		if err := os.MkdirAll(d, os.ModePerm); err != nil {
			t.Fatalf("error creating directory: %v", err)
		}
	}

	utils.WriteConfig()

	fetcherConfig := fetchers.FetcherConfig{Name: "mock", From: "container/iris.csv"}
	dataset := datasets.Dataset{File: "iris.csv", FetcherConfig: fetcherConfig}

	if err := fetch(dataset); err != nil {
		t.Fatalf("failed to fetch dataset: %v", err)
	}

	equal, err := test.FileContentsEquals(filepath.Join(paths.RawDatasets(), "iris.csv"), contents)
	if err != nil {
		t.Fatal(err)
	}

	if !equal {
		t.Errorf("found unexpected file contents")
	}
}

func TestFetchAll(t *testing.T) {
	fetchers.Factories["mock"] = NewMockFetcher
	mockDatasets := map[string]datasets.Dataset{}
	fetcherConfig := fetchers.FetcherConfig{Name: "mock", From: "container/iris.csv"}

	for i := 0; i < 5; i++ {
		name := strconv.Itoa(i)
		mockDatasets[name] = datasets.Dataset{File: name + ".csv", FetcherConfig: fetcherConfig}
	}

	err := test.CreateTestDir("_test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll("_test")
	defer os.Chdir("../")

	for _, d := range []string{".mimir", paths.RawDatasets(), paths.DatasetsCode()} {
		if err := os.MkdirAll(d, os.ModePerm); err != nil {
			t.Fatalf("error creating directory: %v", err)
		}
	}

	viper.Set("Datasets", mockDatasets)
	utils.WriteConfig()

	if err := fetchAll(); err != nil {
		t.Fatalf("error fetching all datasets: %v", err)
	}

	for _, d := range mockDatasets {
		path := filepath.Join(paths.RawDatasets(), d.File)

		equal, err := test.FileContentsEquals(path, contents)
		if err != nil {
			t.Error(err)
		}

		if !equal {
			t.Errorf("found unexpected file contents for: %s", d.File)
		}
	}
}

const contents = `sepal_length,sepal_width,petal_length,petal_width,species
5.1,3.5,1.4,0.2,setosa
4.9,3,1.4,0.2,setosa
4.7,3.2,1.3,0.2,setosa
4.6,3.1,1.5,0.2,setosa
5,3.6,1.4,0.2,setosa`

type MockFetcher struct {
	target string
}

func (m MockFetcher) Fetch(writer io.Writer) error {
	if m.target != "container/iris.csv" {
		return errors.New("target not found")
	}

	_, err := writer.Write([]byte(contents))

	return err
}

func NewMockFetcher(config fetchers.FetcherConfig) (fetchers.Fetcher, error) {
	return MockFetcher{target: config.From}, nil
}
