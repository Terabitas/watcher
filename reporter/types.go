package reporter

import (
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"strconv"
	"time"
)

type (

	// Reporter type is responsible for reporting status of node
	Reporter interface {
		Report(payload []byte) error
	}

	// MetadataFetcher type
	MetadataFetcher interface {
		Fetch(path string) (string, error)
	}

	// DigitalOceanMetadataFetcher to use in DO
	DigitalOceanMetadataFetcher struct {
		apiEndpoint string
	}

	// LocalMetadataFetcher type for local use
	LocalMetadataFetcher struct {
		data map[string]string
	}
)

// NewDigitalOceanMetadataFetcher constructor
func NewDigitalOceanMetadataFetcher() MetadataFetcher {
	return &DigitalOceanMetadataFetcher{
		apiEndpoint: "http://169.254.169.254/metadata/v1/",
	}
}

// NewLocalMetadataFetcher constructor
func NewLocalMetadataFetcher() MetadataFetcher {
	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))
	id := strconv.FormatInt(rnd.Int63n(999999), 10)
	ip := strconv.FormatInt(rnd.Int63n(253), 10)
	d := map[string]string{
		"id":                                id,
		"hostname":                          fmt.Sprintf("drople-%s", id),
		"interfaces/private/0/ipv4/address": fmt.Sprintf("1.1.1.%s", ip),
		"interfaces/public/0/ipv4/address":  fmt.Sprintf("1.1.2.%s", ip),
	}

	return &LocalMetadataFetcher{
		data: d,
	}
}

// Fetch impl
func (lmf *LocalMetadataFetcher) Fetch(path string) (string, error) {
	return lmf.data[path], nil
}

// Fetch impl
func (lmf *DigitalOceanMetadataFetcher) Fetch(path string) (string, error) {
	response, err := http.Get(lmf.apiEndpoint + path)
	if err != nil {
		return "", err
	}

	defer response.Body.Close()
	contents, err := ioutil.ReadAll(response.Body)

	if err != nil {
		return "", err
	}

	return string(contents), nil
}
