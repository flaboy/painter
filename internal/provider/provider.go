package provider

import (
	"context"

	"github.com/flaboy/painter/internal/api"
)

type ImageProvider interface {
	Generate(ctx context.Context, req api.GenerateImageRequest) (api.ImageResult, string, string, error)
	Edit(ctx context.Context, req api.EditImageRequest) (api.ImageResult, string, string, error)
}
