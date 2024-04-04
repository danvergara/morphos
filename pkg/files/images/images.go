package images

import (
	"bytes"
	"fmt"
	"github.com/davidbyttow/govips/v2/vips"
	"image"
	"image/gif"
	"image/jpeg"
	"image/png"
	"strings"

	"github.com/signintech/gopdf"
	"golang.org/x/image/bmp"
	"golang.org/x/image/tiff"
	"golang.org/x/image/webp"
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

func toBMP(source string, img []byte) ([]byte, error) {
	var decodedImage image.Image
	var err error

	// decode from source
	switch source {
	case PNG:
		decodedImage, err = png.Decode(bytes.NewReader(img))
	case JPEG, JPG:
		decodedImage, err = jpeg.Decode(bytes.NewReader(img))
	case GIF:
		decodedImage, err = gif.Decode(bytes.NewReader(img))
	case WEBP:
		decodedImage, err = webp.Decode(bytes.NewReader(img))
	case TIFF:
		decodedImage, err = tiff.Decode(bytes.NewReader(img))
	case BMP:
		decodedImage, err = bmp.Decode(bytes.NewReader(img))
	}
	if err != nil {
		return nil, err
	}

	// encode to BMP
	buf := new(bytes.Buffer)

	if err = bmp.Encode(buf, decodedImage); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

// toPDF returns pdf file as an slice of bytes.
// Receives an image.Image as a parameter.
func toPDF(img image.Image) ([]byte, error) {
	// Sets a Rectangle based on the size of the image.
	imgRect := gopdf.Rect{
		W: float64(img.Bounds().Dx()),
		H: float64(img.Bounds().Dy()),
	}

	// Init the pdf obkect.
	pdf := gopdf.GoPdf{}

	// Sets the size of the every pdf page,
	// based on the dimensions of the image.
	pdf.Start(
		gopdf.Config{
			PageSize: imgRect,
		},
	)

	// Add a page to the PDF.
	pdf.AddPage()

	// Draws the image on the rectangle on the page above created.
	if err := pdf.ImageFrom(img, 0, 0, &imgRect); err != nil {
		return nil, err
	}

	// Creates a bytes.Buffer and writes the pdf data to it.
	buf := new(bytes.Buffer)
	if _, err := pdf.WriteTo(buf); err != nil {
		return nil, err
	}

	// Returns the pdf data as slice of bytes.
	return buf.Bytes(), nil
}

func ParseMimeType(mimetype string) string {
	if !strings.Contains(mimetype, imageMimeType) {
		return mimetype
	}

	return strings.TrimPrefix(mimetype, imageMimeType)
}

func convertToImage(source, target string, img []byte) ([]byte, error) {
	// special case BMP
	if target == BMP {
		return toBMP(source, img)
	}

	// create vips image
	vipsImage, err := vips.NewImageFromBuffer(img)
	if err != nil {
		return nil, err
	}
	defer vipsImage.Close()
	if err = vipsImage.AutoRotate(); err != nil {
		return nil, err
	}

	var data []byte

	switch target {
	case PNG:
		data, _, err = vipsImage.ExportPng(nil)
	case JPEG, JPG:
		data, _, err = vipsImage.ExportJpeg(nil)
	case GIF:
		data, _, err = vipsImage.ExportGIF(nil)
	case WEBP:
		data, _, err = vipsImage.ExportWebp(nil)
	case TIFF:
		data, _, err = vipsImage.ExportTiff(nil)
	default:
		return nil, fmt.Errorf("file format to convert to not supported: %s", target)
	}

	return data, err
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
