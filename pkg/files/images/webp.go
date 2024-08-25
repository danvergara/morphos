package images

import (
	"bytes"
	"fmt"
	"image"
	"image/draw"
	"io"
	"slices"
	"strings"

	"golang.org/x/image/webp"
)

// Webp struct implements the File and Image interface from the files pkg.
type Webp struct {
	compatibleFormats   map[string][]string
	compatibleMIMETypes map[string][]string
}

// NewWebp returns a pointer to a Webp instance.
// The Webp object is set with a map with list of supported file formats.
func NewWebp() *Webp {
	w := Webp{
		compatibleFormats: map[string][]string{
			"Image": {
				AVIF,
				JPG,
				JPEG,
				PNG,
				GIF,
				TIFF,
				BMP,
			},
			"Document": {
				PDF,
			},
		},
		compatibleMIMETypes: map[string][]string{
			"Image": {
				AVIF,
				JPG,
				JPEG,
				PNG,
				GIF,
				TIFF,
				BMP,
			},
			"Document": {
				PDF,
			},
		},
	}

	return &w
}

// SupportedFormats returns a map with a slice of supported files.
// Every key of the map represents the kind of a file.
func (w *Webp) SupportedFormats() map[string][]string {
	return w.compatibleFormats
}

// SupportedMIMETypes returns a map with a slice of supported MIME types.
func (w *Webp) SupportedMIMETypes() map[string][]string {
	return w.compatibleMIMETypes
}

// ConvertTo method converts a given file to a target format.
// This method returns a file in form of a slice of bytes.
func (w *Webp) ConvertTo(fileType, subType string, file io.Reader) (io.Reader, error) {

	var result []byte

	compatibleFormats, ok := w.SupportedFormats()[fileType]
	if !ok {
		return nil, fmt.Errorf("ConvertTo: file type not supported: %s", fileType)
	}

	if !slices.Contains(compatibleFormats, subType) {
		return nil, fmt.Errorf("ConvertTo: file sub-type not supported: %s", subType)
	}

	switch strings.ToLower(fileType) {
	case imageType:
		convertedImage, err := convertToImage(subType, file)
		if err != nil {
			return nil, err
		}

		return convertedImage, nil
	case documentType:
		img, err := webp.Decode(file)
		if err != nil {
			return nil, err
		}

		rgba := image.NewRGBA(img.Bounds())
		draw.Draw(rgba, img.Bounds(), img, image.Point{}, draw.Src)

		result, err = convertToDocument(subType, rgba)
		if err != nil {
			return nil, fmt.Errorf(
				"ConvertTo: error at converting image to another format: %w",
				err,
			)
		}
	}

	return bytes.NewReader(result), nil
}

// ImageType method returns the file format of the current image.
// This method implements the Image interface.
func (w *Webp) ImageType() string {
	return WEBP
}
