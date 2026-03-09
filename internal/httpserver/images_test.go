package httpserver

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/flaboy/painter/internal/api"
	"github.com/flaboy/painter/internal/app"
)

type fakeImagesService struct {
	generateFn func(context.Context, api.GenerateImageRequest) (app.Result, *app.ServiceError)
	editFn     func(context.Context, api.EditImageRequest) (app.Result, *app.ServiceError)
	convertFn  func(context.Context, api.ConvertImageRequest) (app.Result, *app.ServiceError)
}

func (f fakeImagesService) Generate(ctx context.Context, req api.GenerateImageRequest) (app.Result, *app.ServiceError) {
	return f.generateFn(ctx, req)
}

func (f fakeImagesService) Edit(ctx context.Context, req api.EditImageRequest) (app.Result, *app.ServiceError) {
	return f.editFn(ctx, req)
}

func (f fakeImagesService) Convert(ctx context.Context, req api.ConvertImageRequest) (app.Result, *app.ServiceError) {
	return f.convertFn(ctx, req)
}

func TestGenerateRouteReturns200(t *testing.T) {
	handler := NewHandler(fakeImagesService{
		generateFn: func(_ context.Context, _ api.GenerateImageRequest) (app.Result, *app.ServiceError) {
			return app.Result{
				Image:    api.ImageResult{Format: "png", BytesBase64: "abc"},
				Provider: "fake",
				Model:    "fake-image-v1",
			}, nil
		},
		editFn:    noopEdit,
		convertFn: noopConvert,
	})

	req := httptest.NewRequest(http.MethodPost, "/v1/images/generate", bytes.NewBufferString(`{"prompt":"poster","size":{"width":1024,"height":1024},"format":"png"}`))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200", rr.Code)
	}
}

func TestEditRouteReturns400ForBadRequest(t *testing.T) {
	handler := NewHandler(fakeImagesService{
		generateFn: noopGenerate,
		editFn:     noopEdit,
		convertFn:  noopConvert,
	})

	req := httptest.NewRequest(http.MethodPost, "/v1/images/edit", bytes.NewBufferString(`{"mode":"variation"}`))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want 400", rr.Code)
	}
}

func TestConvertRouteReturnsResponseEnvelope(t *testing.T) {
	handler := NewHandler(fakeImagesService{
		generateFn: noopGenerate,
		editFn:     noopEdit,
		convertFn: func(_ context.Context, _ api.ConvertImageRequest) (app.Result, *app.ServiceError) {
			return app.Result{Image: api.ImageResult{Format: "png", BytesBase64: "abc"}}, nil
		},
	})

	req := httptest.NewRequest(http.MethodPost, "/v1/images/convert", bytes.NewBufferString(`{"sourceUrl":"https://example.com/a.png","format":"png"}`))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200", rr.Code)
	}
	var out map[string]any
	if err := json.Unmarshal(rr.Body.Bytes(), &out); err != nil {
		t.Fatalf("unmarshal response: %v", err)
	}
	if out["error"] != nil {
		t.Fatalf("error = %+v, want nil", out["error"])
	}
}

func noopGenerate(_ context.Context, _ api.GenerateImageRequest) (app.Result, *app.ServiceError) {
	return app.Result{}, nil
}

func noopEdit(_ context.Context, _ api.EditImageRequest) (app.Result, *app.ServiceError) {
	return app.Result{}, &app.ServiceError{Code: "INVALID_REQUEST", Message: "invalid request"}
}

func noopConvert(_ context.Context, _ api.ConvertImageRequest) (app.Result, *app.ServiceError) {
	return app.Result{}, nil
}
