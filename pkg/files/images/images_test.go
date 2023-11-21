package images_test

import (
	"os"
	"testing"

	"github.com/danvergara/morphos/pkg/files/images"
	"github.com/gabriel-vasile/mimetype"
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
		{
			name: "webp to png",
			input: input{
				filename:     "testdata/gopher.webp",
				mimetype:     "image/webp",
				targetFormat: "png",
			},
			expected: expected{
				mimetype: "image/png",
			},
		},
		{
			name: "png to webp",
			input: input{
				filename:     "testdata/gopher_pirate.png",
				mimetype:     "image/png",
				targetFormat: "webp",
			},
			expected: expected{
				mimetype: "image/webp",
			},
		},
		{
			name: "webp to tiff",
			input: input{
				filename:     "testdata/gopher.webp",
				mimetype:     "image/webp",
				targetFormat: "tiff",
			},
			expected: expected{
				mimetype: "image/tiff",
			},
		},
		{
			name: "bmp to png",
			input: input{
				filename:     "testdata/sunset.bmp",
				mimetype:     "image/bmp",
				targetFormat: "png",
			},
			expected: expected{
				mimetype: "image/png",
			},
		},
		{
			name: "jpg to bmp",
			input: input{
				filename:     "testdata/Golang_Gopher.jpg",
				mimetype:     "image/jpeg",
				targetFormat: "bmp",
			},
			expected: expected{
				mimetype: "image/bmp",
			},
		},
	}
	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			inputImg, err := os.ReadFile(tc.input.filename)
			require.NoError(t, err)

			detectedFileType := mimetype.Detect(inputImg)
			require.Equal(t, tc.input.mimetype, detectedFileType.String())

			convertedImg, err := images.ConverImage(detectedFileType.String(), tc.input.targetFormat, inputImg)
			require.NoError(t, err)

			detectedFileType = mimetype.Detect(convertedImg)
			require.Equal(t, tc.expected.mimetype, detectedFileType.String())
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
					{Name: images.WEBP},
					{Name: images.TIFF},
					{Name: images.BMP},
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
					{Name: images.WEBP},
					{Name: images.TIFF},
					{Name: images.BMP},
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
					{Name: images.WEBP},
					{Name: images.TIFF},
					{Name: images.BMP},
				},
			},
		},
		{
			name: "WEBP",
			input: input{
				format: images.WEBP,
			},
			expected: expected{
				targetFormats: []images.FileFormat{
					{Name: images.JPG},
					{Name: images.PNG},
					{Name: images.GIF},
					{Name: images.TIFF},
					{Name: images.BMP},
				},
			},
		},
		{
			name: "TIFF",
			input: input{
				format: images.TIFF,
			},
			expected: expected{
				targetFormats: []images.FileFormat{
					{Name: images.JPG},
					{Name: images.PNG},
					{Name: images.GIF},
					{Name: images.WEBP},
					{Name: images.BMP},
				},
			},
		},
		{
			name: "BMP",
			input: input{
				format: images.BMP,
			},
			expected: expected{
				targetFormats: []images.FileFormat{
					{Name: images.JPG},
					{Name: images.PNG},
					{Name: images.GIF},
					{Name: images.WEBP},
					{Name: images.TIFF},
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
