package cluster

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
)

// APICallsBootstrapStatus is used by cli command to interact
// with the cluster
func APICallsBootstrapStatus(config ClusterHTTPCallBaseConfig) {
	path := "/api/v1/cluster/bootstrap/status"
	url := fmt.Sprintf("%s%s", config.HTTPAddress, path)
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	defer func() {
		_ = resp.Body.Close()
	}()
	body, _ := io.ReadAll(resp.Body)

	fmt.Println(string(body))
}

// APICallsBootstrapCluster is used by cli command
// to bootstrap the cluster
func APICallsBootstrapCluster(config BootstrapClusterHTTPConfig) {
	path := "/api/v1/cluster/bootstrap/cluster"
	url := fmt.Sprintf("%s%s", config.HTTPAddress, path)

	b, _ := json.Marshal(APIBootstrapClusterRequest{Token: config.Token})
	reqBody := bytes.NewBuffer(b)
	req, _ := http.NewRequest("POST", url, reqBody)
	req.Header.Add("Content-Length", strconv.Itoa(reqBody.Len()))
	req.Header.Add("Content-Type", "application/json; charset=utf-8")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	defer func() {
		_ = resp.Body.Close()
	}()
	body, _ := io.ReadAll(resp.Body)

	fmt.Println(string(body))
}

// APICallsNodesList is used by cli command
// to list cluster nodes
func APICallsNodesList(config NodesListHTTPConfig) {
	path := "/api/v1/cluster/nodes/list"
	url := fmt.Sprintf("%s%s", config.HTTPAddress, path)

	b, _ := json.Marshal(APINodesListRequest{Kind: config.Kind})
	reqBody := bytes.NewBuffer(b)
	req, _ := http.NewRequest("GET", url, reqBody)
	req.Header.Add("Content-Length", strconv.Itoa(reqBody.Len()))
	req.Header.Add("Content-Type", "application/json; charset=utf-8")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	defer func() {
		_ = resp.Body.Close()
	}()
	body, _ := io.ReadAll(resp.Body)

	if config.OutputFormat == "json" {
		fmt.Println(string(body))
		return
	}
	printTableNodesList(body)
}
