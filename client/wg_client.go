package wgc

import (
	"encoding/base64"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"
)

type Client struct {
	apiAddress  string
	apiUsername string
	apiPassword string
}

func NewClient(apiAddress, apiUsername, apiPassword string) *Client {
	wclient := &Client{apiAddress, apiUsername, apiPassword}
	return wclient
}

func basicAuth(username, password string) string {
	creds := username + ":" + password
	return base64.StdEncoding.EncodeToString([]byte(creds))
}

func (c *Client) Get(path string) ([]byte, error) {
	var client http.Client

	fullPath := c.apiAddress + path
	if !strings.Contains(path, "/wanguard-api/v1/") {
		fullPath = c.apiAddress + "/wanguard-api/v1/" + path
	}

	req, err := http.NewRequest("GET", fullPath, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Authorization", "Basic "+basicAuth(c.apiUsername, c.apiPassword))

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return body, nil
}

func (c *Client) GetParsed(path string, obj interface{}) error {
	body, err := c.Get(path)
	if err != nil {
		return err
	}

	err = json.Unmarshal(body, obj)
	return err
}
