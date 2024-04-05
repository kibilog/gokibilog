package gokibilog

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type client struct {
	baseUrl   string
	authToken string
}

func (c *client) SetToken(token string) {
	c.authToken = token
}

func (c *client) Send(logPool *LogPool) error {
	httpClient := &http.Client{
		Transport: &http.Transport{
			IdleConnTimeout: 3 * time.Second,
		},
	}

	body, err := json.Marshal(logPool.messages)
	if err != nil {
		return err
	}

	req, err := http.NewRequest(
		http.MethodPost,
		fmt.Sprintf("%s/%s", c.baseUrl, logPool.getLogId()),
		bytes.NewReader(body),
	)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("apiToken", c.authToken)

	resp, err := httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err = io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if resp.StatusCode != 200 {
		return fmt.Errorf("Kibilog.com returned the %v status code. Response: %s", resp.StatusCode, string(body))
	}

	return nil
}

var clientInstance *client

func getClientInstance() *client {
	if clientInstance == nil {
		clientInstance = &client{
			baseUrl: "https://kibilog.com/api/v1/log/monolog",
		}
	}
	return clientInstance
}
