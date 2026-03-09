package client

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/url"
	"strings"

	"github.com/asciimoo/hister/server/indexer"
)

func (c *Client) AddDocumentJSON(doc *indexer.Document) error {
	data, err := json.Marshal(doc)
	if err != nil {
		return err
	}
	req, err := c.newRequest("POST", "/api/add", bytes.NewReader(data))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer func() { _ = resp.Body.Close() }()
	return checkStatus(resp)
}

func (c *Client) AddPage(u, title, text string) error {
	formData := url.Values{"url": {u}, "title": {title}, "text": {text}}
	req, err := c.newRequest("POST", "/api/add", strings.NewReader(formData.Encode()))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer func() { _ = resp.Body.Close() }()
	return checkStatus(resp)
}

func (c *Client) DocumentExists(u string) (bool, error) {
	req, err := c.newRequest("HEAD", "/api/document?url="+url.QueryEscape(u), nil)
	if err != nil {
		return false, err
	}
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return false, err
	}
	_ = resp.Body.Close()
	return resp.StatusCode == http.StatusOK, nil
}

func (c *Client) DeleteDocument(u string) error {
	formData := url.Values{"url": {u}}
	req, err := c.newRequest("POST", "/api/delete", strings.NewReader(formData.Encode()))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer func() { _ = resp.Body.Close() }()
	return checkStatus(resp)
}
