package cmd

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"
	"strings"
	"testing"

	"github.com/cybera/mimir/internal/test"
)

func TestJupyter(t *testing.T) {
	os.RemoveAll("_test")

	if output, err := test.InitProject("_test", "John Doe", "MIT", "python"); err != nil {
		t.Log(output)
		t.Fatalf("failed to create test directory: %v", err)
	}
	defer os.RemoveAll("_test")
	defer os.Chdir("../")

	output, err := test.GoRun("jupyter", "start")
	if err != nil {
		t.Fatalf("process exited with error: %v", err)
	}
	defer test.GoRun("jupyter", "stop")

	regex, err := regexp.Compile(`\?token=[A-Za-z0-9]*`)
	if err != nil {
		t.Fatalf("error compiling regex: %v", err)
	}

	match := regex.FindString(output)
	if match == "" {
		t.Fatalf("error extracting jupyter auth token")
	}

	url := "http://localhost:8888/api/contents/docker-compose.yml" + match

	resp, err := http.Get(url)
	if err != nil {
		t.Fatalf("error getting file contents from jupyter: %v", err)
	}
	defer resp.Body.Close()

	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("error reading response body: %v", err)
	}

	var body map[string]interface{}

	if err := json.Unmarshal(bytes, &body); err != nil {
		t.Fatalf("error unmarshaling response to json: %v", err)
	}

	content, ok := body["content"]
	if !ok {
		t.Fatal("unexpected response: ", body)
	}

	if !strings.HasPrefix(content.(string), "version:") {
		t.Errorf("unexpected file contents:\n%s", content.(string))
	}
}
