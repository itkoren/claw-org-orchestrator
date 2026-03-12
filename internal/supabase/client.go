package supabase

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/itkoren/claw-org-orchestrator/internal/config"
)

type Client struct {
	base string
	key  string
	http *http.Client
}

func MustNew(cfg *config.Config) *Client {
	return &Client{
		base: cfg.SupabaseURL,
		key:  cfg.SupabaseServiceRoleKey,
		http: &http.Client{},
	}
}

func (c *Client) restURL(path string, q url.Values) string {
	u := c.base + "/rest/v1" + path
	if q != nil {
		u += "?" + q.Encode()
	}
	return u
}

func (c *Client) do(req *http.Request) ([]byte, error) {
	req.Header.Set("apikey", c.key)
	req.Header.Set("Authorization", "Bearer "+c.key)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Prefer", "return=representation")

	resp, err := c.http.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	b, _ := io.ReadAll(resp.Body)
	if resp.StatusCode >= 300 {
		return nil, fmt.Errorf("supabase %d: %s", resp.StatusCode, string(b))
	}
	return b, nil
}

func (c *Client) ListQueuedTasks(limit int) ([]Task, error) {
	q := url.Values{}
	q.Set("status", "eq.queued")
	q.Set("order", "created_at.asc")
	q.Set("limit", fmt.Sprintf("%d", limit))

	req, _ := http.NewRequest("GET", c.restURL("/tasks", q), nil)
	b, err := c.do(req)
	if err != nil {
		return nil, err
	}
	var out []Task
	return out, json.Unmarshal(b, &out)
}

func (c *Client) GetTask(id string) (*Task, error) {
	q := url.Values{}
	q.Set("id", "eq."+id)
	q.Set("limit", "1")

	req, _ := http.NewRequest("GET", c.restURL("/tasks", q), nil)
	b, err := c.do(req)
	if err != nil {
		return nil, err
	}
	var arr []Task
	if err := json.Unmarshal(b, &arr); err != nil {
		return nil, err
	}
	if len(arr) == 0 {
		return nil, fmt.Errorf("task not found: %s", id)
	}
	return &arr[0], nil
}

func (c *Client) PatchTask(id string, patch map[string]any) error {
	q := url.Values{}
	q.Set("id", "eq."+id)

	body, _ := json.Marshal(patch)
	req, _ := http.NewRequest("PATCH", c.restURL("/tasks", q), bytes.NewReader(body))
	_, err := c.do(req)
	return err
}

func (c *Client) MarkTaskRunning(id string) error {
	return c.PatchTask(id, map[string]any{"status": "running"})
}

func (c *Client) MarkTaskDone(id string) error {
	return c.PatchTask(id, map[string]any{"status": "done"})
}

func (c *Client) MarkTaskFailed(id, reason string) error {
	return c.PatchTask(id, map[string]any{"status": "failed", "description": reason})
}

func (c *Client) SetSlackThread(id, channel, threadTS string) error {
	return c.PatchTask(id, map[string]any{
		"slack_channel_id": channel,
		"slack_thread_ts":  threadTS,
	})
}
