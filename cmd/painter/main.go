package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/flaboy/painter/internal/api"
	"github.com/flaboy/painter/internal/app"
	"github.com/flaboy/painter/internal/httpserver"
	"github.com/flaboy/painter/internal/imageops"
	"github.com/flaboy/painter/internal/provider"
	"github.com/flaboy/painter/internal/usage"
)

func main() {
	addr := ":" + portFromEnv()
	imageSvc := app.NewService(
		provider.NewFakeProvider(),
		imageConverter{},
		usage.NewWebhookReporter(
			os.Getenv("PAINTER_USAGE_WEBHOOK_URL"),
			os.Getenv("PAINTER_USAGE_WEBHOOK_TOKEN"),
		),
	)
	handler := httpserver.NewHandlerWithConfig(httpserver.Config{
		InternalToken: os.Getenv("PAINTER_INTERNAL_TOKEN"),
	}, imageSvc)
	if err := http.ListenAndServe(addr, handler); err != nil {
		log.Fatal(err)
	}
}

func portFromEnv() string {
	if v := os.Getenv("PAINTER_PORT"); v != "" {
		return v
	}
	return "7013"
}

type imageConverter struct{}

func (imageConverter) Convert(ctx context.Context, req app.ConvertRequest) (api.ImageResult, error) {
	return imageops.Convert(ctx, imageops.ConvertRequest(req))
}
