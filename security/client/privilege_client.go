package client

import (
	"context"
	"io"
	"net/http"
	"strconv"
)

type PrivilegeClient struct {
	Client *http.Client
	Url    string
}

func NewPrivilegeClient(client *http.Client, url string) *PrivilegeClient {
	return &PrivilegeClient{Client: client, Url: url}
}

func (c *PrivilegeClient) Privilege(ctx context.Context, userId string, privilegeId string) int32 {
	url := c.Url + "/" + userId + "/" + privilegeId
	resp, er1 := c.Client.Get(url)
	if er1 != nil {
		return 0
	}
	defer resp.Body.Close()
	body, er2 := io.ReadAll(resp.Body)
	if er2 != nil {
		return 0
	}
	s := string(body)
	i64, er3 := strconv.ParseInt(s, 10, 32)
	if er3 != nil {
		return 0
	}
	i := int32(i64)
	return i
}
