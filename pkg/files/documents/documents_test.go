package documents_test

import (
	"bytes"
	"io"
	"os"
	"testing"

	"github.com/gabriel-vasile/mimetype"
	"github.com/stretchr/testify/require"

	"github.com/danvergara/morphos/pkg/files/documents"
)

type filer interface {
	SupportedFormats() map[string][]string
	ConvertTo(string, string, io.Reader) (io.Reader, error)
}

type documenter interface {
	filer
	DocumentType() string
}

func TestPDFTConvertTo(t *testing.T) {
	type input struct {
		filename       string
		mimetype       string
		targetFileType string
		targetFormat   string
		documenter     documenter
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
			name: "pdf to mobi",
			input: input{
				filename:       "testdata/bitcoin.pdf",
				mimetype:       "application/pdf",
				targetFileType: "Ebook",
				targetFormat:   "mobi",
				documenter:     documents.NewPdf("bitcoin.pdf"),
			},
			expected: expected{
				mimetype: "application/zip",
			},
		},
		{
			name: "pdf to epub",
			input: input{
				filename:       "testdata/bitcoin.pdf",
				mimetype:       "application/pdf",
				targetFileType: "Ebook",
				targetFormat:   "epub",
				documenter:     documents.NewPdf("bitcoin.pdf"),
			},
			expected: expected{
				mimetype: "application/zip",
			},
		},
		{
			name: "pdf to jpeg",
			input: input{
				filename:       "testdata/bitcoin.pdf",
				mimetype:       "application/pdf",
				targetFileType: "Image",
				targetFormat:   "jpeg",
				documenter:     documents.NewPdf("bitcoin.pdf"),
			},
			expected: expected{
				mimetype: "application/zip",
			},
		},
		{
			name: "pdf to png",
			input: input{
				filename:       "testdata/bitcoin.pdf",
				mimetype:       "application/pdf",
				targetFileType: "Image",
				targetFormat:   "png",
				documenter:     documents.NewPdf("bitcoin.pdf"),
			},
			expected: expected{
				mimetype: "application/zip",
			},
		},
		{
			name: "pdf to gif",
			input: input{
				filename:       "testdata/bitcoin.pdf",
				mimetype:       "application/pdf",
				targetFileType: "Image",
				targetFormat:   "gif",
				documenter:     documents.NewPdf("bitcoin.pdf"),
			},
			expected: expected{
				mimetype: "application/zip",
			},
		},
		{
			name: "pdf to webp",
			input: input{
				filename:       "testdata/bitcoin.pdf",
				mimetype:       "application/pdf",
				targetFileType: "Image",
				targetFormat:   "webp",
				documenter:     documents.NewPdf("bitcoin.pdf"),
			},
			expected: expected{
				mimetype: "application/zip",
			},
		},
		{
			name: "pdf to docx",
			input: input{
				filename:       "testdata/bitcoin.pdf",
				mimetype:       "application/pdf",
				targetFileType: "Document",
				targetFormat:   "docx",
				documenter:     documents.NewPdf("bitcoin.pdf"),
			},
			expected: expected{
				mimetype: "application/zip",
			},
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			inputDoc, err := os.ReadFile(tc.input.filename)
			require.NoError(t, err)

			detectedFileType := mimetype.Detect(inputDoc)
			require.Equal(t, tc.input.mimetype, detectedFileType.String())

			outoutFile, err := tc.input.documenter.ConvertTo(
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

func TestDOCXTConvertTo(t *testing.T) {
	type input struct {
		filename       string
		mimetype       string
		targetFileType string
		targetFormat   string
		documenter     documenter
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
			name: "docx to pdf",
			input: input{
				filename:       "testdata/file_sample.docx",
				mimetype:       "application/vnd.openxmlformats-officedocument.wordprocessingml.document",
				targetFileType: "Document",
				targetFormat:   "pdf",
				documenter:     documents.NewDocx("file_sample.docx"),
			},
			expected: expected{
				mimetype: "application/zip",
			},
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			inputDoc, err := os.ReadFile(tc.input.filename)
			require.NoError(t, err)

			detectedFileType := mimetype.Detect(inputDoc)
			require.Equal(t, tc.input.mimetype, detectedFileType.String())

			resultFile, err := tc.input.documenter.ConvertTo(
				tc.input.targetFileType,
				tc.input.targetFormat,
				bytes.NewReader(inputDoc),
			)

			require.NoError(t, err)

			buf := new(bytes.Buffer)
			_, err = buf.ReadFrom(resultFile)
			require.NoError(t, err)

			resultFileBytes := buf.Bytes()
			detectedFileType = mimetype.Detect(resultFileBytes)
			require.Equal(t, tc.expected.mimetype, detectedFileType.String())
		})
	}
}

func TestCSVTConvertTo(t *testing.T) {
	type input struct {
		filename       string
		mimetype       string
		targetFileType string
		targetFormat   string
		documenter     documenter
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
			name: "csv to xlsx",
			input: input{
				filename:       "testdata/student.csv",
				mimetype:       "text/csv",
				targetFileType: "Document",
				targetFormat:   "xlsx",
				documenter:     documents.NewCsv("student.csv"),
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

			resultFile, err := tc.input.documenter.ConvertTo(
				tc.input.targetFileType,
				tc.input.targetFormat,
				bytes.NewReader(inputDoc),
			)

			require.NoError(t, err)

			buf := new(bytes.Buffer)
			_, err = buf.ReadFrom(resultFile)
			require.NoError(t, err)

			resultFileBytes := buf.Bytes()
			detectedFileType = mimetype.Detect(resultFileBytes)
			require.Equal(t, tc.expected.mimetype, detectedFileType.String())
		})
	}
}

func TestXLSXTConvertTo(t *testing.T) {
	type input struct {
		filename       string
		mimetype       string
		targetFileType string
		targetFormat   string
		documenter     documenter
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
			name: "xlsx to csv",
			input: input{
				filename:       "testdata/movies.xlsx",
				mimetype:       "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet",
				targetFileType: "Document",
				targetFormat:   "csv",
				documenter:     documents.NewXlsx("movies.xlsx"),
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

			resultFile, err := tc.input.documenter.ConvertTo(
				tc.input.targetFileType,
				tc.input.targetFormat,
				bytes.NewReader(inputDoc),
			)

			require.NoError(t, err)

			buf := new(bytes.Buffer)
			_, err = buf.ReadFrom(resultFile)
			require.NoError(t, err)

			resultFileBytes := buf.Bytes()
			detectedFileType = mimetype.Detect(resultFileBytes)
			require.Equal(t, tc.expected.mimetype, detectedFileType.String())
		})
	}
}
