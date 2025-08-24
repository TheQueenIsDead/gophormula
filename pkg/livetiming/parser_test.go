package livetiming

import (
	"bytes"
	"compress/flate"
	"encoding/base64"
	"testing"
)

func TestDecompress(t *testing.T) {
	t.Run("it should decompress a valid payload", func(t *testing.T) {

		payload := `{"key":"value"}`

		// Compress and encode the payload to simulate the input for Decompress
		var b bytes.Buffer
		w, _ := flate.NewWriter(&b, flate.DefaultCompression)
		w.Write([]byte(payload))
		w.Close()
		encoded := base64.StdEncoding.EncodeToString(b.Bytes())

		decompressed, err := Decompress([]byte(encoded))
		if err != nil {
			t.Fatalf("Decompress failed with error: %v", err)
		}

		// Check if the output matches the original payload
		if string(decompressed) != payload {
			t.Errorf("expected %s, got %s", payload, string(decompressed))
		}
	})

	t.Run("it should return an error for invalid base64 data", func(t *testing.T) {
		invalidData := "this is not valid base64"
		_, err := Decompress([]byte(invalidData))
		if err == nil {
			t.Error("expected an error for invalid base64 data, but got nil")
		}
	})
}
