package images

import (
	"bytes"
	"fmt"
	"image"
	"slices"
	"strings"

	_ "github.com/strukturag/libheif/go/heif"
)

// Avif struct implements the File and Image interface from the files pkg.
type Avif struct {
	compatibleFormats map[string][]string
}

// NewAvif returns a pointer to a Avif instance.
// The Avif object is set with a map with list of supported file formats.
func NewAvif() *Avif {
	a := Avif{
		compatibleFormats: map[string][]string{
			"Image": {
				JPG,
				JPEG,
				PNG,
				GIF,
				WEBP,
				TIFF,
				BMP,
			},
		},
	}

	return &a
}

// SupportedFormats returns a map with a slice of supported files.
// Every key of the map represents the kind of a file.
func (a *Avif) SupportedFormats() map[string][]string {
	return a.compatibleFormats
}

// ConvertTo method converts a given file to a target format.
// This method returns a file in form of a slice of bytes.
func (a *Avif) ConvertTo(fileType, subType string, fileBytes []byte) ([]byte, error) {
	var result []byte

	compatibleFormats, ok := a.SupportedFormats()[fileType]
	if !ok {
		return nil, fmt.Errorf("ConvertTo: file type not supported: %s", fileType)
	}

	if !slices.Contains(compatibleFormats, subType) {
		return nil, fmt.Errorf("ConvertTo: file sub-type not supported: %s", subType)
	}

	switch strings.ToLower(fileType) {
	case imageType:
		img, _, err := image.Decode(bytes.NewReader(fileBytes))
		if err != nil {
			return nil, fmt.Errorf(
				"ConvertTo: error at decoding avif image: %w",
				err,
			)
		}

		result, err = convertToImage(subType, img)
		if err != nil {
			return nil, fmt.Errorf(
				"ConvertTo: error at converting image to another format: %w",
				err,
			)
		}
	}

	return result, nil
}

// ImageType returns the file format of the current image.
// This method implements the Image interface.
func (a *Avif) ImageType() string {
	return AVIF
}
