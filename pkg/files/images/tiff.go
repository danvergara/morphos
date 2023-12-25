package images

import (
	"bytes"
	"fmt"
	"slices"
	"strings"

	"golang.org/x/image/tiff"
)

// Tiff struct implements the File and Image interface from the files pkg.
type Tiff struct {
	compatibleFormats map[string][]string
}

// NewTiff returns a pointer to a Tiff instance.
// The Tiff object is set with a map with list of supported file formats.
func NewTiff() *Tiff {
	t := Tiff{
		compatibleFormats: map[string][]string{
			"Image": {
				AVIF,
				JPG,
				JPEG,
				PNG,
				GIF,
				WEBP,
				BMP,
			},
		},
	}

	return &t
}

// SupportedFormats method returns a map with a slice of supported files.
// Every key of the map represents the kind of a file.
func (t *Tiff) SupportedFormats() map[string][]string {
	return t.compatibleFormats
}

// ConvertTo method converts a given file to a target format.
// This method returns a file in form of a slice of bytes.
func (t *Tiff) ConvertTo(fileType, subType string, fileBytes []byte) ([]byte, error) {

	var result []byte

	compatibleFormats, ok := t.SupportedFormats()[fileType]
	if !ok {
		return nil, fmt.Errorf("ConvertTo: file type not supported: %s", fileType)
	}

	if !slices.Contains(compatibleFormats, subType) {
		return nil, fmt.Errorf("ConvertTo: file sub-type not supported: %s", subType)
	}

	switch strings.ToLower(fileType) {
	case imageType:
		img, err := tiff.Decode(bytes.NewReader(fileBytes))
		if err != nil {
			return nil, err
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
func (t *Tiff) ImageType() string {
	return TIFF
}
