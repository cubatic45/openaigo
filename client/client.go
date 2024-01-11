package client

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"time"
)

var defaultClient = &http.Client{Timeout: 5 * time.Minute}

type Client struct {
	key string
	*http.Client
}

func NewClient(key string) *Client {
	return &Client{
		key:    key,
		Client: defaultClient,
	}
}

func NewClientWithClient(key string, client *http.Client) *Client {
	return &Client{
		key:    key,
		Client: client,
	}
}

func (c *Client) newRequest(method, url string, body any) (*http.Request, error) {
	var bodyReader io.Reader
	switch v := body.(type) {
	case string:
		bodyReader = strings.NewReader(v)
	case []byte:
		bodyReader = bytes.NewReader(v)
	case io.Reader:
		bodyReader = v
	}
	req, err := http.NewRequest(method, url, bodyReader)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.key)
	return req, nil
}

func (c *Client) sendRequest(req *http.Request, response any) error {
	resp, err := c.Client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 400 {
		// TODO: handle error
	}
	return json.NewDecoder(resp.Body).Decode(response)
}
