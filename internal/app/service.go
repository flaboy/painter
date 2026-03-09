package app

import (
	"context"

	"github.com/flaboy/painter/internal/api"
	"github.com/flaboy/painter/internal/imageops"
	"github.com/flaboy/painter/internal/provider"
)

type Result struct {
	Image    api.ImageResult `json:"image"`
	Provider string          `json:"provider,omitempty"`
	Model    string          `json:"model,omitempty"`
}

type ServiceError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

type Converter interface {
	Convert(ctx context.Context, req ConvertRequest) (api.ImageResult, error)
}

type ConvertRequest = imageops.ConvertRequest

type Service struct {
	provider  provider.ImageProvider
	converter Converter
}

func NewService(provider provider.ImageProvider, converter Converter) *Service {
	return &Service{provider: provider, converter: converter}
}

func (s *Service) Generate(ctx context.Context, req api.GenerateImageRequest) (Result, *ServiceError) {
	image, providerName, modelName, err := s.provider.Generate(ctx, req)
	if err != nil {
		return Result{}, mapError("IMAGE_GENERATE_FAILED", err)
	}
	return Result{Image: image, Provider: providerName, Model: modelName}, nil
}

func (s *Service) Edit(ctx context.Context, req api.EditImageRequest) (Result, *ServiceError) {
	image, providerName, modelName, err := s.provider.Edit(ctx, req)
	if err != nil {
		return Result{}, mapError("IMAGE_EDIT_FAILED", err)
	}
	return Result{Image: image, Provider: providerName, Model: modelName}, nil
}

func (s *Service) Convert(ctx context.Context, req api.ConvertImageRequest) (Result, *ServiceError) {
	image, err := s.converter.Convert(ctx, ConvertRequest{
		SourceURL:  req.SourceUrl,
		Format:     req.Format,
		Resize:     req.Resize,
		Quality:    req.Quality,
		Background: req.Background,
	})
	if err != nil {
		return Result{}, mapError("IMAGE_CONVERT_FAILED", err)
	}
	return Result{Image: image}, nil
}

func mapError(defaultCode string, err error) *ServiceError {
	if err == nil {
		return nil
	}
	code := err.Error()
	switch code {
	case "INVALID_REQUEST", "UNSUPPORTED_MODE", "IMAGE_FETCH_FAILED", "IMAGE_CONVERT_FAILED", "IMAGE_GENERATE_FAILED", "IMAGE_EDIT_FAILED", "IMAGE_DECODE_FAILED", "IMAGE_ENCODE_FAILED", "UNSUPPORTED_FORMAT":
	default:
		code = defaultCode
	}
	return &ServiceError{Code: code, Message: messageForCode(code)}
}

func messageForCode(code string) string {
	switch code {
	case "UNSUPPORTED_MODE":
		return "unsupported image edit mode"
	case "IMAGE_FETCH_FAILED":
		return "failed to fetch source image"
	case "IMAGE_DECODE_FAILED":
		return "failed to decode image"
	case "IMAGE_ENCODE_FAILED":
		return "failed to encode image"
	case "UNSUPPORTED_FORMAT":
		return "unsupported image format"
	case "IMAGE_EDIT_FAILED":
		return "failed to edit image"
	case "IMAGE_GENERATE_FAILED":
		return "failed to generate image"
	case "IMAGE_CONVERT_FAILED":
		return "failed to convert image"
	default:
		return "request failed"
	}
}
