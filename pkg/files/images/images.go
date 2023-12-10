package images

import (
	"bytes"
	"image"
	"image/gif"
	"image/jpeg"
	"image/png"
	"strings"

	"github.com/chai2010/webp"
	"golang.org/x/image/bmp"
	"golang.org/x/image/tiff"
)

const (
	PNG  = "png"
	JPEG = "jpeg"
	JPG  = "jpg"
	GIF  = "gif"
	WEBP = "webp"
	TIFF = "tiff"
	BMP  = "bmp"

	imageMimeType = "image/"
)

func toPNG(img image.Image) ([]byte, error) {
	buf := new(bytes.Buffer)

	// encode the image as a PNG file.
	if err := png.Encode(buf, img); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func toGIF(img image.Image) ([]byte, error) {
	buf := new(bytes.Buffer)

	// encode the image as a GIF file.
	if err := gif.Encode(buf, img, nil); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func toJPG(img image.Image) ([]byte, error) {
	buf := new(bytes.Buffer)

	// encode the image as a JPEG file.
	if err := jpeg.Encode(buf, img, nil); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func toWEBP(img image.Image) ([]byte, error) {
	buf := new(bytes.Buffer)

	// encode the image as a WEPB file.
	if err := webp.Encode(buf, img, nil); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func toTIFF(img image.Image) ([]byte, error) {
	buf := new(bytes.Buffer)

	// encode the image as a TIFF file.
	if err := tiff.Encode(buf, img, nil); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func toBMP(img image.Image) ([]byte, error) {
	buf := new(bytes.Buffer)

	if err := bmp.Encode(buf, img); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func ParseMimeType(mimetype string) string {
	if !strings.Contains(mimetype, imageMimeType) {
		return mimetype
	}

	return strings.TrimPrefix(mimetype, imageMimeType)
}
