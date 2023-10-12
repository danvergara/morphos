package image_test

import (
	"net/http"
	"os"
	"testing"

	"github.com/danvergara/morphos/pkg/image"
	"github.com/stretchr/testify/require"
)

func TestPngToJpeg(t *testing.T) {
	imagePng, err := os.ReadFile("testdata/gopher_pirate.png")
	require.NoError(t, err)

	imageJpp, err := image.PngToJpeg(imagePng)
	require.NoError(t, err)

	detectedFileType := http.DetectContentType(imageJpp)
	require.Equal(t, "image/jpeg", detectedFileType)
}

func TestJpegToPng(t *testing.T) {
	imagePng, err := os.ReadFile("testdata/Golang_Gopher.jpg")
	require.NoError(t, err)

	imageJpp, err := image.JpegToPng(imagePng)
	require.NoError(t, err)

	detectedFileType := http.DetectContentType(imageJpp)
	require.Equal(t, "image/png", detectedFileType)
}
