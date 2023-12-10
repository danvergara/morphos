package images_test

import (
	"os"
	"testing"

	"github.com/gabriel-vasile/mimetype"
	"github.com/stretchr/testify/require"

	"github.com/danvergara/morphos/pkg/files/images"
)

type filer interface {
	SupportedFormats() map[string][]string
	ConvertTo(string, []byte) ([]byte, error)
}

type imager interface {
	filer
	ImageType() string
}

func TestConvertImage(t *testing.T) {
	type input struct {
		filename     string
		mimetype     string
		targetFormat string
		imager       imager
	}
	type expected struct {
		mimetype         string
		supportedFormats map[string][]string
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
				imager:       new(images.Png),
			},
			expected: expected{
				mimetype: "image/jpeg",
				supportedFormats: map[string][]string{
					"Image": {
						images.JPG,
						images.GIF,
						images.WEBP,
						images.TIFF,
						images.BMP,
					},
				},
			},
		},
		{
			name: "jpeg to png",
			input: input{
				filename:     "testdata/Golang_Gopher.jpg",
				mimetype:     "image/jpeg",
				targetFormat: "png",
				imager:       new(images.Jpeg),
			},
			expected: expected{
				mimetype: "image/png",
				supportedFormats: map[string][]string{
					"Image": {
						images.PNG,
						images.GIF,
						images.WEBP,
						images.TIFF,
						images.BMP,
					},
				},
			},
		},
		{
			name: "webp to png",
			input: input{
				filename:     "testdata/gopher.webp",
				mimetype:     "image/webp",
				targetFormat: "png",
				imager:       new(images.Webp),
			},
			expected: expected{
				mimetype: "image/png",
				supportedFormats: map[string][]string{
					"Image": {
						images.JPG,
						images.PNG,
						images.GIF,
						images.TIFF,
						images.BMP,
					},
				},
			},
		},
		{
			name: "png to webp",
			input: input{
				filename:     "testdata/gopher_pirate.png",
				mimetype:     "image/png",
				targetFormat: "webp",
				imager:       new(images.Png),
			},
			expected: expected{
				mimetype: "image/webp",
				supportedFormats: map[string][]string{
					"Image": {
						images.JPG,
						images.GIF,
						images.WEBP,
						images.TIFF,
						images.BMP,
					},
				},
			},
		},
		{
			name: "webp to tiff",
			input: input{
				filename:     "testdata/gopher.webp",
				mimetype:     "image/webp",
				targetFormat: "tiff",
				imager:       new(images.Webp),
			},
			expected: expected{
				mimetype: "image/tiff",
				supportedFormats: map[string][]string{
					"Image": {
						images.JPG,
						images.PNG,
						images.GIF,
						images.TIFF,
						images.BMP,
					},
				},
			},
		},
		{
			name: "bmp to png",
			input: input{
				filename:     "testdata/sunset.bmp",
				mimetype:     "image/bmp",
				targetFormat: "png",
				imager:       new(images.Bmp),
			},
			expected: expected{
				mimetype: "image/png",
				supportedFormats: map[string][]string{
					"Image": {
						images.JPG,
						images.PNG,
						images.GIF,
						images.TIFF,
						images.WEBP,
					},
				},
			},
		},
		{
			name: "jpg to bmp",
			input: input{
				filename:     "testdata/Golang_Gopher.jpg",
				mimetype:     "image/jpeg",
				targetFormat: "bmp",
				imager:       new(images.Jpeg),
			},
			expected: expected{
				mimetype: "image/bmp",
				supportedFormats: map[string][]string{
					"Image": {
						images.PNG,
						images.GIF,
						images.WEBP,
						images.TIFF,
						images.BMP,
					},
				},
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

			convertedImg, err := tc.input.imager.ConvertTo(
				tc.input.targetFormat,
				inputImg,
			)
			require.NoError(t, err)

			detectedFileType = mimetype.Detect(convertedImg)
			require.Equal(t, tc.expected.mimetype, detectedFileType.String())
			formats := tc.input.imager.SupportedFormats()
			require.EqualValues(t, tc.expected.supportedFormats, formats)
		})
	}
}

func TestParseMimeType(t *testing.T) {
	parsedType := images.ParseMimeType("image/png")
	require.Equal(t, parsedType, "png")
}
