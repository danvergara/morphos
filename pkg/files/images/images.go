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
	"github.com/signintech/gopdf"
	"golang.org/x/image/bmp"
	"golang.org/x/image/tiff"
)

const (
	// Images.
	PNG  = "png"
	JPEG = "jpeg"
	JPG  = "jpg"
	GIF  = "gif"
	WEBP = "webp"
	TIFF = "tiff"
	BMP  = "bmp"

	imageMimeType = "image/"
	imageType     = "image"

	// Documents.
	PDF = "pdf"

	documentMimeType = "application/"
	documentType     = "document"
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

func toPDF(img image.Image) ([]byte, error) {
	pdf := gopdf.GoPdf{}
	pdf.Start(
		gopdf.Config{
			PageSize: gopdf.Rect{W: float64(img.Bounds().Dx()), H: float64(img.Bounds().Dy())},
		},
	)

	// Add a page to the PDF
	pdf.AddPage()

	if err := pdf.ImageFrom(img, 0, 0, &gopdf.Rect{W: float64(img.Bounds().Dx()), H: float64(img.Bounds().Dy())}); err != nil {
		return nil, err
	}

	buf := new(bytes.Buffer)
	if _, err := pdf.WriteTo(buf); err != nil {
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

func convertToImage(target string, img image.Image) ([]byte, error) {
	var err error
	var result []byte

	switch target {
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
	case BMP:
		result, err = toBMP(img)
		if err != nil {
			return nil, err
		}
	default:
		return nil, fmt.Errorf("file format to convert to not supported: %s", target)
	}

	return result, nil
}

func convertToDocument(target string, img image.Image) ([]byte, error) {
	var err error
	var result []byte

	switch target {
	case PDF:
		result, err = toPDF(img)
		if err != nil {
			return nil, err
		}
	}

	return result, nil
}
