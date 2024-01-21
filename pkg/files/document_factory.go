package files

import (
	"fmt"

	"github.com/danvergara/morphos/pkg/files/documents"
)

// DocumentFactory implements the FileFactory interface.
type DocumentFactory struct {
	filename string
}

func NewDocumentFactory(filename string) *DocumentFactory {
	return &DocumentFactory{filename: filename}
}

// NewFile method returns an object that implements the File interface,
// given a document format as input.
// If not supported, it will error out.
func (d *DocumentFactory) NewFile(f string) (File, error) {
	switch f {
	case documents.PDF:
		return documents.NewPdf(d.filename), nil
	case documents.DOCX:
		return new(documents.Docx), nil
	default:
		return nil, fmt.Errorf("type file file  %s not recognized", f)
	}
}
