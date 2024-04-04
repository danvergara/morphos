package documents

import (
	"archive/zip"
	"bytes"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"slices"
	"strings"
)

// Docx struct implements the File and Document interface from the file package.
type Docx struct {
	filename            string
	compatibleFormats   map[string][]string
	compatibleMIMETypes map[string][]string
	OutDir              string
}

// NewDocx returns a pointer to Docx.
func NewDocx(filename string) *Docx {
	d := Docx{
		filename: filename,
		compatibleFormats: map[string][]string{
			"Document": {
				PDF,
			},
		},
		compatibleMIMETypes: map[string][]string{
			"Document": {
				PDF,
			},
		},
	}

	return &d
}

// SupportedFormats returns a map witht the compatible formats that Docx is
// compatible to be converted to.
func (d *Docx) SupportedFormats() map[string][]string {
	return d.compatibleFormats
}

// SupportedMIMETypes returns a map witht the compatible MIME types that Docx is
// compatible to be converted to.
func (d *Docx) SupportedMIMETypes() map[string][]string {
	return d.compatibleMIMETypes
}

func (d *Docx) ConvertTo(fileType, subType string, file io.Reader) (io.Reader, error) {
	compatibleFormats, ok := d.SupportedFormats()[fileType]
	if !ok {
		return nil, fmt.Errorf("file type not supported: %s", fileType)
	}

	if !slices.Contains(compatibleFormats, subType) {
		return nil, fmt.Errorf("sub-type not supported: %s", subType)
	}

	buf := new(bytes.Buffer)
	if _, err := buf.ReadFrom(file); err != nil {
		return nil, fmt.Errorf(
			"error getting the content of the docx file in form of slice of bytes: %w",
			err,
		)
	}
	fileBytes := buf.Bytes()

	switch strings.ToLower(fileType) {
	case documentType:
		switch subType {
		case PDF:
			var (
				stdout bytes.Buffer
				stderr bytes.Buffer
			)

			docxFilename := filepath.Join("/tmp", d.filename)
			pdfFileName := fmt.Sprintf(
				"%s.pdf",
				strings.TrimSuffix(d.filename, filepath.Ext(d.filename)),
			)
			tmpPdfFileName := filepath.Join("/tmp", fmt.Sprintf(
				"%s.pdf",
				strings.TrimSuffix(d.filename, filepath.Ext(d.filename)),
			))

			// Parses the file name of the Zip file.
			zipFileName := filepath.Join("/tmp", fmt.Sprintf(
				"%s.zip",
				strings.TrimSuffix(d.filename, filepath.Ext(d.filename)),
			))

			docxFile, err := os.Create(docxFilename)
			if err != nil {
				return nil, fmt.Errorf(
					"error creating file to store the incoming docx locally %s: %w",
					d.filename,
					err,
				)
			}
			defer docxFile.Close()

			if _, err := docxFile.Write(fileBytes); err != nil {
				return nil, fmt.Errorf(
					"error storing the incoming pdf file %s: %w",
					d.filename,
					err,
				)
			}

			tmpPdfFile, err := os.Create(tmpPdfFileName)
			if err != nil {
				return nil, fmt.Errorf(
					"error at creating the pdf file to store the pdf content: %w",
					err,
				)
			}

			cmdStr := "libreoffice --headless --convert-to pdf:writer_pdf_Export --outdir %s %q"
			cmd := exec.Command(
				"bash",
				"-c",
				fmt.Sprintf(cmdStr, "/tmp", docxFilename),
			)

			cmd.Stdout = &stdout
			cmd.Stderr = &stderr

			if err := cmd.Run(); err != nil {
				return nil, fmt.Errorf(
					"error converting docx to pdf using libreoffice: %s",
					err,
				)
			}

			if stderr.String() != "" {
				return nil, fmt.Errorf(
					"error converting docx to pdf calling libreoffice: %s",
					stderr.String(),
				)
			}

			log.Println(stdout.String())

			tmpPdfFile.Close()

			tmpPdfFile, err = os.Open(tmpPdfFileName)
			if err != nil {
				return nil, fmt.Errorf(
					"error at opening the pdf file: %w",
					err,
				)
			}
			defer tmpPdfFile.Close()

			// Creates the zip file that will be returned.
			archive, err := os.Create(zipFileName)
			if err != nil {
				return nil, fmt.Errorf(
					"error at creating the zip file to store the pdf file: %w",
					err,
				)
			}

			// Creates a Zip Writer to add files later on.
			zipWriter := zip.NewWriter(archive)

			w1, err := zipWriter.Create(pdfFileName)
			if err != nil {
				return nil, fmt.Errorf(
					"eror at creating a zip file: %w",
					err,
				)
			}

			if _, err := io.Copy(w1, tmpPdfFile); err != nil {
				return nil, fmt.Errorf(
					"error at writing the pdf file content to the zip writer: %w",
					err,
				)
			}

			// Closes both zip writer and the zip file after its done with the writing.
			zipWriter.Close()
			archive.Close()

			// Reads the zip file as an slice of bytes.
			zipFile, err := os.ReadFile(zipFileName)
			if err != nil {
				return nil, fmt.Errorf("error reading zip file: %v", err)
			}

			return bytes.NewReader(zipFile), nil
		}
	}

	return nil, errors.New("not implemented")
}

func (d *Docx) DocumentType() string {
	return DOCX
}
