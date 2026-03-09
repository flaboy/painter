package imageops

import (
	"bytes"
	"context"
	"image"
	"image/color"
	"image/png"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/flaboy/painter/internal/api"
)

func TestImageOps(t *testing.T) {
	sourcePNG := makeTestPNG(t)
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "image/png")
		_, _ = w.Write(sourcePNG)
	}))
	defer server.Close()

	t.Run("fetch source image", func(t *testing.T) {
		data, contentType, err := FetchSource(context.Background(), server.URL)
		if err != nil {
			t.Fatalf("FetchSource: %v", err)
		}
		if len(data) == 0 {
			t.Fatal("expected non-empty data")
		}
		if contentType != "image/png" {
			t.Fatalf("contentType = %q, want image/png", contentType)
		}
	})

	t.Run("convert png to jpeg with white background", func(t *testing.T) {
		out, err := Convert(context.Background(), ConvertRequest{
			SourceURL:  server.URL,
			Format:     "jpeg",
			Background: "white",
		})
		if err != nil {
			t.Fatalf("Convert: %v", err)
		}
		if out.Format != "jpeg" {
			t.Fatalf("format = %q, want jpeg", out.Format)
		}
		if out.MimeType != "image/jpeg" {
			t.Fatalf("mimeType = %q, want image/jpeg", out.MimeType)
		}
		if out.BytesBase64 == "" {
			t.Fatal("bytesBase64 is empty")
		}
	})

	t.Run("resize with fit inside", func(t *testing.T) {
		out, err := Convert(context.Background(), ConvertRequest{
			SourceURL: server.URL,
			Format:    "png",
			Resize: api.Resize{
				Width:  1,
				Height: 1,
				Fit:    "inside",
			},
		})
		if err != nil {
			t.Fatalf("Convert: %v", err)
		}
		if out.Width != 1 {
			t.Fatalf("width = %d, want 1", out.Width)
		}
		if out.Height != 1 {
			t.Fatalf("height = %d, want 1", out.Height)
		}
	})

	t.Run("convert png to webp", func(t *testing.T) {
		out, err := Convert(context.Background(), ConvertRequest{
			SourceURL: server.URL,
			Format:    "webp",
		})
		if err != nil {
			t.Fatalf("Convert: %v", err)
		}
		if out.Format != "webp" {
			t.Fatalf("format = %q, want webp", out.Format)
		}
		if out.MimeType != "image/webp" {
			t.Fatalf("mimeType = %q, want image/webp", out.MimeType)
		}
		if out.BytesBase64 == "" {
			t.Fatal("bytesBase64 is empty")
		}
	})
}

func makeTestPNG(t *testing.T) []byte {
	t.Helper()
	img := image.NewNRGBA(image.Rect(0, 0, 2, 1))
	img.Set(0, 0, color.NRGBA{R: 255, G: 0, B: 0, A: 255})
	img.Set(1, 0, color.NRGBA{R: 0, G: 0, B: 255, A: 128})

	var buf bytes.Buffer
	if err := png.Encode(&buf, img); err != nil {
		t.Fatalf("png.Encode: %v", err)
	}
	return buf.Bytes()
}
