package documents

import "errors"

type Pdf struct{}

func (p *Pdf) SupportedFormats() map[string][]string {
	return make(map[string][]string)
}

func (p *Pdf) ConvertTo(fileType, subType string, fileBytes []byte) ([]byte, error) {
	return nil, errors.New("not implemented")
}

func (p *Pdf) DocumentType() string {
	return PDF
}
