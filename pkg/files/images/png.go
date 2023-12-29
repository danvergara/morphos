package images

import (
	"bytes"
	"fmt"
	"image/png"
	"slices"
	"strings"
)

// Png struct implements the File and Image interface from the files pkg.
type Png struct {
	compatibleFormats map[string][]string
}

// NewPng returns a pointer to a Png instance.
// The Png object is set with a map with list of supported file formats.
func NewPng() *Png {
	p := Png{
		compatibleFormats: map[string][]string{
			"Image": {
				AVIF,
				JPG,
				JPEG,
				GIF,
				WEBP,
				TIFF,
				BMP,
			},
		},
	}

	return &p
}

// SupportedFormats returns a map with a slice of supported files.
// Every key of the map represents the kind of a file.
func (p *Png) SupportedFormats() map[string][]string {
	return p.compatibleFormats
}

// ConvertTo method converts a given file to a target format.
// This method returns a file in form of a slice of bytes.
func (p *Png) ConvertTo(fileType, subType string, fileBytes []byte) ([]byte, error) {
	var result []byte

	compatibleFormats, ok := p.SupportedFormats()[fileType]
	if !ok {
		return nil, fmt.Errorf("ConvertTo: file type not supported: %s", fileType)
	}

	if !slices.Contains(compatibleFormats, subType) {
		return nil, fmt.Errorf("ConvertTo: file sub-type not supported: %s", subType)
	}

	switch strings.ToLower(fileType) {
	case imageType:
		img, err := png.Decode(bytes.NewReader(fileBytes))
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
func (p *Png) ImageType() string {
	return PNG
}
