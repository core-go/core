package client

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
)

type PrivilegesClient struct {
	Client *http.Client
	Url    string
}

func NewPrivilegesClient(client *http.Client, url string) *PrivilegesClient {
	return &PrivilegesClient{Client: client, Url: url}
}

func (c PrivilegesClient) Privileges(ctx context.Context, userId string) []string {
	url := c.Url + "/" + userId
	resp, err := c.Client.Get(url)
	privileges := make([]string, 0)
	if err != nil {
		return privileges
	}
	defer resp.Body.Close()
	body, er2 := ioutil.ReadAll(resp.Body)
	if er2 != nil {
		return privileges
	}
	var ps []string
	er3 := json.Unmarshal(body, &ps)
	if er3 != nil {
		return privileges
	}
	return ps
}
