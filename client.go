package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type (
	Client struct {
		Client *http.Client
		Host string
	}

	RequestParamPostShellgei struct {
		Code string `json:"code"`
		Images []string `json:"images"`
	}
	ResponseParamPostShellgei struct {
		Stdout string `json:"stdout"`
		Stderr string `json:"stderr"`
	}
)

const (
	WebshHost = "https://websh.jiro4989.com/api"
)

func NewClient(host string) *Client {
	c := &http.Client{
		Timeout: 10 * time.Second,
	}
	return &Client{
		Client: c,
		Host: host,
	}
}

func (c *Client) PostShellgei(p *RequestParamPostShellgei) (*ResponseParamPostShellgei, error) {
	jsonBody, err := json.Marshal(p)
	if err != nil {
		return nil, err
	}

	body := bytes.NewBuffer(jsonBody)
	resp,err := c.Client.Post(c.url("/shellgei"), "application/json", body)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	st := resp.StatusCode
	if 400 <= st && st < 600 {
		return nil, fmt.Errorf("error occured: status = %d", st)
	}

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var re ResponseParamPostShellgei
	err = json.Unmarshal(b, &re)
	if err != nil {
		return nil, err
	}

	return &re, nil
}

func (c *Client) url(path string) string {
	return fmt.Sprintf("%s%s", c.Host, path)
}
