package extractor

import (
	"fmt"
	"strings"

	"github.com/internetarchive/Zeno/internal/pkg/utils"
	"github.com/internetarchive/Zeno/pkg/models"
)

// All the supported object storage servers
var ObjectStorageServers = func() (s []string) {
	s = append(s, s3CompatibleServers...)
	s = append(s, azureServers...)
	return s
}()

type ObjectStorageOutlinkExtractor struct{}

// Check if the response is from an object storage server
func (ObjectStorageOutlinkExtractor) Match(URL *models.URL) bool {
	return utils.StringContainsSliceElements(URL.GetResponse().Header.Get("Server"), ObjectStorageServers) &&
		URL.GetMIMEType() != nil &&
		strings.Contains(URL.GetMIMEType().String(), "/xml") // tricky match both application/xml and text/xml
}

// ObjectStorage decides which helper to call based on the object storage server
func (ObjectStorageOutlinkExtractor) Extract(URL *models.URL) ([]*models.URL, error) {
	defer URL.RewindBody()

	server := URL.GetResponse().Header.Get("Server")
	if utils.StringContainsSliceElements(server, s3CompatibleServers) {
		return s3Compatible(URL)
	} else if utils.StringContainsSliceElements(server, azureServers) {
		return azure(URL)
	} else {
		return nil, fmt.Errorf("unknown object storage server: %s", server)
	}
}
