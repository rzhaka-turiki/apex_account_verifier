package api

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/rzhaka-turiki/apex_account_verifier/internal/model"
)

type Client struct {
	baseURL   string
	authToken string
	http      *http.Client
}

func NewClient(baseURL, authToken string) *Client {
	return &Client{
		baseURL:   baseURL,
		authToken: authToken,
		http:      &http.Client{},
	}
}

type bridgeResponse struct {
	Global struct {
		UID      string `json:"uid"`
		Name     string `json:"name"`
		Platform string `json:"platform"`
		Level    int    `json:"level"`
	} `json:"global"`
}

func (c *Client) GetAccount(ctx context.Context, player, platform string) (*model.Account, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.baseURL, nil)
	if err != nil {
		return nil, err
	}
	q := req.URL.Query()
	q.Set("auth", c.authToken)
	q.Set("player", player)
	q.Set("platform", platform)

	req.URL.RawQuery = q.Encode()

	resp, err := c.http.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("apex api returned %d", resp.StatusCode)
	}
	var data bridgeResponse
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, err
	}
	if data.Global.UID == "" {
		return nil, fmt.Errorf("account uid not found")
	}
	return &model.Account{
		UID:      data.Global.UID,
		Player:   data.Global.Name,
		Platform: data.Global.Platform,
		Level:    data.Global.Level,
	}, nil
}
