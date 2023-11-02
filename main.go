package main

import (
	"embed"
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

	"github.com/danvergara/morphos/pkg/files/images"
)

const (
	uploadPath          = "/tmp"
	uploadFileFormField = "uploadFile"
)

var (
	//go:embed all:templates
	templatesHTML embed.FS

	//go:embed all:static
	staticFiles embed.FS
)

type ConvertedFile struct {
	Filename string
}

type FileFormat struct {
	Name string
	ID   int
}

func index(w http.ResponseWriter, _ *http.Request) {
	files := []string{
		"templates/base.tmpl",
		"templates/partials/htmx.tmpl",
		"templates/partials/style.tmpl",
		"templates/partials/nav.tmpl",
		"templates/partials/form.tmpl",
		"templates/partials/modal.tmpl",
		"templates/partials/js.tmpl",
	}

	tmpl, err := template.ParseFS(templatesHTML, files...)
	if err != nil {
		log.Printf("error ocurred parsing templates: %v", err)
		renderError(w, "INTERNAL_ERROR", http.StatusInternalServerError)
		return
	}

	err = tmpl.ExecuteTemplate(w, "base", nil)
	if err != nil {
		log.Printf("error ocurred executing template: %v", err)
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
		log.Printf("error ocurred getting file from form: %v", err)
		renderError(w, "INVALID_FILE", http.StatusBadRequest)
		return
	}
	defer file.Close()

	fileBytes, err := io.ReadAll(file)
	if err != nil {
		log.Printf("error ocurred reading file: %v", err)
		renderError(w, "INVALID_FILE", http.StatusBadRequest)
		return
	}

	fileType := r.FormValue("input_format")

	detectedFileType := http.DetectContentType(fileBytes)
	convertedFile, err = images.ConverImage(detectedFileType, fileType, fileBytes)
	if err != nil {
		log.Printf("error ocurred while converting image %v", err)
		renderError(w, "INVALID_FILE", http.StatusBadRequest)
		return
	}

	convertedFileName = filename(fileHeader.Filename, fileType)
	convertedFilePath = filepath.Join(uploadPath, convertedFileName)

	newFile, err := os.Create(convertedFilePath)
	if err != nil {
		log.Printf("error occurred converting file: %v", err)
		renderError(w, "CANT_WRITE_FILE", http.StatusInternalServerError)
		return
	}
	defer newFile.Close()
	if _, err := newFile.Write(convertedFile); err != nil {
		log.Printf("error occurred writing file: %v", err)
		renderError(w, "CANT_WRITE_FILE", http.StatusInternalServerError)
		return
	}

	files := []string{
		"templates/partials/card_file.tmpl",
		"templates/partials/modal.tmpl",
	}

	tmpl, err := template.ParseFS(templatesHTML, files...)
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

func handleFileFormat(w http.ResponseWriter, r *http.Request) {

	file, _, err := r.FormFile(uploadFileFormField)
	if err != nil {
		log.Printf("error ocurred getting file from form: %v", err)
		renderError(w, "INVALID_FILE", http.StatusBadRequest)
		return
	}

	fileBytes, err := io.ReadAll(file)
	if err != nil {
		log.Printf("error occurred executing template files: %v", err)
		renderError(w, "INVALID_FILE", http.StatusBadRequest)
		return
	}

	detectedFileType := http.DetectContentType(fileBytes)

	files := []string{
		"templates/partials/form.tmpl",
	}

	tmpl, err := template.ParseFS(templatesHTML, files...)
	formats := images.FileFormatsToConvert(detectedFileType)

	err = tmpl.ExecuteTemplate(w, "format-elements", formats)
	if err != nil {
		log.Printf("error occurred parsing template files: %v", err)
		renderError(w, "FINTERNAL_ERROR", http.StatusInternalServerError)
		return
	}
}

func handleModal(w http.ResponseWriter, r *http.Request) {
	filename := r.URL.Query().Get("filename")

	files := []string{
		"templates/partials/active_modal.tmpl",
	}

	tmpl, err := template.ParseFS(templatesHTML, files...)
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

	fsUpload := http.FileServer(http.Dir(uploadPath))

	var staticFS = http.FS(staticFiles)
	fs := http.FileServer(staticFS)

	r.Handle("/static/*", fs)
	r.Handle("/files/*", http.StripPrefix("/files", fsUpload))
	r.Get("/", index)
	r.Post("/upload", handleUploadFile)
	r.Post("/format", handleFileFormat)
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
