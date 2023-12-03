package images

import "errors"

type Png struct{}

func (p *Png) SupportedFormats() map[string]string {
	return make(map[string]string)
}

func (p *Png) ConvertTo(format string) ([]byte, error) {
	return nil, errors.New("not implemented")
}

func (p *Png) ImageType() string {
	return PNG
}
