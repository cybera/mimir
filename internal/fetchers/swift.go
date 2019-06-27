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

func NewSwiftFetcher(target string, args interface{}) (Fetcher, error) {
	split := strings.Split(target, "/")
	if len(split) < 2 {
		return nil, errors.New("no container specified")
	}

	return SwiftFetcher{containerName: split[0], objectName: split[1]}, nil
}

func (f SwiftFetcher) Fetch() ([]byte, error) {
	opts := &clientconfig.ClientOpts{}
	client, err := clientconfig.NewServiceClient("object-store", opts)
	if err != nil {
		return []byte{}, errors.Wrap(err, "swift client creation failed")
	}

	object := objects.Download(client, f.containerName, f.objectName, nil)
	content, err := object.ExtractContent()
	if err != nil {
		return []byte{}, errors.Wrap(err, "failed to extract object content")
	}

	return content, nil
}

func init() {
	factories["swift"] = NewSwiftFetcher
}
