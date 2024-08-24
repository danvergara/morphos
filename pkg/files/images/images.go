package images

import (
	"bytes"
	"fmt"
	"image"
	"io"
	"math/rand"
	"os"
	"strings"
	"time"

	"github.com/signintech/gopdf"
	ffmpeg "github.com/u2takey/ffmpeg-go"
)

const (
	// Images.
	PNG  = "png"
	JPEG = "jpeg"
	JPG  = "jpg"
	GIF  = "gif"
	WEBP = "webp"
	TIFF = "tiff"
	BMP  = "bmp"
	AVIF = "avif"

	imageMimeType = "image/"
	imageType     = "image"

	// Documents.
	PDF = "pdf"

	documentMimeType = "application/"
	documentType     = "document"

	// letters is constant that used as a pool of letters to generate a random string.
	letters = "abcdefghijklmnopqrstuvwxyz" +
		"ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
)

var seededRand *rand.Rand = rand.New(rand.NewSource(time.Now().UnixNano()))

// toPDF returns pdf file as an slice of bytes.
// Receives an image.Image as a parameter.
func toPDF(img image.Image) ([]byte, error) {
	// Sets a Rectangle based on the size of the image.
	imgRect := gopdf.Rect{
		W: float64(img.Bounds().Dx()),
		H: float64(img.Bounds().Dy()),
	}

	// Init the pdf obkect.
	pdf := gopdf.GoPdf{}

	// Sets the size of the every pdf page,
	// based on the dimensions of the image.
	pdf.Start(
		gopdf.Config{
			PageSize: imgRect,
		},
	)

	// Add a page to the PDF.
	pdf.AddPage()

	// Draws the image on the rectangle on the page above created.
	if err := pdf.ImageFrom(img, 0, 0, &imgRect); err != nil {
		return nil, err
	}

	// Creates a bytes.Buffer and writes the pdf data to it.
	buf := new(bytes.Buffer)
	if _, err := pdf.WriteTo(buf); err != nil {
		return nil, err
	}

	// Returns the pdf data as slice of bytes.
	return buf.Bytes(), nil
}

func ParseMimeType(mimetype string) string {
	if !strings.Contains(mimetype, imageMimeType) {
		return mimetype
	}

	return strings.TrimPrefix(mimetype, imageMimeType)
}

// stringWithCharset returns a random string based a length and a charset.
func stringWithCharset(length int, charset string) string {
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b)
}

// randString returns a random string calling stringWithCharset and using the letters constant.
func randString(length int) string {
	return stringWithCharset(length, letters)
}

// convertToImage retuns an image as io.Reader and error if something goes wrong.
// It gets the target format as input alongside the image to be converted to that format.
func convertToImage(target string, file io.Reader) (io.Reader, error) {
	// Create a buffer meant to store the input file data.
	buf := new(bytes.Buffer)
	if _, err := buf.ReadFrom(file); err != nil {
		return nil, fmt.Errorf(
			"error reading from the image file: %w",
			err,
		)
	}

	// Get the bytes off the input image.
	inputReaderBytes := buf.Bytes()

	// Create a temporary empty file where the input image is gonna be stored.
	tmpInputImage, err := os.CreateTemp("/tmp", fmt.Sprintf("*.%s", target))
	if err != nil {
		return nil, fmt.Errorf("error creating temporary image file: %w", err)
	}
	defer os.Remove(tmpInputImage.Name())

	// Write the content of the input image into the temporary file.
	if _, err = tmpInputImage.Write(inputReaderBytes); err != nil {
		return nil, fmt.Errorf("error writting the input reader to the temporary image file")
	}

	tmpConvertedFilename := fmt.Sprintf("/tmp/%s.%s", randString(10), target)

	// Convert the input image to the target format.
	// This is calling the ffmpeg command under the hood.
	// The reason behind this is that we could avoid using different libraries,
	// when we can use a use a single tool for multiple things.
	if err = ffmpeg.Input(tmpInputImage.Name()).
		Output(tmpConvertedFilename).
		OverWriteOutput().ErrorToStdOut().Run(); err != nil {
		return nil, err
	}

	// Open the converted file to get the bytes out of it,
	// and then turning them into a io.Reader.
	cf, err := os.Open(tmpConvertedFilename)
	if err != nil {
		return nil, err
	}
	defer os.Remove(cf.Name())

	fileBytes, err := io.ReadAll(cf)
	if err != nil {
		return nil, err
	}

	return bytes.NewReader(fileBytes), nil
}

func convertToDocument(target string, img image.Image) ([]byte, error) {
	var err error
	var result []byte

	switch target {
	case PDF:
		result, err = toPDF(img)
		if err != nil {
			return nil, err
		}
	}

	return result, nil
}
