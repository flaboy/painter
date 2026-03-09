package provider

import (
	"context"
	"testing"

	"github.com/flaboy/painter/internal/api"
)

func TestFakeProviderGenerate(t *testing.T) {
	p := NewFakeProvider()

	image, providerName, modelName, err := p.Generate(context.Background(), api.GenerateImageRequest{
		Prompt: "a red square icon",
		Size:   api.ImageSize{Width: 1, Height: 1},
		Format: "png",
	})
	if err != nil {
		t.Fatalf("Generate: %v", err)
	}
	if providerName != "fake" {
		t.Fatalf("provider = %q, want fake", providerName)
	}
	if modelName == "" {
		t.Fatalf("modelName is empty")
	}
	if image.BytesBase64 == "" {
		t.Fatalf("bytesBase64 is empty")
	}
}

func TestFakeProviderEdit(t *testing.T) {
	p := NewFakeProvider()

	image, providerName, modelName, err := p.Edit(context.Background(), api.EditImageRequest{
		Mode:      "variation",
		SourceUrl: "https://example.com/source.png",
		Format:    "png",
		Size:      api.ImageSize{Width: 1, Height: 1},
	})
	if err != nil {
		t.Fatalf("Edit: %v", err)
	}
	if providerName != "fake" {
		t.Fatalf("provider = %q, want fake", providerName)
	}
	if modelName == "" {
		t.Fatalf("modelName is empty")
	}
	if image.BytesBase64 == "" {
		t.Fatalf("bytesBase64 is empty")
	}
}

func TestFakeProviderRejectsUnsupportedEditMode(t *testing.T) {
	p := NewFakeProvider()

	_, _, _, err := p.Edit(context.Background(), api.EditImageRequest{
		Mode:      "replace_sky",
		SourceUrl: "https://example.com/source.png",
	})
	if err == nil {
		t.Fatal("expected error")
	}
	if err.Error() != "UNSUPPORTED_MODE" {
		t.Fatalf("error = %q, want UNSUPPORTED_MODE", err.Error())
	}
}
