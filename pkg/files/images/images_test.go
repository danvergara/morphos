package images_test

import (
	"net/http"
	"os"
	"testing"

	"github.com/danvergara/morphos/pkg/files/images"
	"github.com/stretchr/testify/require"
)

func TestConvertImage(t *testing.T) {
	type input struct {
		filename     string
		mimetype     string
		targetFormat string
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
			name: "png to jpeg",
			input: input{
				filename:     "testdata/gopher_pirate.png",
				mimetype:     "image/png",
				targetFormat: "jpeg",
			},
			expected: expected{
				mimetype: "image/jpeg",
			},
		},
		{
			name: "jpeg to png",
			input: input{
				filename:     "testdata/Golang_Gopher.jpg",
				mimetype:     "image/jpeg",
				targetFormat: "png",
			},
			expected: expected{
				mimetype: "image/png",
			},
		},
	}
	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			inputImg, err := os.ReadFile(tc.input.filename)
			require.NoError(t, err)

			detectedFileType := http.DetectContentType(inputImg)
			require.Equal(t, tc.input.mimetype, detectedFileType)

			convertedImg, err := images.ConverImage(detectedFileType, tc.input.targetFormat, inputImg)
			require.NoError(t, err)

			detectedFileType = http.DetectContentType(convertedImg)
			require.Equal(t, tc.expected.mimetype, detectedFileType)
		})
	}

}

func TestFileFormatsToConvert(t *testing.T) {
	type input struct {
		format string
	}
	type expected struct {
		targetFormats []images.FileFormat
	}

	var tests = []struct {
		name     string
		input    input
		expected expected
	}{
		{
			name: "JPEG",
			input: input{
				format: images.JPEG,
			},
			expected: expected{
				targetFormats: []images.FileFormat{
					{Name: images.PNG},
					{Name: images.GIF},
				},
			},
		},
		{
			name: "PNG",
			input: input{
				format: images.PNG,
			},
			expected: expected{
				targetFormats: []images.FileFormat{
					{Name: images.JPG},
					{Name: images.GIF},
				},
			},
		},
		{
			name: "GIF",
			input: input{
				format: images.GIF,
			},
			expected: expected{
				targetFormats: []images.FileFormat{
					{Name: images.JPG},
					{Name: images.PNG},
				},
			},
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			formats := images.FileFormatsToConvert(tc.input.format)
			require.EqualValues(t, tc.expected.targetFormats, formats["Formats"])
		})
	}
}

func TestParseMimeType(t *testing.T) {
	parsedType := images.ParseMimeType("image/png")
	require.Equal(t, parsedType, "png")
}
