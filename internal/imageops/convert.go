package imageops

import (
	"bytes"
	"context"
	"encoding/base64"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/jpeg"
	"image/png"
	"strings"

	"github.com/chai2010/webp"
	"github.com/disintegration/imaging"
	"github.com/flaboy/painter/internal/api"
)

type ConvertRequest struct {
	SourceURL  string
	Format     string
	Resize     api.Resize
	Quality    int
	Background string
}

func Convert(ctx context.Context, req ConvertRequest) (api.ImageResult, error) {
	body, _, err := FetchSource(ctx, req.SourceURL)
	if err != nil {
		return api.ImageResult{}, err
	}

	src, _, err := image.Decode(bytes.NewReader(body))
	if err != nil {
		return api.ImageResult{}, fmt.Errorf("IMAGE_DECODE_FAILED")
	}

	img := imaging.Clone(src)
	if req.Resize.Width > 0 || req.Resize.Height > 0 {
		width := req.Resize.Width
		height := req.Resize.Height
		if width <= 0 {
			width = 0
		}
		if height <= 0 {
			height = 0
		}
		img = imaging.Fit(img, width, height, imaging.Lanczos)
	}

	format := normalizeFormat(req.Format)
	var encoded []byte
	switch format {
	case "png":
		encoded, err = encodePNG(img)
	case "jpeg":
		if strings.EqualFold(req.Background, "white") {
			img = flattenOnWhite(img)
		}
		encoded, err = encodeJPEG(img, req.Quality)
	case "webp":
		encoded, err = encodeWEBP(img, req.Quality)
	default:
		return api.ImageResult{}, fmt.Errorf("UNSUPPORTED_FORMAT")
	}
	if err != nil {
		return api.ImageResult{}, err
	}

	return api.ImageResult{
		MimeType:    mimeTypeFor(format),
		Format:      format,
		Width:       img.Bounds().Dx(),
		Height:      img.Bounds().Dy(),
		BytesBase64: base64.StdEncoding.EncodeToString(encoded),
	}, nil
}

func normalizeFormat(format string) string {
	switch strings.ToLower(strings.TrimSpace(format)) {
	case "", "png":
		return "png"
	case "jpg", "jpeg":
		return "jpeg"
	case "webp":
		return "webp"
	default:
		return strings.ToLower(strings.TrimSpace(format))
	}
}

func mimeTypeFor(format string) string {
	switch format {
	case "jpeg":
		return "image/jpeg"
	case "webp":
		return "image/webp"
	default:
		return "image/png"
	}
}

func encodePNG(img image.Image) ([]byte, error) {
	var buf bytes.Buffer
	if err := png.Encode(&buf, img); err != nil {
		return nil, fmt.Errorf("IMAGE_ENCODE_FAILED")
	}
	return buf.Bytes(), nil
}

func encodeJPEG(img image.Image, quality int) ([]byte, error) {
	if quality <= 0 {
		quality = 85
	}
	var buf bytes.Buffer
	if err := jpeg.Encode(&buf, img, &jpeg.Options{Quality: quality}); err != nil {
		return nil, fmt.Errorf("IMAGE_ENCODE_FAILED")
	}
	return buf.Bytes(), nil
}

func encodeWEBP(img image.Image, quality int) ([]byte, error) {
	if quality <= 0 {
		quality = 85
	}
	var buf bytes.Buffer
	if err := webp.Encode(&buf, img, &webp.Options{Quality: float32(quality)}); err != nil {
		return nil, fmt.Errorf("IMAGE_ENCODE_FAILED")
	}
	return buf.Bytes(), nil
}

func flattenOnWhite(src image.Image) *image.NRGBA {
	bounds := src.Bounds()
	dst := image.NewNRGBA(bounds)
	draw.Draw(dst, bounds, &image.Uniform{C: color.White}, image.Point{}, draw.Src)
	draw.Draw(dst, bounds, src, bounds.Min, draw.Over)
	return dst
}
