package ebooks

import (
	"bytes"
	"io"
	"os"
	"testing"

	"github.com/gabriel-vasile/mimetype"
	"github.com/stretchr/testify/require"
)

type file interface {
	SupportedFormats() map[string][]string
	ConvertTo(string, string, io.Reader) (io.Reader, error)
}

type ebook interface {
	file
	EbookType() string
}

func TestEbookTConvertTo(t *testing.T) {
	type input struct {
		filename       string
		mimetype       string
		targetFileType string
		targetFormat   string
		ebook          ebook
	}
	type expected struct {
		mimetype string
	}

	var tests = []struct {
		name     string
		input    input
		expected expected
	}{
		{
			name: "epub to pdf",
			input: input{
				filename:       "testdata/no-man-s-land.epub",
				mimetype:       "application/epub+zip",
				targetFileType: "Document",
				targetFormat:   "pdf",
				ebook:          NewEpub("no-man-s-land.epub"),
			},
			expected: expected{
				mimetype: "application/zip",
			},
		},
		{
			name: "epub to mobi",
			input: input{
				filename:       "testdata/no-man-s-land.epub",
				mimetype:       "application/epub+zip",
				targetFileType: "Ebook",
				targetFormat:   "mobi",
				ebook:          NewEpub("no-man-s-land.epub"),
			},
			expected: expected{
				mimetype: "application/zip",
			},
		},
	}
	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {

			inputDoc, err := os.ReadFile(tc.input.filename)
			require.NoError(t, err)

			detectedFileType := mimetype.Detect(inputDoc)
			require.Equal(t, tc.input.mimetype, detectedFileType.String())

			outoutFile, err := tc.input.ebook.ConvertTo(
				tc.input.targetFileType,
				tc.input.targetFormat,
				bytes.NewReader(inputDoc),
			)
			require.NoError(t, err)

			buf := new(bytes.Buffer)
			_, err = buf.ReadFrom(outoutFile)
			require.NoError(t, err)

			outoutFileBytes := buf.Bytes()
			detectedFileType = mimetype.Detect(outoutFileBytes)
			require.Equal(t, tc.expected.mimetype, detectedFileType.String())
		})
	}
}
