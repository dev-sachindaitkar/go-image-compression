package compression

import (
	"bytes"
	"fmt"
	"image"
	"image/jpeg"
	_ "image/png" // Blank import registers the PNG decoder framework automatically
)

// CompressConfig holds the parameter configurations coming from the UI sliders
type CompressConfig struct {
	Quality int // Value from 1 to 100 (Higher means better quality, larger size)
}

// CompressImage converts raw input bytes into an optimized JPEG format using standard lib buffers
func CompressImage(inputData []byte, config CompressConfig) ([]byte, error) {
	// 1. Convert the raw byte slice into an io.Reader stream
	reader := bytes.NewReader(inputData)

	// 2. Decode the image. This automatically routes to the JPEG or PNG decoder
	// based on the image's magic bytes header.
	img, format, err := image.Decode(reader)
	if err != nil {
		return nil, fmt.Errorf("failed to decode image structure: %w", err)
	}

	// 3. Set up an in-memory byte buffer to hold our newly compressed binary data
	var outputBuffer bytes.Buffer

	// 4. Execute encoding configuration tuning
	switch format {
	case "jpeg", "jpg", "png":
		// To optimize file size natively, we transcode all incoming files (including PNGs)
		// into compressed JPEGs with fine-tuned baseline optimization.
		options := &jpeg.Options{
			Quality: config.Quality,
		}

		err = jpeg.Encode(&outputBuffer, img, options)
		if err != nil {
			return nil, fmt.Errorf("failed to encode compressed jpeg bytes: %w", err)
		}

	default:
		return nil, fmt.Errorf("unsupported file format format: %s", format)
	}

	// 5. Extract the optimized byte payload
	return outputBuffer.Bytes(), nil
}
