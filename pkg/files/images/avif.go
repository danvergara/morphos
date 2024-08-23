package images

import (
	"fmt"
	"io"
	"slices"
	"strings"
)

// Avif struct implements the File and Image interface from the files pkg.
type Avif struct {
	compatibleFormats   map[string][]string
	compatibleMIMETypes map[string][]string
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

		compatibleMIMETypes: map[string][]string{
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

// SupportedMIMETypes returns a map with a slice of supported MIME types.
func (a *Avif) SupportedMIMETypes() map[string][]string {
	return a.compatibleMIMETypes
}

// ConvertTo method converts a given file to a target format.
// This method returns a file in form of a slice of bytes.
func (a *Avif) ConvertTo(fileType, subType string, file io.Reader) (io.Reader, error) {
	compatibleFormats, ok := a.SupportedFormats()[fileType]
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
	default:
		return nil, fmt.Errorf("not supported file type %s", fileType)
	}
}

// ImageType returns the file format of the current image.
// This method implements the Image interface.
func (a *Avif) ImageType() string {
	return AVIF
}
