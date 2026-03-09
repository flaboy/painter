package usage

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/flaboy/painter/internal/api"
)

func TestWebhookReporterPostsUsageEvent(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if got := r.Header.Get("Authorization"); got != "Bearer usage-secret" {
			t.Fatalf("authorization = %q, want Bearer usage-secret", got)
		}
		var body api.UsageReportRequest
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Fatalf("decode body: %v", err)
		}
		if body.Service != "painter" {
			t.Fatalf("service = %q, want painter", body.Service)
		}
		w.WriteHeader(http.StatusNoContent)
	}))
	defer server.Close()

	reporter := NewWebhookReporter(server.URL, "usage-secret")
	if err := reporter.Report(context.Background(), api.UsageReportRequest{Service: "painter"}); err != nil {
		t.Fatalf("report: %v", err)
	}
}
