package images

import (
	"bytes"
	"fmt"
	"image/gif"
	"slices"
	"strings"
)

// Gif struct implements the File and Image interface from the files pkg.
type Gif struct {
	compatibleFormats   map[string][]string
	compatibleMIMETypes map[string][]string
}

// NewGif returns a pointer to a Gif instance.
// The Gif object is set with a map with list of supported file formats.
func NewGif() *Gif {
	g := Gif{
		compatibleFormats: map[string][]string{
			"Image": {
				JPG,
				JPEG,
				PNG,
				WEBP,
				TIFF,
				BMP,
			},
			"Document": {
				PDF,
			},
		},
		compatibleMIMETypes: map[string][]string{
			"Image": {
				JPG,
				JPEG,
				PNG,
				WEBP,
				TIFF,
				BMP,
			},
			"Document": {
				PDF,
			},
		},
	}

	return &g
}

// SupportedFormats returns a map with a slice of supported files.
// Every key of the map represents the kind of a file.
func (g *Gif) SupportedFormats() map[string][]string {
	return g.compatibleFormats
}

// SupportedMIMETypes returns a map with a slice of supported MIME types.
func (g *Gif) SupportedMIMETypes() map[string][]string {
	return g.compatibleMIMETypes
}

// ConvertTo method converts a given file to a target format.
// This method returns a file in form of a slice of bytes.
// The methd receives a file type and the sub-type of the target format and the file as array of bytes.
func (g *Gif) ConvertTo(fileType, subType string, fileBytes []byte) (result []byte, err error) {
	compatibleFormats, ok := g.SupportedFormats()[fileType]
	if !ok {
		return nil, fmt.Errorf("ConvertTo: file type not supported: %s", fileType)
	}

	if !slices.Contains(compatibleFormats, subType) {
		return nil, fmt.Errorf("ConvertTo: file sub-type not supported: %s", subType)
	}

	switch strings.ToLower(fileType) {
	case imageType:
		result, err = convertToImage(g.ImageType(), subType, fileBytes)
		if err != nil {
			err = fmt.Errorf("ConvertTo: error at converting image to another format: %w", err)
		}
	case documentType:
		img, err := gif.Decode(bytes.NewReader(fileBytes))
		if err != nil {
			return nil, err
		}

		result, err = convertToDocument(subType, img)
		if err != nil {
			return nil, fmt.Errorf(
				"ConvertTo: error at converting image to another format: %w",
				err,
			)
		}
	}

	return
}

// ImageType returns the file format of the current image.
// This method implements the Image interface.
func (g *Gif) ImageType() string {
	return GIF
}
