package files

import (
	"fmt"

	"github.com/danvergara/morphos/pkg/files/images"
)

// ImageFactory implements the FileFactory interface.
type ImageFactory struct{}

// NewFile method returns an object that implements the File interface,
// given an image format as input.
// If not supported, it will error out.
func (i *ImageFactory) NewFile(f string) (File, error) {
	switch f {
	case images.PNG:
		return images.NewPng(), nil
	case images.JPEG:
		return images.NewJpeg(), nil
	case images.GIF:
		return images.NewGif(), nil
	case images.WEBP:
		return images.NewWebp(), nil
	case images.TIFF:
		return images.NewTiff(), nil
	case images.BMP:
		return images.NewBmp(), nil
	default:
		return nil, fmt.Errorf("type file %s not recognized", f)
	}
}
