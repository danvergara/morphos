package files

import "fmt"

type FileFactory interface {
	NewFile(string) (File, error)
}

const (
	Img = "image"
	Doc = "document"
)

func BuildFactory(f string) (FileFactory, error) {
	switch f {
	case Img:
		return new(ImageFactory), nil
	case Doc:
		return new(DocumentFactory), nil
	default:
		return nil, fmt.Errorf("factory with id %s not recognized", f)
	}
}
