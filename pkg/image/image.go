package image

import (
	"bytes"
	"image/jpeg"
	"image/png"
)

// PngToJpeg converts a PNG image to JPEG format.
func PngToJpeg(imageBytes []byte) ([]byte, error) {
	// Decode the PNG image bytes.
	img, err := png.Decode(bytes.NewReader(imageBytes))

	if err != nil {
		return nil, err
	}

	buf := new(bytes.Buffer)

	// encode the image as a JPEG file.
	if err := jpeg.Encode(buf, img, nil); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

// JpegToPng converts a JPEG image to PNG format.
func JpegToPng(imageBytes []byte) ([]byte, error) {
	img, err := jpeg.Decode(bytes.NewReader(imageBytes))

	if err != nil {
		return nil, err
	}

	buf := new(bytes.Buffer)

	if err := png.Encode(buf, img); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}
