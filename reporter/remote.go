package reporter

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

type (
	remoteReporter struct {
		endpoint string

		nodeID    string
		hostname  string
		privateIP string
		publicIP  string
	}
)

func NewRemoteReporter(endpoint string, metadataFetcher MetadataFetcher) Reporter {

	nodeID, err := metadataFetcher.Fetch("id")
	if err != nil {
		return nil
	}

	hostName, err := metadataFetcher.Fetch("hostname")
	if err != nil {
		return nil
	}

	privateIP, err := metadataFetcher.Fetch("interfaces/private/0/ipv4/address")
	if err != nil {
		return nil
	}

	publicIP, err := metadataFetcher.Fetch("interfaces/public/0/ipv4/address")
	if err != nil {
		return nil
	}

	fmt.Printf("[%s] [%s] [%s] [%s]", nodeID, hostName, publicIP, privateIP)

	return &remoteReporter{
		endpoint:  endpoint,
		nodeID:    nodeID,
		hostname:  hostName,
		privateIP: privateIP,
		publicIP:  publicIP,
	}
}

func (rr *remoteReporter) Report(payload []byte) error {

	var jsonStr = []byte(`{"ID":"asg-1","NodeID":"` + rr.nodeID + `", "Metrics":[{"Value":1.0, "Time":"` + time.Now().Format(time.RFC3339Nano) + `"}]}`)
	req, err := http.NewRequest("POST", rr.endpoint, bytes.NewBuffer(jsonStr))
	if err != nil {
		return err
	}
	req.Header.Set("X-API-Key", "every-node-has-it-own")
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	fmt.Println("response Status:", resp.Status)
	fmt.Println("response Headers:", resp.Header)
	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println("response Body:", string(body))

	return nil
}
