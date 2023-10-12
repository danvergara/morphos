package main

import (
	"bytes"
	"fmt"
	"image/jpeg"
	"image/png"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

const (
	uploadPath          = "/tmp"
	uploadFileFormField = "uploadFile"
)

func handleUploadFile(w http.ResponseWriter, r *http.Request) {
	var (
		convertedFile     []byte
		convertedFilePath string
		err               error
	)

	// parse and validate file and post parameters.
	file, fileHeader, err := r.FormFile(uploadFileFormField)
	if err != nil {
		renderError(w, "INVALID_FILE", http.StatusBadRequest)
		return
	}
	defer file.Close()

	fileBytes, err := io.ReadAll(file)
	if err != nil {
		renderError(w, "INVALID_FILE", http.StatusBadRequest)
		return
	}

	detectedFileType := http.DetectContentType(fileBytes)
	switch detectedFileType {
	case "image/jpeg", "image/jpg":
		convertedFile, err = JpegToPng(fileBytes)
		if err != nil {
			renderError(w, "INVALID_FILE", http.StatusBadRequest)
			return
		}

		convertedFilePath = filepath.Join(uploadPath, fmt.Sprintf("%s.%s", fileNameWithoutExtension(fileHeader.Filename), "png"))
	case "image/png":
		convertedFile, err = PngToJpeg(fileBytes)
		if err != nil {
			renderError(w, "INVALID_FILE", http.StatusBadRequest)
			return
		}

		convertedFilePath = filepath.Join(uploadPath, fmt.Sprintf("%s.%s", fileNameWithoutExtension(fileHeader.Filename), "jpg"))
	default:
		renderError(w, "INVALID_FILE_TYPE", http.StatusBadRequest)
		return
	}

	newFile, err := os.Create(convertedFilePath)
	if err != nil {
		renderError(w, "CANT_WRITE_FILE", http.StatusInternalServerError)
		return
	}
	defer newFile.Close()
	if _, err := newFile.Write(convertedFile); err != nil {
		renderError(w, "CANT_WRITE_FILE", http.StatusInternalServerError)
		return
	}

	w.Write([]byte("SUCCESS"))
}

func main() {
	r := chi.NewRouter()
	r.Use(middleware.Logger)

	fs := http.FileServer(http.Dir(uploadPath))
	r.Handle("/files/*", http.StripPrefix("/files", fs))

	r.Post("/upload", handleUploadFile)

	http.ListenAndServe("localhost:8080", r)
}

// PngToJpeg converts a PNG image to JPEG format.
func PngToJpeg(imageBytes []byte) ([]byte, error) {
	// Decode the PNG image bytes.
	img, err := png.Decode(bytes.NewReader(imageBytes))

	if err != nil {
		return nil, err
	}

	buf := new(bytes.Buffer)

	// encode the image as a JPEG file.
	if err := jpeg.Encode(buf, img, nil); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

// JpegToPng converts a JPEG image to PNG format.
func JpegToPng(imageBytes []byte) ([]byte, error) {
	img, err := jpeg.Decode(bytes.NewReader(imageBytes))

	if err != nil {
		return nil, err
	}

	buf := new(bytes.Buffer)

	if err := png.Encode(buf, img); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func renderError(w http.ResponseWriter, message string, statusCode int) {
	w.WriteHeader(statusCode)
	w.Write([]byte(message))
}

func fileNameWithoutExtension(fileName string) string {
	return strings.TrimSuffix(filepath.Base(fileName), filepath.Ext(fileName))
}
