package httpserver

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/flaboy/painter/internal/api"
	"github.com/flaboy/painter/internal/app"
)

type envelope struct {
	Data  any `json:"data"`
	Error any `json:"error"`
}

type imagesService interface {
	Generate(ctx context.Context, req api.GenerateImageRequest) (app.Result, *app.ServiceError)
	Edit(ctx context.Context, req api.EditImageRequest) (app.Result, *app.ServiceError)
	Convert(ctx context.Context, req api.ConvertImageRequest) (app.Result, *app.ServiceError)
}

func NewHandler(services ...imagesService) http.Handler {
	var imageSvc imagesService
	if len(services) > 0 {
		imageSvc = services[0]
	}
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.Method == http.MethodGet && r.URL.Path == "/healthz":
			writeOK(w, http.StatusOK, map[string]any{"ok": true})
		case imageSvc != nil && r.Method == http.MethodPost && r.URL.Path == "/v1/images/generate":
			handle(w, r, imageSvc.Generate)
		case imageSvc != nil && r.Method == http.MethodPost && r.URL.Path == "/v1/images/edit":
			handle(w, r, imageSvc.Edit)
		case imageSvc != nil && r.Method == http.MethodPost && r.URL.Path == "/v1/images/convert":
			handle(w, r, imageSvc.Convert)
		default:
			writeErr(w, http.StatusNotFound, "NOT_FOUND", "route not found")
		}
	})
}

type validator interface {
	Validate() error
}

func handle[T validator, R any](w http.ResponseWriter, r *http.Request, fn func(context.Context, T) (R, *app.ServiceError)) {
	var req T
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeErr(w, http.StatusBadRequest, "INVALID_REQUEST", "invalid json body")
		return
	}
	if err := req.Validate(); err != nil {
		if vErr, ok := err.(*api.ValidationError); ok {
			writeErr(w, http.StatusBadRequest, vErr.Code, vErr.Message)
			return
		}
		writeErr(w, http.StatusBadRequest, "INVALID_REQUEST", err.Error())
		return
	}
	out, svcErr := fn(r.Context(), req)
	if svcErr != nil {
		writeErr(w, statusCodeFor(svcErr.Code), svcErr.Code, svcErr.Message)
		return
	}
	writeOK(w, http.StatusOK, out)
}

func writeOK(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(envelope{Data: data, Error: nil})
}

func writeErr(w http.ResponseWriter, status int, code, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(envelope{
		Data: nil,
		Error: map[string]any{
			"code":    code,
			"message": message,
		},
	})
}

func statusCodeFor(code string) int {
	switch code {
	case "INVALID_REQUEST", "IMAGE_FETCH_FAILED", "IMAGE_DECODE_FAILED", "IMAGE_ENCODE_FAILED", "UNSUPPORTED_FORMAT", "UNSUPPORTED_MODE":
		return http.StatusBadRequest
	default:
		return http.StatusInternalServerError
	}
}
