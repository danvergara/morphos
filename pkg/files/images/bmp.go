package images

import (
	"bytes"
	"fmt"
	"slices"
	"strings"

	"golang.org/x/image/bmp"
)

// Bmp struct implements the File and Image interface from the files pkg.
type Bmp struct {
	compatibleFormats map[string][]string
}

// NewBmp returns a pointer to a Bmp instance.
// The Bmp object is set with a map with list of supported file formats.
func NewBmp() *Bmp {
	b := Bmp{
		compatibleFormats: map[string][]string{
			"Image": {
				AVIF,
				JPG,
				JPEG,
				PNG,
				GIF,
				TIFF,
				WEBP,
			},
		},
	}

	return &b
}

// SupportedFormats returns a map with a slice of supported files.
// Every key of the map represents a kind of a file.
func (b *Bmp) SupportedFormats() map[string][]string {
	return b.compatibleFormats
}

// ConvertTo method converts a given file to a target format.
// This method returns a file in form of a slice of bytes.
// The methd receives a file type and the sub-type of the target format and the file as array of bytes.
func (b *Bmp) ConvertTo(fileType, subType string, fileBytes []byte) ([]byte, error) {
	var result []byte

	compatibleFormats, ok := b.SupportedFormats()[fileType]
	if !ok {
		return nil, fmt.Errorf("ConvertTo: file type not supported: %s", fileType)
	}

	if !slices.Contains(compatibleFormats, subType) {
		return nil, fmt.Errorf("ConvertTo: file sub-type not supported: %s", subType)
	}

	switch strings.ToLower(fileType) {
	case imageType:
		img, err := bmp.Decode(bytes.NewReader(fileBytes))
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
func (b *Bmp) ImageType() string {
	return BMP
}
