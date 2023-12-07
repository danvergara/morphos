package images

import (
	"bytes"
	"fmt"
	"image/png"
)

type Png struct{}

func (p *Png) SupportedFormats() map[string][]string {
	return map[string][]string{
		"Image": {
			JPG,
			GIF,
			WEBP,
			TIFF,
			BMP,
		},
	}
}

func (p *Png) ConvertTo(format string, fileBytes []byte) ([]byte, error) {
	var result []byte

	img, err := png.Decode(bytes.NewReader(fileBytes))
	if err != nil {
		return nil, err
	}

	switch format {
	case JPEG, JPG:
		result, err = toJPG(img)
		if err != nil {
			return nil, err
		}
	case GIF:
		result, err = toGIF(img)
		if err != nil {
			return nil, err
		}
	case WEBP:
		result, err = toWEBP(img)
		if err != nil {
			return nil, err
		}
	case TIFF:
		result, err = toTIFF(img)
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

func (p *Png) ImageType() string {
	return PNG
}
