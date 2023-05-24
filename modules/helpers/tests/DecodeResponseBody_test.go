package helpers_test

import (
	"bytes"
	"io"
	"testing"
	"web-crawler/modules/helpers"
)

type MockResponseData struct {
	Data string `json:"data"`
	Name string `json:"name"`
}

func TestDecodeResponseBody(t *testing.T) {
	responseData := &MockResponseData{}

	reader := io.NopCloser(bytes.NewReader(
		[]byte(`{"data":"test","name":"test"}`),
	))

	decodeErr := helpers.DecodeResponseBody(responseData, reader)

	if decodeErr != nil {
		t.Error(decodeErr)
	}

	if responseData.Data != "test" {
		t.Error("DecodeResponseBody failed: data field mismatch")
	}

	if responseData.Name != "test" {
		t.Error("DecodeResponseBody failed: name field mismatch")
	}
}
