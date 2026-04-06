package cluster

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
)

// APICallsBootstrapStatus is used by cli command to interact with the cluster
func APICallsBootstrapStatus(config ClusterHTTPCallBaseConfig) error {
	path := "/api/v1/cluster/bootstrap/status"
	url := fmt.Sprintf("%s%s", config.HTTPAddress, path)
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer func() {
		_ = resp.Body.Close()
	}()
	body, _ := io.ReadAll(resp.Body)

	if resp.StatusCode == 200 {
		fmt.Println(string(body))
		return nil
	}

	return decodeError(body)
}

// APICallsBootstrapCluster is used by cli command
// to bootstrap the cluster
func APICallsBootstrapCluster(config BootstrapClusterHTTPConfig) error {
	path := "/api/v1/cluster/bootstrap/cluster"
	url := fmt.Sprintf("%s%s", config.HTTPAddress, path)

	b, _ := json.Marshal(APIBootstrapClusterRequest{Token: config.Token})
	reqBody := bytes.NewBuffer(b)
	req, _ := http.NewRequest("POST", url, reqBody)
	req.Header.Add("Content-Length", strconv.Itoa(reqBody.Len()))
	req.Header.Add("Content-Type", "application/json; charset=utf-8")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer func() {
		_ = resp.Body.Close()
	}()
	body, _ := io.ReadAll(resp.Body)

	if resp.StatusCode == 200 {
		fmt.Println(string(body))
		return nil
	}

	return decodeError(body)
}

// APICallsNodesList is used by cli command to list cluster nodes
func APICallsNodesList(config NodesListHTTPConfig) error {
	path := "/api/v1/cluster/nodes/list"
	url := fmt.Sprintf("%s%s", config.HTTPAddress, path)

	b, _ := json.Marshal(APINodesListRequest{Kind: config.Kind})
	reqBody := bytes.NewBuffer(b)
	req, _ := http.NewRequest("GET", url, reqBody)
	req.Header.Add("Content-Length", strconv.Itoa(reqBody.Len()))
	req.Header.Add("Content-Type", "application/json; charset=utf-8")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer func() {
		_ = resp.Body.Close()
	}()
	body, _ := io.ReadAll(resp.Body)

	if resp.StatusCode == 200 {
		if config.OutputFormat == "json" {
			fmt.Println(string(body))
			return nil
		}
		printTableNodesList(body)
		return nil
	}

	return decodeError(body)
}

// APICallsDeploymentApply is used by cli command to apply a new deployment
func APICallsDeploymentApply(config DeploymentApplyHTTPConfig) error {
	if err := parseDeploymentFileSyntax(config.File); err != nil {
		return err
	}

	var (
		data []byte
		err  error
	)
	if config.osReadFile == nil {
		data, err = os.ReadFile(config.File)
	} else {
		data, err = config.osReadFile(config.File)
	}

	if err != nil {
		return err
	}

	path := "/api/v1/deployment/apply"
	url := fmt.Sprintf("%s%s", config.HTTPAddress, path)

	b, _ := json.Marshal(APIDeploymentApplyRequest{HCLContent: string(data)})
	reqBody := bytes.NewBuffer(b)
	req, _ := http.NewRequest("POST", url, reqBody)
	req.Header.Add("Content-Length", strconv.Itoa(reqBody.Len()))
	req.Header.Add("Content-Type", "application/json; charset=utf-8")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer func() {
		_ = resp.Body.Close()
	}()
	body, _ := io.ReadAll(resp.Body)

	if resp.StatusCode == 200 {
		if config.OutputFormat == "json" {
			fmt.Println(string(body))
			return nil
		}
		return nil
	}

	return decodeError(body)
}

// APICallsPodsList is used by cli command to list pods
func APICallsPodsList(config PodsListHTTPConfig) error {
	path := "/api/v1/pods/list"
	url := fmt.Sprintf("%s%s", config.HTTPAddress, path)

	b, _ := json.Marshal(APIPodsListRequest{Namespace: config.Namespace})
	reqBody := bytes.NewBuffer(b)
	req, _ := http.NewRequest("GET", url, reqBody)
	req.Header.Add("Content-Length", strconv.Itoa(reqBody.Len()))
	req.Header.Add("Content-Type", "application/json; charset=utf-8")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer func() {
		_ = resp.Body.Close()
	}()
	body, _ := io.ReadAll(resp.Body)

	if resp.StatusCode == 200 || resp.StatusCode == 404 {
		if config.OutputFormat == "json" {
			fmt.Println(string(body))
			return nil
		}
		printTablePodsList(body)
		return nil
	}

	return decodeError(body)
}

// APICallsPodsDelete is used by cli command to delete pods
func APICallsPodsDelete(config PodsDeleteHTTPConfig) error {
	path := "/api/v1/pods/delete"
	url := fmt.Sprintf("%s%s", config.HTTPAddress, path)

	b, _ := json.Marshal(APIPodsDeleteRequest{
		Namespace: config.Namespace,
		Pods:      config.Pods,
		Detached:  config.Detached,
	})

	reqBody := bytes.NewBuffer(b)
	req, _ := http.NewRequest("DELETE", url, reqBody)
	req.Header.Add("Content-Length", strconv.Itoa(reqBody.Len()))
	req.Header.Add("Content-Type", "application/json; charset=utf-8")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer func() {
		_ = resp.Body.Close()
	}()
	body, _ := io.ReadAll(resp.Body)

	if resp.StatusCode == 200 {
		fmt.Println(string(body))
		return nil
	}

	return decodeError(body)
}
