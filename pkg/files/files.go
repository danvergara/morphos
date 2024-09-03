package files

import "io"

// File interface is the main interface of the package,
// that defines what a file is in this context.
// It's moslty responsible to say other entitites what formats it can be converted to
// and provides a method to convert the current file given a target format, if supported.
// SupportedMIMETypes was added to tell between how we see files, categorized by
// extension, and how they are registered as MIME types.
// e.g.
// Kind of document: 	Microsoft Word (OpenXML)
// Extension: docx
// MIME Type: application/vnd.openxmlformats-officedocument.wordprocessingml.document
type File interface {
	SupportedFormats() map[string][]string
	SupportedMIMETypes() map[string][]string
	ConvertTo(string, string, io.Reader) (io.Reader, error)
}

// SupportedFileTypes returns a map with the underlying file type,
// given a sub-type.
func SupportedFileTypes() map[string]string {
	return map[string]string{
		"avif": "image",
		"png":  "image",
		"jpg":  "image",
		"jpeg": "image",
		"gif":  "image",
		"webp": "image",
		"tiff": "image",
		"bmp":  "image",
		"docx": "document",
		"pdf":  "document",
		"xlsx": "document",
		"csv":  "document",
		"epub": "ebook",
	}
}
