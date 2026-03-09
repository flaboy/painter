package provider

import (
	"context"
	"errors"

	"github.com/flaboy/painter/internal/api"
)

const fakePNGBase64 = "iVBORw0KGgoAAAANSUhEUgAAAAEAAAABCAQAAAC1HAwCAAAAC0lEQVR42mP8/x8AAusB9WnR6j0AAAAASUVORK5CYII="

type FakeProvider struct{}

func NewFakeProvider() ImageProvider {
	return FakeProvider{}
}

func (FakeProvider) Generate(_ context.Context, req api.GenerateImageRequest) (api.ImageResult, string, string, error) {
	return api.ImageResult{
		MimeType:    "image/png",
		Format:      normalizeFormat(req.Format),
		Width:       normalizeDim(req.Size.Width),
		Height:      normalizeDim(req.Size.Height),
		BytesBase64: fakePNGBase64,
	}, "fake", "fake-image-v1", nil
}

func (FakeProvider) Edit(_ context.Context, req api.EditImageRequest) (api.ImageResult, string, string, error) {
	switch req.Mode {
	case "variation", "expand", "mask_edit":
	default:
		return api.ImageResult{}, "", "", errors.New("UNSUPPORTED_MODE")
	}
	return api.ImageResult{
		MimeType:    "image/png",
		Format:      normalizeFormat(req.Format),
		Width:       normalizeDim(req.Size.Width),
		Height:      normalizeDim(req.Size.Height),
		BytesBase64: fakePNGBase64,
	}, "fake", "fake-image-v1", nil
}

func normalizeFormat(format string) string {
	if format == "" {
		return "png"
	}
	return format
}

func normalizeDim(v int) int {
	if v <= 0 {
		return 1
	}
	return v
}
