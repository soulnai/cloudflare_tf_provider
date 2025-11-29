package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
	// ... other imports
)

type Client struct {
	HTTPClient *http.Client
	ApiToken   string
	AccountId  string
	BaseURL    string
}

type Tunnel struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	Secret string `json:"tunnel_secret,omitempty"` //Sets the password required to run a locally-managed tunnel. Must be at least 32 bytes and encoded as a base64 string.
	// Add other fields if needed, but these are the core ones
}

type TunnelResponse struct {
	Result  Tunnel `json:"result"`
	Success bool   `json:"success"`
	Errors  []struct {
		Message string `json:"message"`
	} `json:"errors"`
}

func (c *Client) doRequest(method, path string, body interface{}) (*http.Response, error) {
	var reqBody io.Reader
	if body != nil {
		b, err := json.Marshal(body)
		if err != nil {
			return nil, err
		}
		reqBody = bytes.NewBuffer(b)
	}

	// path should be something like "/accounts/{id}/tunnels"
	// We append it to BaseURL
	url := c.BaseURL + path

	req, err := http.NewRequest(method, url, reqBody)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+c.ApiToken)
	req.Header.Set("Content-Type", "application/json")

	return c.HTTPClient.Do(req)
}

func NewClient(apiToken, accountId string) *Client {
	if apiToken == "" {
		return nil
	}
	if accountId == "" {
		return nil
	}
	return &Client{
		HTTPClient: &http.Client{Timeout: 10 * time.Second},
		ApiToken:   apiToken,
		AccountId:  accountId,
		BaseURL:    "https://api.cloudflare.com/client/v4",
	}
}

func (c *Client) CreateTunnel(name, secret string) (*Tunnel, error) {
	reqData := Tunnel{
		Name:   name,
		Secret: secret,
	}

	path := fmt.Sprintf("/accounts/%s/cfd_tunnel", c.AccountId)
	resp, err := c.doRequest("POST", path, reqData)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("failed to create tunnel, status: %d, body: %s", resp.StatusCode, string(bodyBytes))
	}

	var tunnelResp TunnelResponse
	if err := json.NewDecoder(resp.Body).Decode(&tunnelResp); err != nil {
		return nil, err
	}

	if !tunnelResp.Success {
		// Handle Cloudflare API errors
		return nil, fmt.Errorf("cloudflare api error: %v", tunnelResp.Errors)
	}

	return &tunnelResp.Result, nil
}

func (c *Client) DeleteTunnel(id string) (*Tunnel, error) {
	path := fmt.Sprintf("/accounts/%s/cfd_tunnel/%s", c.AccountId, id)
	resp, err := c.doRequest("DELETE", path, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("failed to delete tunnel, status: %d, body: %s", resp.StatusCode, string(bodyBytes))
	}

	var tunnelResp TunnelResponse
	if err := json.NewDecoder(resp.Body).Decode(&tunnelResp); err != nil {
		return nil, err
	}

	if !tunnelResp.Success {
		// Handle Cloudflare API errors
		return nil, fmt.Errorf("cloudflare api error: %v", tunnelResp.Errors)
	}

	return &tunnelResp.Result, nil
}

func (c *Client) GetTunnel(id string) (*Tunnel, error) {
	path := fmt.Sprintf("/accounts/%s/cfd_tunnel/%s", c.AccountId, id)
	resp, err := c.doRequest("GET", path, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("failed to get tunnel, status: %d, body: %s", resp.StatusCode, string(bodyBytes))
	}

	var tunnelResp TunnelResponse
	if err := json.NewDecoder(resp.Body).Decode(&tunnelResp); err != nil {
		return nil, err
	}

	if !tunnelResp.Success {
		// Handle Cloudflare API errors
		return nil, fmt.Errorf("cloudflare api error: %v", tunnelResp.Errors)
	}

	return &tunnelResp.Result, nil
}

func (c *Client) UpdateTunnel(id string, name, secret string) (*Tunnel, error) {
	reqData := Tunnel{
		ID:     id,
		Name:   name,
		Secret: secret,
	}

	path := fmt.Sprintf("/accounts/%s/cfd_tunnel/%s", c.AccountId, id)
	resp, err := c.doRequest("PATCH", path, reqData)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("failed to update tunnel, status: %d, body: %s, reqData: %v", resp.StatusCode, string(bodyBytes), reqData)
	}

	var tunnelResp TunnelResponse
	if err := json.NewDecoder(resp.Body).Decode(&tunnelResp); err != nil {
		return nil, err
	}

	if !tunnelResp.Success {
		// Handle Cloudflare API errors
		return nil, fmt.Errorf("cloudflare api error: %v", tunnelResp.Errors)
	}

	return &tunnelResp.Result, nil
}
