package orchestration

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type RunnerClient struct {
	base string
	http *http.Client
}

func NewRunnerClient(base string) *RunnerClient {
	return &RunnerClient{
		base: base,
		http: &http.Client{Timeout: 120 * time.Second},
	}
}

type RunRequest struct {
	AgentName string `json:"agentName"`
	ModelTier string `json:"modelTier"` // tier1|tier2|tier3 (sidecar decides actual model)
	TaskID    string `json:"taskId"`
	Title     string `json:"title"`
	Prompt    string `json:"prompt"`
}

type RunResponse struct {
	Output string `json:"output"`
}

func (r *RunnerClient) Run(req RunRequest) (*RunResponse, error) {
	b, _ := json.Marshal(req)
	httpReq, _ := http.NewRequest("POST", r.base+"/run", bytes.NewReader(b))
	httpReq.Header.Set("Content-Type", "application/json")
	resp, err := r.http.Do(httpReq)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	raw, _ := io.ReadAll(resp.Body)
	if resp.StatusCode >= 300 {
		return nil, fmt.Errorf("runner %d: %s", resp.StatusCode, string(raw))
	}
	var out RunResponse
	if err := json.Unmarshal(raw, &out); err != nil {
		return nil, err
	}
	return &out, nil
}
