package main

import (
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"github.com/danvergara/morphos/pkg/image"
)

const (
	uploadPath          = "/tmp"
	uploadFileFormField = "uploadFile"
)

type ConvertedFile struct {
	Filename string
}

func index(w http.ResponseWriter, _ *http.Request) {
	files := []string{
		"./templates/base.tmpl",
		"./templates/partials/htmx.tmpl",
		"./templates/partials/style.tmpl",
		"./templates/partials/nav.tmpl",
		"./templates/partials/form.tmpl",
		"./templates/partials/modal.tmpl",
		"./templates/partials/js.tmpl",
	}

	tmpl, err := template.ParseFiles(files...)
	if err != nil {
		renderError(w, "INTERNAL_ERROR", http.StatusInternalServerError)
		return
	}

	err = tmpl.ExecuteTemplate(w, "base", nil)
	if err != nil {
		renderError(w, "INTERNAL_ERROR", http.StatusInternalServerError)
		return
	}
}

func handleUploadFile(w http.ResponseWriter, r *http.Request) {
	var (
		convertedFile     []byte
		convertedFilePath string
		convertedFileName string
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
		convertedFile, err = image.JpegToPng(fileBytes)
		if err != nil {
			renderError(w, "INVALID_FILE", http.StatusBadRequest)
			return
		}

		convertedFileName = filename(fileHeader.Filename, "png")
		convertedFilePath = filepath.Join(uploadPath, convertedFileName)
	case "image/png":
		convertedFile, err = image.PngToJpeg(fileBytes)
		if err != nil {
			renderError(w, "INVALID_FILE", http.StatusBadRequest)
			return
		}

		convertedFileName = filename(fileHeader.Filename, "jpg")
		convertedFilePath = filepath.Join(uploadPath, convertedFileName)
	default:
		renderError(w, "INVALID_FILE_TYPE", http.StatusBadRequest)
		return
	}

	newFile, err := os.Create(convertedFilePath)
	if err != nil {
		log.Printf("error occurred converting file: %v", err)
		renderError(w, "CANT_WRITE_FILE", http.StatusInternalServerError)
		return
	}
	defer newFile.Close()
	if _, err := newFile.Write(convertedFile); err != nil {
		renderError(w, "CANT_WRITE_FILE", http.StatusInternalServerError)
		return
	}

	files := []string{
		"./templates/partials/card_file.tmpl",
		"./templates/partials/modal.tmpl",
	}

	tmpl, err := template.ParseFiles(files...)
	if err != nil {
		log.Printf("error occurred parsing template files: %v", err)
		renderError(w, "INTERNAL_ERROR", http.StatusInternalServerError)
		return
	}

	err = tmpl.ExecuteTemplate(w, "content", ConvertedFile{Filename: convertedFileName})
	if err != nil {
		log.Printf("error occurred executing template files: %v", err)
		renderError(w, "INTERNAL_ERROR", http.StatusInternalServerError)
		return
	}
}

func handleModal(w http.ResponseWriter, r *http.Request) {
	filename := r.URL.Query().Get("filename")

	files := []string{
		"./templates/partials/active_modal.tmpl",
	}

	tmpl, err := template.ParseFiles(files...)
	if err != nil {
		log.Printf("error occurred parsing template files: %v", err)
		renderError(w, "INTERNAL_ERROR", http.StatusInternalServerError)
		return
	}

	err = tmpl.ExecuteTemplate(w, "content", ConvertedFile{Filename: filename})
	if err != nil {
		log.Printf("error occurred executing template files: %v", err)
		renderError(w, "INTERNAL_ERROR", http.StatusInternalServerError)
		return
	}
}

func main() {
	r := chi.NewRouter()
	r.Use(middleware.Logger)

	fs := http.FileServer(http.Dir(uploadPath))

	r.Handle("/files/*", http.StripPrefix("/files", fs))
	r.Get("/", index)
	r.Post("/upload", handleUploadFile)
	r.Get("/modal", handleModal)

	http.ListenAndServe("localhost:8080", r)
}

func renderError(w http.ResponseWriter, message string, statusCode int) {
	w.WriteHeader(statusCode)
	w.Write([]byte(message))
}

func fileNameWithoutExtension(fileName string) string {
	return strings.TrimSuffix(filepath.Base(fileName), filepath.Ext(fileName))
}

func filename(filename, extension string) string {
	return fmt.Sprintf("%s.%s", fileNameWithoutExtension(filename), extension)
}
