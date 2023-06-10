package helpers_test

import (
	"net/url"
	"testing"
	"web-crawler/modules/helpers"
)

func TestTransformLink(t *testing.T) {
	mockURL, parseErr := url.Parse("http://example.com/foo/bar?test=123&bar=456")
	if parseErr != nil {
		t.Fatal(parseErr)
	}

	expected := "http://example.com/foo/bar"

	if actual := helpers.TransformLink(mockURL); actual.String() != expected {
		t.Fatalf("expected %s, got %s", expected, actual)
	}
}
