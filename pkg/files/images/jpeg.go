package images

import (
	"bytes"
	"fmt"
	"image"
	"image/draw"
	"image/jpeg"
	"slices"
	"strings"
)

// Jpeg struct implements the File and Image interface from the files pkg.
type Jpeg struct {
	compatibleFormats   map[string][]string
	compatibleMIMETypes map[string][]string
}

// NewJpeg returns a pointer to a Jpeg instance.
// The Jpeg object is set with a map with list of supported file formats.
func NewJpeg() *Jpeg {
	j := Jpeg{
		compatibleFormats: map[string][]string{
			"Image": {
				PNG,
				GIF,
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
				PNG,
				GIF,
				WEBP,
				TIFF,
				BMP,
			},
			"Document": {
				PDF,
			},
		},
	}

	return &j
}

// SupportedFormats returns a map with a slice of supported files.
// Every key of the map represents a kind of a file.
func (j *Jpeg) SupportedFormats() map[string][]string {
	return j.compatibleFormats
}

// SupportedMIMETypes returns a map with a slice of supported MIME types.
func (j *Jpeg) SupportedMIMETypes() map[string][]string {
	return j.compatibleMIMETypes
}

// ConvertTo method converts a given file to a target format.
// This method returns a file in form of a slice of bytes.
// The methd receives a file type and the sub-type of the target format and the file as array of bytes.
func (j *Jpeg) ConvertTo(fileType, subType string, fileBytes []byte) ([]byte, error) {
	var result []byte

	compatibleFormats, ok := j.SupportedFormats()[fileType]
	if !ok {
		return nil, fmt.Errorf("ConvertTo: file type not supported: %s", fileType)
	}

	if !slices.Contains(compatibleFormats, subType) {
		return nil, fmt.Errorf("ConvertTo: file sub-type not supported: %s", subType)
	}

	switch strings.ToLower(fileType) {
	case imageType:
		img, err := jpeg.Decode(bytes.NewReader(fileBytes))
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
	case documentType:
		img, err := jpeg.Decode(bytes.NewReader(fileBytes))
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

	return result, nil
}

// ImageType returns the file format of the current image.
// This method implements the Image interface.
func (j *Jpeg) ImageType() string {
	return JPEG
}
