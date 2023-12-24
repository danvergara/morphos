package documents

import "errors"

type Docx struct{}

func (p *Docx) SupportedFormats() map[string][]string {
	return make(map[string][]string)
}

func (p *Docx) ConvertTo(fileType, subType string, fileBytes []byte) ([]byte, error) {
	return nil, errors.New("not implemented")
}

func (p *Docx) DocumentType() string {
	return DOCX
}
