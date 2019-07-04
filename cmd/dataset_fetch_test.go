package cmd

import (
	"os"
	"path/filepath"
	"strconv"
	"testing"

	"github.com/cybera/ccds/internal/datasets"
	"github.com/cybera/ccds/internal/fetchers"
	"github.com/cybera/ccds/internal/paths"
	"github.com/cybera/ccds/internal/test"
	"github.com/cybera/ccds/internal/utils"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
)

func TestFetch(t *testing.T) {
	fetchers.Factories["mock"] = NewMockFetcher

	testDir, err := test.CreateTestDir()
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(testDir)
	defer os.Chdir("../")

	for _, d := range []string{".ccds", paths.RawDatasets(), paths.DatasetsCode()} {
		if err := os.MkdirAll(d, os.ModePerm); err != nil {
			t.Fatalf("error creating directory: %v", err)
		}
	}

	utils.WriteConfig()

	source := datasets.Source{Name: "mock", Target: "container/iris.csv"}
	dataset := datasets.Dataset{File: "iris.csv", Source: source}

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
	source := datasets.Source{Name: "mock", Target: "container/iris.csv"}

	for i := 0; i < 5; i++ {
		name := strconv.Itoa(i)
		mockDatasets[name] = datasets.Dataset{File: name + ".csv", Source: source}
	}

	testDir, err := test.CreateTestDir()
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(testDir)
	defer os.Chdir("../")

	for _, d := range []string{".ccds", paths.RawDatasets(), paths.DatasetsCode()} {
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

func (m MockFetcher) Fetch() ([]byte, error) {
	if m.target != "container/iris.csv" {
		return nil, errors.New("target not found")
	}

	return []byte(contents), nil
}

func NewMockFetcher(target string, args interface{}) (fetchers.Fetcher, error) {
	return MockFetcher{target: target}, nil
}
