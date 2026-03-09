package httpserver

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestHealthz(t *testing.T) {
	handler := NewHandler()

	req := httptest.NewRequest(http.MethodGet, "/healthz", nil)
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", rr.Code, http.StatusOK)
	}

	var body struct {
		Data struct {
			OK bool `json:"ok"`
		} `json:"data"`
		Error any `json:"error"`
	}
	if err := json.Unmarshal(rr.Body.Bytes(), &body); err != nil {
		t.Fatalf("unmarshal response: %v", err)
	}
	if !body.Data.OK {
		t.Fatalf("data.ok = false, want true")
	}
	if body.Error != nil {
		t.Fatalf("error = %+v, want nil", body.Error)
	}
}

func TestImageRouteRejectsMissingInternalToken(t *testing.T) {
	handler := NewHandlerWithConfig(Config{InternalToken: "secret"}, fakeImagesService{
		generateFn: noopGenerate,
		editFn:     noopEdit,
		convertFn:  noopConvert,
	})

	req := httptest.NewRequest(http.MethodPost, "/v1/images/generate", strings.NewReader(`{"prompt":"poster","size":{"width":1024,"height":1024}}`))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusUnauthorized {
		t.Fatalf("status = %d, want %d", rr.Code, http.StatusUnauthorized)
	}
}
