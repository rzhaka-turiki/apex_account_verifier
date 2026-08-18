package api

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/rzhaka-turiki/apex_account_verifier/internal/dto"
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
		http: &http.Client{
			Transport: &http.Transport{
				TLSHandshakeTimeout: 10 * time.Second,
			},
			Timeout: 30 * time.Second,
		},
	}
}

type APIError struct {
	StatusCode int
}

func (e *APIError) Error() string {
	return fmt.Sprintf("apex api returned %d", e.StatusCode)
}

func (c *Client) GetAccount(ctx context.Context, player, platform string) (*model.Account, error) {
	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		c.baseURL,
		nil,
	)
	if err != nil {
		return nil, err
	}

	q := req.URL.Query()
	q.Set("auth", c.authToken)
	q.Set("player", player)
	q.Set("platform", platform)

	req.URL.RawQuery = q.Encode()

	fmt.Printf(
		"REQUEST URL: %s?auth=***&player=%s&platform=%s\n",
		c.baseURL,
		player,
		platform,
	)

	resp, err := c.http.Do(req)
	if err != nil {
		fmt.Printf("HTTP REQUEST ERROR: %v\n", err)
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		fmt.Printf("APEX STATUS: %d\n", resp.StatusCode)
		fmt.Printf("APEX HEADERS: %v\n", resp.Header)
		fmt.Printf("APEX BODY: %s\n", string(body))

		return nil, &APIError{
			StatusCode: resp.StatusCode,
		}
	}

	var data dto.BridgeResponse

	if err := json.Unmarshal(body, &data); err != nil {
		return nil, err
	}

	if data.Global.UID == "" {
		return nil, fmt.Errorf("account uid not found")
	}

	const levelsPerPrestige = 500

	totalLevel := data.Global.Level +
		data.Global.LevelPrestige*levelsPerPrestige

	return &model.Account{
		UID:      data.Global.UID,
		Player:   data.Global.Name,
		Platform: data.Global.Platform,
		Level:    totalLevel,
	}, nil
}
