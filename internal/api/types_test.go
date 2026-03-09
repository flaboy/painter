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
