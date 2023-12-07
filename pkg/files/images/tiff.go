package images

import (
	"bytes"
	"fmt"

	"golang.org/x/image/tiff"
)

type Tiff struct{}

func (t *Tiff) SupportedFormats() map[string][]string {
	return map[string][]string{
		"Image": {
			JPG,
			PNG,
			GIF,
			WEBP,
			BMP,
		},
	}
}

func (t *Tiff) ConvertTo(format string, fileBytes []byte) ([]byte, error) {
	var result []byte

	img, err := tiff.Decode(bytes.NewReader(fileBytes))
	if err != nil {
		return nil, err
	}

	switch format {
	case JPEG, JPG:
		result, err = toJPG(img)
		if err != nil {
			return nil, err
		}
	case PNG:
		result, err = toPNG(img)
		if err != nil {
			return nil, err
		}
	case WEBP:
		result, err = toWEBP(img)
		if err != nil {
			return nil, err
		}
	case GIF:
		result, err = toGIF(img)
		if err != nil {
			return nil, err
		}
	case BMP:
		result, err = toBMP(img)
		if err != nil {
			return nil, err
		}
	default:
		return nil, fmt.Errorf("file format to conver to %s not supported", format)
	}

	return result, nil
}

func (t *Tiff) ImageType() string {
	return TIFF
}
