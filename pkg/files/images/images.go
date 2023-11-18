package images

import (
	"bytes"
	"fmt"
	"image"
	"image/gif"
	"image/jpeg"
	"image/png"
	"strings"

	"github.com/chai2010/webp"
	"golang.org/x/image/tiff"
	webpx "golang.org/x/image/webp"
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

// FileFormat is the file format representation meant to be shown in the
// form template as an option.
type FileFormat struct {
	// Name of the file format to be shown in the option tag and as option value.
	Name string
}

// ConverImage returns a image converted as an array of bytes,
// if somethings wrong happens, the functions will error out.
// The functions receives the format from to be converted,
// the file format to be converted to and the image to be converted.
func ConverImage(from, to string, imageBytes []byte) ([]byte, error) {
	var (
		img    image.Image
		result []byte
		err    error
	)

	to = ParseMimeType(to)
	from = ParseMimeType(from)

	switch from {
	case PNG:
		img, err = png.Decode(bytes.NewReader(imageBytes))
		if err != nil {
			return nil, err
		}
	case JPEG, JPG:
		img, err = jpeg.Decode(bytes.NewReader(imageBytes))
		if err != nil {
			return nil, err
		}
	case GIF:
		img, err = gif.Decode(bytes.NewReader(imageBytes))
		if err != nil {
			return nil, err
		}
	case WEBP:
		img, err = webpx.Decode(bytes.NewReader(imageBytes))
	case TIFF:
		img, err = tiff.Decode(bytes.NewReader(imageBytes))
	default:
		return nil, fmt.Errorf("file format %s not supported", from)
	}

	switch to {
	case PNG:
		result, err = toPNG(img)
		if err != nil {
			return nil, err
		}
	case JPEG, JPG:
		result, err = toJPG(img)
		if err != nil {
			return nil, err
		}
	case GIF:
		result, err = toGIF(img)
		if err != nil {
			return nil, err
		}
	case WEBP:
		result, err = toWEBP(img)
		if err != nil {
			return nil, err
		}
	case TIFF:
		result, err = toTIFF(img)
		if err != nil {
			return nil, err
		}
	default:
		return nil, fmt.Errorf("file format to conver to %s not supported", to)
	}

	return result, nil
}

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

func FileFormatsToConvert(to string) map[string][]FileFormat {
	formats := make(map[string][]FileFormat)

	to = ParseMimeType(to)

	switch to {
	case JPEG, JPG:
		formats = map[string][]FileFormat{
			"Formats": {
				{Name: PNG},
				{Name: GIF},
				{Name: WEBP},
				{Name: TIFF},
			},
		}
	case PNG:
		formats = map[string][]FileFormat{
			"Formats": {
				{Name: JPG},
				{Name: GIF},
				{Name: WEBP},
				{Name: TIFF},
			},
		}
	case GIF:
		formats = map[string][]FileFormat{
			"Formats": {
				{Name: JPG},
				{Name: PNG},
				{Name: WEBP},
				{Name: TIFF},
			},
		}
	case WEBP:
		formats = map[string][]FileFormat{
			"Formats": {
				{Name: JPG},
				{Name: PNG},
				{Name: GIF},
				{Name: TIFF},
			},
		}
	case TIFF:
		formats = map[string][]FileFormat{
			"Formats": {
				{Name: JPG},
				{Name: PNG},
				{Name: GIF},
				{Name: WEBP},
			},
		}
	}

	return formats
}

func ParseMimeType(mimetype string) string {
	if !strings.Contains(mimetype, imageMimeType) {
		return mimetype
	}

	return strings.TrimPrefix(mimetype, imageMimeType)
}
