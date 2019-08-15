package fetchers

import (
	"io"
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

func (f SwiftFetcher) Fetch(writer io.Writer) error {
	client, err := clientconfig.NewServiceClient("object-store", nil)
	if err != nil {
		return errors.Wrap(err, "swift client creation failed")
	}

	res := objects.Download(client, f.containerName, f.objectName, nil)
	if res.Err != nil {
		return errors.Wrapf(res.Err, "error getting object %s/%s", f.containerName, f.objectName)
	}

	_, err = io.Copy(writer, res.Body)
	if err != nil {
		return errors.Wrapf(err, "error writing object %s/%s", f.containerName, f.objectName)
	}

	return nil
}

func init() {
	Factories["swift"] = NewSwiftFetcher
}
