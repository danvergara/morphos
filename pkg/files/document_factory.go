package files

import (
	"fmt"

	"github.com/danvergara/morphos/pkg/files/documents"
)

// DocumentFactory implements the FileFactory interface.
type DocumentFactory struct{}

// NewFile method returns an object that implements the File interface,
// given a document format as input.
// If not supported, it will error out.
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
