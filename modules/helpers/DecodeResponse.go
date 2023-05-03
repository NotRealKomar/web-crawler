package helpers

import (
	"bytes"
	"encoding/json"
	"io"
)

func DecodeResponseBody[T any](output *T, body io.ReadCloser) error {
	defer body.Close()

	buffer := new(bytes.Buffer)
	_, copyErr := io.Copy(buffer, body)

	if copyErr != nil {
		return copyErr
	}

	unmarshalErr := json.Unmarshal(buffer.Bytes(), output)

	if unmarshalErr != nil {
		return unmarshalErr
	}

	return nil
}
