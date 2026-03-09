package app

import (
	"context"
	"errors"
	"testing"

	"github.com/flaboy/painter/internal/api"
)

type fakeImageProvider struct {
	generateFn func(context.Context, api.GenerateImageRequest) (api.ImageResult, string, string, error)
	editFn     func(context.Context, api.EditImageRequest) (api.ImageResult, string, string, error)
}

func (f fakeImageProvider) Generate(ctx context.Context, req api.GenerateImageRequest) (api.ImageResult, string, string, error) {
	return f.generateFn(ctx, req)
}

func (f fakeImageProvider) Edit(ctx context.Context, req api.EditImageRequest) (api.ImageResult, string, string, error) {
	return f.editFn(ctx, req)
}

type fakeConverter struct {
	convertFn func(context.Context, ConvertRequest) (api.ImageResult, error)
}

func (f fakeConverter) Convert(ctx context.Context, req ConvertRequest) (api.ImageResult, error) {
	return f.convertFn(ctx, req)
}

func TestPainterServiceGenerate(t *testing.T) {
	svc := NewService(fakeImageProvider{
		generateFn: func(_ context.Context, _ api.GenerateImageRequest) (api.ImageResult, string, string, error) {
			return api.ImageResult{Format: "png", BytesBase64: "abc"}, "fake", "fake-image-v1", nil
		},
		editFn: func(_ context.Context, _ api.EditImageRequest) (api.ImageResult, string, string, error) {
			return api.ImageResult{}, "", "", nil
		},
	}, fakeConverter{})

	out, svcErr := svc.Generate(context.Background(), api.GenerateImageRequest{Prompt: "poster"})
	if svcErr != nil {
		t.Fatalf("Generate returned error: %+v", svcErr)
	}
	if out.Provider != "fake" {
		t.Fatalf("provider = %q, want fake", out.Provider)
	}
	if out.Image.BytesBase64 == "" {
		t.Fatal("bytesBase64 is empty")
	}
}

func TestPainterServiceEdit(t *testing.T) {
	svc := NewService(fakeImageProvider{
		generateFn: func(_ context.Context, _ api.GenerateImageRequest) (api.ImageResult, string, string, error) {
			return api.ImageResult{}, "", "", nil
		},
		editFn: func(_ context.Context, _ api.EditImageRequest) (api.ImageResult, string, string, error) {
			return api.ImageResult{Format: "png", BytesBase64: "abc"}, "fake", "fake-image-v1", nil
		},
	}, fakeConverter{})

	out, svcErr := svc.Edit(context.Background(), api.EditImageRequest{
		Mode:      "variation",
		SourceUrl: "https://example.com/source.png",
	})
	if svcErr != nil {
		t.Fatalf("Edit returned error: %+v", svcErr)
	}
	if out.Provider != "fake" {
		t.Fatalf("provider = %q, want fake", out.Provider)
	}
}

func TestPainterServiceConvert(t *testing.T) {
	svc := NewService(fakeImageProvider{
		generateFn: func(_ context.Context, _ api.GenerateImageRequest) (api.ImageResult, string, string, error) {
			return api.ImageResult{}, "", "", nil
		},
		editFn: func(_ context.Context, _ api.EditImageRequest) (api.ImageResult, string, string, error) {
			return api.ImageResult{}, "", "", nil
		},
	}, fakeConverter{
		convertFn: func(_ context.Context, _ ConvertRequest) (api.ImageResult, error) {
			return api.ImageResult{Format: "webp", BytesBase64: "abc"}, nil
		},
	})

	out, svcErr := svc.Convert(context.Background(), api.ConvertImageRequest{
		SourceUrl: "https://example.com/source.png",
		Format:    "webp",
	})
	if svcErr != nil {
		t.Fatalf("Convert returned error: %+v", svcErr)
	}
	if out.Image.Format != "webp" {
		t.Fatalf("format = %q, want webp", out.Image.Format)
	}
}

func TestPainterServiceConvertMapsFetchFailure(t *testing.T) {
	svc := NewService(fakeImageProvider{
		generateFn: func(_ context.Context, _ api.GenerateImageRequest) (api.ImageResult, string, string, error) {
			return api.ImageResult{}, "", "", nil
		},
		editFn: func(_ context.Context, _ api.EditImageRequest) (api.ImageResult, string, string, error) {
			return api.ImageResult{}, "", "", nil
		},
	}, fakeConverter{
		convertFn: func(_ context.Context, _ ConvertRequest) (api.ImageResult, error) {
			return api.ImageResult{}, errors.New("IMAGE_FETCH_FAILED")
		},
	})

	_, svcErr := svc.Convert(context.Background(), api.ConvertImageRequest{
		SourceUrl: "https://example.com/source.png",
		Format:    "webp",
	})
	if svcErr == nil {
		t.Fatal("expected service error")
	}
	if svcErr.Code != "IMAGE_FETCH_FAILED" {
		t.Fatalf("code = %q, want IMAGE_FETCH_FAILED", svcErr.Code)
	}
}
