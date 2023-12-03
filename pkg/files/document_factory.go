package files

import (
	"fmt"

	"github.com/danvergara/morphos/pkg/files/documents"
)

type DocumentFactory struct{}

func (d *DocumentFactory) NewFile(f string) (File, error) {
	switch f {
	case documents.PDF:
		return new(documents.Pdf), nil
	case documents.DOCX:
		return new(documents.Docx), nil
	default:
		return nil, fmt.Errorf("file of type %s not recognized", f)
	}
}
