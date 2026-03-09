package usage

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/flaboy/painter/internal/api"
)

type Reporter interface {
	Report(context.Context, api.UsageReportRequest) error
}

type noopReporter struct{}

func (noopReporter) Report(context.Context, api.UsageReportRequest) error {
	return nil
}

type WebhookReporter struct {
	url    string
	token  string
	client *http.Client
}

func NewWebhookReporter(url, token string) Reporter {
	trimmedURL := strings.TrimSpace(url)
	if trimmedURL == "" {
		return noopReporter{}
	}
	return &WebhookReporter{
		url:   trimmedURL,
		token: strings.TrimSpace(token),
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

func (r *WebhookReporter) Report(ctx context.Context, req api.UsageReportRequest) error {
	body, err := json.Marshal(req)
	if err != nil {
		return err
	}
	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, r.url, bytes.NewReader(body))
	if err != nil {
		return err
	}
	httpReq.Header.Set("Content-Type", "application/json")
	if r.token != "" {
		httpReq.Header.Set("Authorization", "Bearer "+r.token)
	}
	resp, err := r.client.Do(httpReq)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("usage_webhook_failed:%d", resp.StatusCode)
	}
	return nil
}
