package ebooks

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"slices"
	"strings"

	"github.com/danvergara/morphos/pkg/files/documents"
	"github.com/danvergara/morphos/pkg/util"
)

// Epub struct implements the File and Document interface from the file package.
type Epub struct {
	filename            string
	compatibleFormats   map[string][]string
	compatibleMIMETypes map[string][]string
}

func NewEpub(filename string) *Epub {
	e := Epub{
		filename: filename,
		compatibleFormats: map[string][]string{
			"Document": {
				documents.PDF,
			},
			"Ebook": {
				MOBI,
			},
		},
		compatibleMIMETypes: map[string][]string{
			"Document": {
				documents.PDF,
			},
			"Ebook": {
				MobiMimeType,
			},
		},
	}

	return &e
}

// SupportedFormats returns a map witht the compatible formats that Pdf is
// compatible to be converted to.
func (e *Epub) SupportedFormats() map[string][]string {
	return e.compatibleFormats
}

// SupportedMIMETypes returns a map witht the compatible MIME types that Pdf is
// compatible to be converted to.
func (e *Epub) SupportedMIMETypes() map[string][]string {
	return e.compatibleMIMETypes
}

func (e *Epub) ConvertTo(fileType, subtype string, file io.Reader) (io.Reader, error) {
	// These are guard clauses that check if the target file type is valid.
	compatibleFormats, ok := e.SupportedFormats()[fileType]
	if !ok {
		return nil, fmt.Errorf("ConvertTo: file type not supported: %s", fileType)
	}

	if !slices.Contains(compatibleFormats, subtype) {
		return nil, fmt.Errorf("ConvertTo: file sub-type not supported: %s", subtype)
	}

	buf := new(bytes.Buffer)
	if _, err := buf.ReadFrom(file); err != nil {
		return nil, fmt.Errorf(
			"error getting the content of the pdf file in form of slice of bytes: %w",
			err,
		)
	}

	fileBytes := buf.Bytes()

	switch strings.ToLower(fileType) {
	case documentType:
		switch subtype {
		case PDF:
			return util.EbookConvert(e.filename, EPUB, PDF, fileBytes)
		}
	case ebookType:
		switch subtype {
		case MOBI:
			return util.EbookConvert(e.filename, EPUB, MOBI, fileBytes)
		}
	}

	return nil, errors.New("not implemented")
}

// EbookType returns the Ebook type which is Epub in this case.
func (e *Epub) EbookType() string {
	return EPUB
}
