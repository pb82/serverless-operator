package stub

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"github.com/sirupsen/logrus"
)

// Openwhisk Actions Base URL
const BASE_URL = "http://%s/api/v1/namespaces/%s/actions/%s%s"

type OpenwhiskClient struct {
	clusterIp string
	username  string
	password  string
}

// Create Action Payload
type OpenwhiskPayload struct {
	Namespace string            `json:"namespace"`
	Name      string            `json:"name"`
	Exec      map[string]string `json:"exec,omitempty"`
}

func (client *OpenwhiskClient) createAction(payload *OpenwhiskPayload) error {
	url := fmt.Sprintf(BASE_URL, client.clusterIp, payload.Namespace, payload.Name, "?overwrite=true")

	action, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	req, err := http.NewRequest(http.MethodPut, url, bytes.NewBuffer(action))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.SetBasicAuth(client.username, client.password)

	httpClient := http.Client{}
	resp, err := httpClient.Do(req)
	if err != nil {
		return err
	}

	success := resp.StatusCode == http.StatusOK
	if success {
		return nil
	}

	return errors.New(fmt.Sprintf("Error creating action, status code %d", resp.StatusCode))
}

func (client *OpenwhiskClient) deleteAction(name string, namespace string) error {
	url := fmt.Sprintf(BASE_URL, client.clusterIp, namespace, name, "")
	req, err := http.NewRequest(http.MethodDelete, url, nil)
	if err != nil {
		logrus.Error(err.Error())
		return err
	}

	req.SetBasicAuth(client.username, client.password)
	httpClient := http.Client{}
	resp, err := httpClient.Do(req)
	if err != nil {
		return err
	}

	success := resp.StatusCode == http.StatusOK
	if success {
		return nil
	}

	return errors.New(fmt.Sprintf("Error deleting action, status code %d", resp.StatusCode))
}