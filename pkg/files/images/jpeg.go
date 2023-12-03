package images

import "errors"

type Jpeg struct{}

func (p *Jpeg) SupportedFormats() map[string]string {
	return make(map[string]string)
}

func (p *Jpeg) ConvertTo(format string) ([]byte, error) {
	return nil, errors.New("not implemented")
}

func (p *Jpeg) ImageType() string {
	return JPEG
}
