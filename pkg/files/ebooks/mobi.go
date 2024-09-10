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

type Mobi struct {
	filename            string
	compatibleFormats   map[string][]string
	compatibleMIMETypes map[string][]string
}

func NewMobi(filename string) Mobi {
	m := Mobi{
		filename: filename,
		compatibleFormats: map[string][]string{
			"Document": {
				documents.PDF,
			},
			"Ebook": {
				EPUB,
			},
		},
		compatibleMIMETypes: map[string][]string{
			"Document": {
				documents.PDF,
			},
			"Ebook": {
				EpubMimeType,
			},
		},
	}

	return m
}

// SupportedFormats returns a map witht the compatible formats that MOBI is
// compatible to be converted to.
func (m Mobi) SupportedFormats() map[string][]string {
	return m.compatibleFormats
}

// SupportedMIMETypes returns a map witht the compatible MIME types that MOBI is
// compatible to be converted to.
func (m Mobi) SupportedMIMETypes() map[string][]string {
	return m.compatibleMIMETypes
}

func (m Mobi) ConvertTo(fileType, subtype string, file io.Reader) (io.Reader, error) {
	// These are guard clauses that check if the target file type is valid.
	compatibleFormats, ok := m.SupportedFormats()[fileType]
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
			return util.EbookConvert(m.filename, MOBI, PDF, fileBytes)
		}
	case ebookType:
		switch subtype {
		case EPUB:
			return util.EbookConvert(m.filename, MOBI, EPUB, fileBytes)
		}
	}

	return nil, errors.New("file format not implemented")
}

// EbookType returns the Ebook type which is MOBI in this case.
func (m Mobi) EbookType() string {
	return EPUB
}
