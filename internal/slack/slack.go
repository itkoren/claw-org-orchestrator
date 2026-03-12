package slack

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/itkoren/claw-org-orchestrator/internal/config"
)

type Client struct {
	token string
	http  *http.Client
	orgCh string
}

func MustNew(cfg *config.Config) *Client {
	return &Client{token: cfg.SlackBotToken, http: &http.Client{}, orgCh: cfg.SlackOrgChannelID}
}

type postMessageResp struct {
	OK       bool   `json:"ok"`
	Error    string `json:"error"`
	Channel  string `json:"channel"`
	TS       string `json:"ts"`
	ThreadTS string `json:"thread_ts"`
}

func (c *Client) post(endpoint string, payload any) (*postMessageResp, error) {
	b, _ := json.Marshal(payload)
	req, _ := http.NewRequest("POST", "https://slack.com/api/"+endpoint, bytes.NewReader(b))
	req.Header.Set("Authorization", "Bearer "+c.token)
	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	resp, err := c.http.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	raw, _ := io.ReadAll(resp.Body)

	var out postMessageResp
	_ = json.Unmarshal(raw, &out)
	if resp.StatusCode >= 300 || !out.OK {
		return nil, fmt.Errorf("slack error (%s): %s", endpoint, string(raw))
	}
	return &out, nil
}

func (c *Client) PostOrgAnnouncement(text string) error {
	_, err := c.post("chat.postMessage", map[string]any{
		"channel": c.orgCh,
		"text":    text,
	})
	return err
}

func (c *Client) PostRoot(channel, text string) (ts string, err error) {
	resp, err := c.post("chat.postMessage", map[string]any{"channel": channel, "text": text})
	if err != nil {
		return "", err
	}
	return resp.TS, nil
}

func (c *Client) Reply(channel, threadTS, text string) error {
	_, err := c.post("chat.postMessage", map[string]any{
		"channel":   channel,
		"text":      text,
		"thread_ts": threadTS,
	})
	return err
}
