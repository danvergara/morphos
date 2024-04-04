package main

import (
	"embed"
	"errors"
	"fmt"
	"github.com/davidbyttow/govips/v2/vips"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/gabriel-vasile/mimetype"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"

	"github.com/danvergara/morphos/pkg/files"
)

const (
	uploadFileFormField = "uploadFile"
)

var (
	//go:embed all:templates
	templatesHTML embed.FS

	//go:embed all:static
	staticFiles embed.FS
	// Upload path.
	// It is a variable now, which means that can be
	// cofigurable through a environment variable.
	uploadPath string
)

func init() {
	uploadPath = os.Getenv("TMP_DIR")
	if uploadPath == "" {
		uploadPath = "/tmp"
	}
}

// statusError struct is the error representation
// at the HTTP layer.
type statusError struct {
	error
	status int
}

// Unwrap method returns the inner error.
func (e statusError) Unwrap() error { return e.error }

// HTTPStatus returns a HTTP status code.
func HTTPStatus(err error) int {
	if err == nil {
		return 0
	}

	var statusErr interface {
		error
		HTTPStatus() int
	}

	// Checks if err implements the statusErr interface.
	if errors.As(err, &statusErr) {
		return statusErr.HTTPStatus()
	}

	// Returns a default status code if none is provided.
	return http.StatusInternalServerError
}

// WithHTTPStatus returns an error with the original error and the status code.
func WithHTTPStatus(err error, status int) error {
	return statusError{
		error:  err,
		status: status,
	}
}

// toHandler is a wrapper for functions that have the following signature:
// func(http.ResponseWriter, *http.Request) error
// So, regular handlers can return an error that can be unwrapped.
// If an errors is received at the time to execute the original handler,
// renderError function is called.
func toHandler(f func(http.ResponseWriter, *http.Request) error) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := f(w, r)
		if err != nil {
			renderError(w, err.Error(), HTTPStatus(err))
		}
	}
}

type ConvertedFile struct {
	Filename string
	FileType string
}

func index(w http.ResponseWriter, _ *http.Request) error {
	tmpls := []string{
		"templates/base.tmpl",
		"templates/partials/htmx.tmpl",
		"templates/partials/style.tmpl",
		"templates/partials/nav.tmpl",
		"templates/partials/form.tmpl",
		"templates/partials/modal.tmpl",
		"templates/partials/js.tmpl",
	}

	tmpl, err := template.ParseFS(templatesHTML, tmpls...)
	if err != nil {
		log.Printf("error ocurred parsing templates: %v", err)
		return WithHTTPStatus(err, http.StatusInternalServerError)
	}

	err = tmpl.ExecuteTemplate(w, "base", nil)
	if err != nil {
		log.Printf("error ocurred executing template: %v", err)
		return WithHTTPStatus(err, http.StatusInternalServerError)
	}

	return nil
}

func handleUploadFile(w http.ResponseWriter, r *http.Request) error {
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
		return WithHTTPStatus(err, http.StatusBadRequest)
	}
	defer file.Close()

	fileBytes, err := io.ReadAll(file)
	if err != nil {
		log.Printf("error ocurred reading file: %v", err)
		return WithHTTPStatus(err, http.StatusBadRequest)
	}

	targetFileSubType := r.FormValue("input_format")

	detectedFileType := mimetype.Detect(fileBytes)

	fileType, subType, err := files.TypeAndSupType(detectedFileType.String())
	if err != nil {
		log.Printf("error occurred getting type and subtype from mimetype: %v", err)
		return WithHTTPStatus(err, http.StatusBadRequest)
	}

	fileFactory, err := files.BuildFactory(fileType, fileHeader.Filename)
	if err != nil {
		log.Printf("error occurred while getting a file factory: %v", err)
		return WithHTTPStatus(err, http.StatusBadRequest)
	}

	f, err := fileFactory.NewFile(subType)
	if err != nil {
		log.Printf("error occurred getting the file object: %v", err)
		return WithHTTPStatus(err, http.StatusBadRequest)
	}

	targetFileType := files.SupportedFileTypes()[targetFileSubType]
	convertedFile, err = f.ConvertTo(
		cases.Title(language.English).String(targetFileType),
		targetFileSubType,
		fileBytes,
	)
	if err != nil {
		log.Printf("error ocurred while converting image %v", err)
		return WithHTTPStatus(err, http.StatusInternalServerError)
	}

	if fileType == "application" {
		targetFileSubType = "zip"
	}

	convertedFileName = filename(fileHeader.Filename, targetFileSubType)
	convertedFilePath = filepath.Join(uploadPath, convertedFileName)

	newFile, err := os.Create(convertedFilePath)
	if err != nil {
		log.Printf("error occurred converting file: %v", err)
		return WithHTTPStatus(err, http.StatusInternalServerError)
	}
	defer newFile.Close()
	if _, err := newFile.Write(convertedFile); err != nil {
		log.Printf("error occurred writing file: %v", err)
		return WithHTTPStatus(err, http.StatusInternalServerError)
	}

	tmpls := []string{
		"templates/partials/card_file.tmpl",
		"templates/partials/modal.tmpl",
	}

	tmpl, err := template.ParseFS(templatesHTML, tmpls...)
	if err != nil {
		log.Printf("error occurred parsing template files: %v", err)
		return WithHTTPStatus(err, http.StatusInternalServerError)
	}

	convertedFileMimeType := mimetype.Detect(convertedFile)

	convertedFileType, _, err := files.TypeAndSupType(convertedFileMimeType.String())
	if err != nil {
		log.Printf("error occurred getting the file type of the result file: %v", err)
		return WithHTTPStatus(err, http.StatusInternalServerError)
	}

	err = tmpl.ExecuteTemplate(
		w,
		"content",
		ConvertedFile{Filename: convertedFileName, FileType: convertedFileType},
	)
	if err != nil {
		log.Printf("error occurred executing template files: %v", err)
		return WithHTTPStatus(err, http.StatusInternalServerError)
	}

	return nil
}

func handleFileFormat(w http.ResponseWriter, r *http.Request) error {
	file, _, err := r.FormFile(uploadFileFormField)
	if err != nil {
		log.Printf("error ocurred while getting file from form: %v", err)
		return WithHTTPStatus(err, http.StatusBadRequest)
	}

	fileBytes, err := io.ReadAll(file)
	if err != nil {
		log.Printf("error occurred while executing template files: %v", err)
		return WithHTTPStatus(err, http.StatusBadRequest)
	}

	detectedFileType := mimetype.Detect(fileBytes)

	templates := []string{
		"templates/partials/form.tmpl",
	}

	fileType, subType, err := files.TypeAndSupType(detectedFileType.String())
	if err != nil {
		log.Printf("error occurred getting type and subtype from mimetype: %v", err)
		return WithHTTPStatus(err, http.StatusBadRequest)
	}

	fileFactory, err := files.BuildFactory(fileType, "")
	if err != nil {
		log.Printf("error occurred while getting a file factory: %v", err)
		return WithHTTPStatus(err, http.StatusBadRequest)
	}

	f, err := fileFactory.NewFile(subType)
	if err != nil {
		log.Printf("error occurred getting the file object: %v", err)
		return WithHTTPStatus(err, http.StatusInternalServerError)
	}

	tmpl, err := template.ParseFS(templatesHTML, templates...)
	if err = tmpl.ExecuteTemplate(w, "format-elements", f.SupportedFormats()); err != nil {
		log.Printf("error occurred parsing template files: %v", err)
		return WithHTTPStatus(err, http.StatusInternalServerError)
	}

	return nil
}

func handleModal(w http.ResponseWriter, r *http.Request) error {
	filename := r.URL.Query().Get("filename")
	filetype := r.URL.Query().Get("filetype")

	tmpls := []string{
		"templates/partials/active_modal.tmpl",
	}

	tmpl, err := template.ParseFS(templatesHTML, tmpls...)
	if err != nil {
		log.Printf("error occurred parsing template files: %v", err)
		return WithHTTPStatus(err, http.StatusInternalServerError)
	}

	if err = tmpl.ExecuteTemplate(w, "content", ConvertedFile{Filename: filename, FileType: filetype}); err != nil {
		log.Printf("error occurred executing template files: %v", err)
		return WithHTTPStatus(err, http.StatusInternalServerError)
	}

	return nil
}

func main() {
	// start vips
	vips.Startup(&vips.Config{ConcurrencyLevel: runtime.NumCPU()})
	defer vips.Shutdown()

	port := os.Getenv("MORPHOS_PORT")
	// default port.
	if port == "" {
		port = "8080"
	}

	r := chi.NewRouter()
	r.Use(middleware.Logger)

	fsUpload := http.FileServer(http.Dir(uploadPath))

	var staticFS = http.FS(staticFiles)
	fs := http.FileServer(staticFS)

	r.Handle("/static/*", fs)
	r.Handle("/files/*", http.StripPrefix("/files", fsUpload))
	r.Get("/", toHandler(index))
	r.Post("/upload", toHandler(handleUploadFile))
	r.Post("/format", toHandler(handleFileFormat))
	r.Get("/modal", toHandler(handleModal))

	http.ListenAndServe(fmt.Sprintf(":%s", port), r)
}

// renderError functions executes the error template.
func renderError(w http.ResponseWriter, message string, statusCode int) {
	w.WriteHeader(statusCode)
	tmpl, _ := template.ParseFS(templatesHTML, "templates/partials/error.tmpl")
	tmpl.ExecuteTemplate(w, "error", struct{ ErrorMessage string }{ErrorMessage: message})
}

func fileNameWithoutExtension(fileName string) string {
	return strings.TrimSuffix(filepath.Base(fileName), filepath.Ext(fileName))
}

func filename(filename, extension string) string {
	return fmt.Sprintf("%s.%s", fileNameWithoutExtension(filename), extension)
}
