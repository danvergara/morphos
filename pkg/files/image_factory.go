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
		return new(images.Png), nil
	case images.JPEG:
		return new(images.Jpeg), nil
	case images.GIF:
		return new(images.Gif), nil
	case images.WEBP:
		return new(images.Webp), nil
	case images.TIFF:
		return new(images.Tiff), nil
	case images.BMP:
		return new(images.Bmp), nil
	default:
		return nil, fmt.Errorf("file of type %s not recognized\n", f)
	}
}
