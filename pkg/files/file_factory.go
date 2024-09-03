package files

import "fmt"

// FileFactory interface is responsible for defining how a FileFactory behaves.
// It defines a NewFile method that returns an entity
// that implements the File interface.
type FileFactory interface {
	NewFile(string) (File, error)
}

const (
	Img = "image"
	// Application is provide because the type from the document's mimetype
	// is defined as application, not document. Both are supported.
	Application = "application"
	Doc         = "document"
	Text        = "text"
	Ebook       = "ebook"
)

// BuildFactory is a function responsible to return a FileFactory,
// given a supported and valid file type, otherwise, it will error out.
func BuildFactory(f string, filename string) (FileFactory, error) {
	switch f {
	case Img:
		return new(ImageFactory), nil
	case Doc, Application, Text:
		return NewDocumentFactory(filename), nil
	default:
		return nil, fmt.Errorf("factory with type file %s not recognized", f)
	}
}
