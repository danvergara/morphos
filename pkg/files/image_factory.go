package files

import (
	"fmt"

	"github.com/danvergara/morphos/pkg/files/images"
)

type ImageFactory struct{}

func (i *ImageFactory) NewFile(f string) (File, error) {
	switch f {
	case images.PNG:
		return new(images.Png), nil
	case images.JPEG:
		return new(images.Jpeg), nil
	default:
		return nil, fmt.Errorf("file of type %s not recognized\n", f)
	}
}
