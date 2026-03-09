package api

import (
	"reflect"
	"testing"
)

func TestImageAPIContracts(t *testing.T) {
	assertHasJSONTag(t, reflect.TypeOf(GenerateImageRequest{}), "prompt")
	assertHasJSONTag(t, reflect.TypeOf(EditImageRequest{}), "sourceUrl")
	assertHasJSONTag(t, reflect.TypeOf(EditImageRequest{}), "maskUrl")
	assertHasJSONTag(t, reflect.TypeOf(ConvertImageRequest{}), "sourceUrl")
	assertHasJSONTag(t, reflect.TypeOf(ImageResult{}), "bytesBase64")

	assertNoJSONTag(t, reflect.TypeOf(GenerateImageRequest{}), "inputPath")
	assertNoJSONTag(t, reflect.TypeOf(GenerateImageRequest{}), "outputPath")
	assertNoJSONTag(t, reflect.TypeOf(EditImageRequest{}), "inputPath")
	assertNoJSONTag(t, reflect.TypeOf(EditImageRequest{}), "outputPath")
	assertNoJSONTag(t, reflect.TypeOf(ConvertImageRequest{}), "inputPath")
	assertNoJSONTag(t, reflect.TypeOf(ConvertImageRequest{}), "outputPath")
	assertNoJSONTag(t, reflect.TypeOf(ConvertImageRequest{}), "vfs")
	assertNoJSONTag(t, reflect.TypeOf(ConvertImageRequest{}), "bucket")
}

func TestValidateRejectsUnsupportedFormat(t *testing.T) {
	if err := (GenerateImageRequest{
		Prompt: "poster",
		Size:   ImageSize{Width: 1024, Height: 1024},
		Format: "webp",
	}).Validate(); err == nil {
		t.Fatal("expected generate validation error")
	}

	if err := (EditImageRequest{
		Mode:      "variation",
		SourceUrl: "https://example.com/source.png",
		Format:    "webp",
	}).Validate(); err == nil {
		t.Fatal("expected edit validation error")
	}

	if err := (ConvertImageRequest{
		SourceUrl: "https://example.com/source.png",
		Format:    "webp",
	}).Validate(); err == nil {
		t.Fatal("expected convert validation error")
	}
}

func TestUsageContextJSONTags(t *testing.T) {
	assertHasJSONTag(t, reflect.TypeOf(GenerateImageRequest{}), "usageContext")
	assertHasJSONTag(t, reflect.TypeOf(EditImageRequest{}), "usageContext")
	assertHasJSONTag(t, reflect.TypeOf(ConvertImageRequest{}), "usageContext")
	assertHasJSONTag(t, reflect.TypeOf(UsageReportRequest{}), "usageContext")
}

func assertHasJSONTag(t *testing.T, typ reflect.Type, tagValue string) {
	t.Helper()
	for i := 0; i < typ.NumField(); i++ {
		if typ.Field(i).Tag.Get("json") == tagValue || typ.Field(i).Tag.Get("json") == tagValue+",omitempty" {
			return
		}
	}
	t.Fatalf("json tag %q not found on %s", tagValue, typ.Name())
}

func assertNoJSONTag(t *testing.T, typ reflect.Type, forbidden string) {
	t.Helper()
	for i := 0; i < typ.NumField(); i++ {
		if typ.Field(i).Tag.Get("json") == forbidden {
			t.Fatalf("unexpected json tag %q on %s", forbidden, typ.Name())
		}
	}
}
