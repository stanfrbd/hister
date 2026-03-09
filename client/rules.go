package client

import (
	"encoding/json"
	"net/url"
	"strings"
)

func (c *Client) FetchRules() (*RulesResponse, error) {
	req, err := c.newRequest("GET", "/api/rules", nil)
	if err != nil {
		return nil, err
	}
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer func() { _ = resp.Body.Close() }()
	if err := checkStatus(resp); err != nil {
		return nil, err
	}
	var data RulesResponse
	err = json.NewDecoder(resp.Body).Decode(&data)
	return &data, err
}

func (c *Client) SaveRules(skip, priority string) error {
	formData := url.Values{"skip": {skip}, "priority": {priority}}
	req, err := c.newRequest("POST", "/api/rules", strings.NewReader(formData.Encode()))
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

func (c *Client) AddAlias(keyword, value string) error {
	formData := url.Values{"alias-keyword": {keyword}, "alias-value": {value}}
	req, err := c.newRequest("POST", "/api/add_alias", strings.NewReader(formData.Encode()))
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

func (c *Client) DeleteAlias(alias string) error {
	formData := url.Values{"alias": {alias}}
	req, err := c.newRequest("POST", "/api/delete_alias", strings.NewReader(formData.Encode()))
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
