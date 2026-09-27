package wan

import (
	"context"
	"fmt"
	"time"

	"dramastudio/internal/platform/ai/contracts"
	"dramastudio/internal/platform/ai/providers/httpjson"
)

// Adapter is an async video-generation adapter in the shape providers
// hosting Wan-class models expose: submit returns a provider job id, status
// is retrieved via poll or verified webhook (Backend.md §61).
type Adapter struct {
	name   string
	client *httpjson.Client
}

func New(providerName, baseURL, apiKey string) *Adapter {
	name := providerName
	if name == "" {
		name = "wan"
	}
	return &Adapter{
		name:   name,
		client: httpjson.New(name, baseURL, apiKey, httpjson.AuthBearer, 60*time.Second),
	}
}

func (a *Adapter) ProviderName() string { return a.name }

type submitRequest struct {
	Model       string                 `json:"model"`
	Input       contracts.GenerationSpec `json:"input"`
	CallbackURL string                 `json:"callback_url,omitempty"`
}

type submitResponse struct {
	JobID  string `json:"job_id"`
	Status string `json:"status"`
}

func (a *Adapter) Submit(ctx context.Context, model string, spec contracts.GenerationSpec, callbackURL string) (*contracts.AsyncJob, error) {
	var out submitResponse
	err := a.client.Do(ctx, "POST", "/generations", submitRequest{
		Model:       model,
		Input:       spec,
		CallbackURL: callbackURL,
	}, &out)
	if err != nil {
		return nil, err
	}
	if out.JobID == "" {
		return nil, &contracts.ProviderError{Kind: contracts.ErrUnknown, Provider: a.name, Message: "provider returned no job id"}
	}
	return &contracts.AsyncJob{ProviderJobID: out.JobID, Status: normalizeStatus(out.Status)}, nil
}

type statusResponse struct {
	JobID     string  `json:"job_id"`
	Status    string  `json:"status"`
	OutputURL string  `json:"output_url"`
	Cost      float64 `json:"cost"`
	Error     string  `json:"error"`
}

func (a *Adapter) Poll(ctx context.Context, providerJobID string) (*contracts.AsyncJob, error) {
	var out statusResponse
	if err := a.client.Do(ctx, "GET", fmt.Sprintf("/generations/%s", providerJobID), nil, &out); err != nil {
		return nil, err
	}
	job := &contracts.AsyncJob{ProviderJobID: providerJobID, Status: normalizeStatus(out.Status)}
	if job.Status == "succeeded" {
		job.Result = &contracts.Result{
			Provider:      a.name,
			OutputURI:     out.OutputURL,
			Cost:          out.Cost,
			ProviderJobID: providerJobID,
		}
	}
	if job.Status == "failed" && out.Error != "" {
		return job, &contracts.ProviderError{Kind: contracts.ErrRejected, Provider: a.name, Message: out.Error}
	}
	return job, nil
}

func (a *Adapter) Cancel(ctx context.Context, providerJobID string) error {
	return a.client.Do(ctx, "DELETE", fmt.Sprintf("/generations/%s", providerJobID), nil, nil)
}

func normalizeStatus(s string) string {
	switch s {
	case "succeeded", "completed", "SUCCESS":
		return "succeeded"
	case "failed", "error", "FAILED":
		return "failed"
	case "running", "processing", "RUNNING":
		return "running"
	default:
		return "pending"
	}
}
