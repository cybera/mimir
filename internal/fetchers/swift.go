package fetchers

import (
	"strings"

	"github.com/gophercloud/gophercloud/openstack/objectstorage/v1/objects"
	"github.com/gophercloud/utils/openstack/clientconfig"
	"github.com/pkg/errors"
)

type SwiftFetcher struct {
	containerName string
	objectName    string
}

func NewSwiftFetcher(config FetcherConfig) (Fetcher, error) {
	containerEnd := strings.LastIndex(config.From, "/")
	if containerEnd == -1 {
		return nil, errors.New("no container specified")
	}

	container := config.From[0:containerEnd]
	object := config.From[containerEnd+1 : len(config.From)]

	return SwiftFetcher{containerName: container, objectName: object}, nil
}

func (f SwiftFetcher) Fetch() ([]byte, error) {
	client, err := clientconfig.NewServiceClient("object-store", nil)
	if err != nil {
		return nil, errors.Wrap(err, "swift client creation failed")
	}

	object := objects.Download(client, f.containerName, f.objectName, nil)
	content, err := object.ExtractContent()
	if err != nil {
		return nil, errors.Wrap(err, "failed to extract object content")
	}

	return content, nil
}

func init() {
	Factories["swift"] = NewSwiftFetcher
}
