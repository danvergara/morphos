package documents

import (
	"archive/zip"
	"bytes"
	"errors"
	"fmt"
	"image/gif"
	"image/jpeg"
	"image/png"
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"slices"
	"strings"

	"github.com/chai2010/webp"
	"github.com/gen2brain/go-fitz"
	"golang.org/x/image/bmp"
	"golang.org/x/image/tiff"

	"github.com/danvergara/morphos/pkg/files/images"
)

// Pdf struct implements the File and Document interface from the file package.
type Pdf struct {
	filename            string
	compatibleFormats   map[string][]string
	compatibleMIMETypes map[string][]string
	OutDir              string
}

// NewPdf returns a pointer to Pdf.
func NewPdf(filename string) *Pdf {
	p := Pdf{
		filename: filename,
		compatibleFormats: map[string][]string{
			"Image": {
				images.JPG,
				images.JPEG,
				images.PNG,
				images.GIF,
				images.WEBP,
				images.TIFF,
				images.BMP,
			},
			"Document": {
				DOCX,
			},
		},
		compatibleMIMETypes: map[string][]string{
			"Image": {
				images.JPG,
				images.JPEG,
				images.PNG,
				images.GIF,
				images.WEBP,
				images.TIFF,
				images.BMP,
			},
			"Document": {
				DOCXMIMEType,
			},
		},
	}

	return &p
}

// SupportedFormats returns a map witht the compatible formats that Pdf is
// compatible to be converted to.
func (p *Pdf) SupportedFormats() map[string][]string {
	return p.compatibleFormats
}

// SupportedMIMETypes returns a map witht the compatible MIME types that Pdf is
// compatible to be converted to.
func (p *Pdf) SupportedMIMETypes() map[string][]string {
	return p.compatibleMIMETypes
}

// ConvertTo converts the current PDF file to another given format.
// This method receives the file type, the sub-type and the file as an slice of bytes.
// Returns the converted file as an slice of bytes, if something wrong happens, an error is returned.
func (p *Pdf) ConvertTo(fileType, subType string, fileBytes []byte) ([]byte, error) {
	// These are guard clauses that check if the target file type is valid.
	compatibleFormats, ok := p.SupportedFormats()[fileType]
	if !ok {
		return nil, fmt.Errorf("ConvertTo: file type not supported: %s", fileType)
	}

	if !slices.Contains(compatibleFormats, subType) {
		return nil, fmt.Errorf("ConvertTo: file sub-type not supported: %s", subType)
	}

	// If the file type is valid, figures out how to go ahead.
	switch strings.ToLower(fileType) {
	case imageType:
		// Creates a PDF Reader based on the pdf file.
		doc, err := fitz.NewFromMemory(fileBytes)
		if err != nil {
			return nil, fmt.Errorf("ConvertTo: error at opening the input pdf: %w", err)
		}

		// Parses the file name of the Zip file.
		zipFileName := filepath.Join("/tmp", fmt.Sprintf(
			"%s.zip",
			strings.TrimSuffix(p.filename, filepath.Ext(p.filename)),
		))

		// Creates the zip file that will be returned.
		archive, err := os.Create(zipFileName)
		if err != nil {
			return nil, fmt.Errorf(
				"ConvertTo: error at creating the zip file to store the images: %w",
				err,
			)
		}

		// Creates a Zip Writer to add files later on.
		zipWriter := zip.NewWriter(archive)

		for n := 0; n < doc.NumPage(); n++ {
			// Parses the file name image.
			imgFileName := fmt.Sprintf(
				"%s_%d.%s",
				strings.TrimSuffix(p.filename, filepath.Ext(p.filename)),
				n,
				subType,
			)

			tmpImgFileMame := filepath.Join("/tmp", imgFileName)

			// Converts the current pdf page to an image.Image.
			img, err := doc.Image(n)
			if err != nil {
				return nil, fmt.Errorf(
					"ConvertTo: error at converting the pdf page number %d to image: %w",
					n,
					err,
				)
			}

			// Saves the image on disk.
			imgFile, err := os.Create(tmpImgFileMame)
			if err != nil {
				return nil, fmt.Errorf(
					"ConvertTo: error at storing the pdf image from the page #%d: %w",
					n,
					err,
				)
			}

			// Encodes the image based on the sub-type of the file.
			// e.g. png.
			switch subType {
			case images.PNG:
				err = png.Encode(imgFile, img)
				if err != nil {
					return nil, fmt.Errorf(
						"ConvertTo: error at encoding the pdf page %d as png: %w",
						n,
						err,
					)
				}
			case images.JPG, images.JPEG:
				err = jpeg.Encode(imgFile, img, nil)
				if err != nil {
					return nil, fmt.Errorf(
						"ConvertTo: error at encoding the pdf page %d as jpeg: %w",
						n,
						err,
					)
				}
			case images.GIF:
				err = gif.Encode(imgFile, img, nil)
				if err != nil {
					return nil, fmt.Errorf(
						"ConvertTo: error at encoding the pdf page %d as gif: %w",
						n,
						err,
					)
				}
			case images.WEBP:
				err = webp.Encode(imgFile, img, nil)
				if err != nil {
					return nil, fmt.Errorf(
						"ConvertTo: error at encoding the pdf page %d as webp: %w",
						n,
						err,
					)
				}
			case images.TIFF:
				err = tiff.Encode(imgFile, img, nil)
				if err != nil {
					return nil, fmt.Errorf(
						"ConvertTo: error at encoding the pdf page %d as tiff: %w",
						n,
						err,
					)
				}
			case images.BMP:
				err = bmp.Encode(imgFile, img)
				if err != nil {
					return nil, fmt.Errorf(
						"ConvertTo: error at encoding the pdf page %d as bmp: %w",
						n,
						err,
					)
				}
			}

			imgFile.Close()

			// Opens the image to add it to the zip file.
			imgFile, err = os.Open(tmpImgFileMame)
			if err != nil {
				return nil, fmt.Errorf(
					"ConvertTo: error at storing the pdf image from the page #%d: %w",
					n,
					err,
				)
			}
			defer imgFile.Close()

			// Adds the image to the zip file.
			w1, err := zipWriter.Create(imgFileName)
			if err != nil {
				return nil, fmt.Errorf(
					"ConvertTo: error at creating a zip writer to store the page #%d: %w",
					n,
					err,
				)
			}

			if _, err := io.Copy(w1, imgFile); err != nil {
				return nil, fmt.Errorf(
					"ConvertTo: error at copying the content of the page #%d to the zipwriter: %w",
					n,
					err,
				)
			}
		}

		// Closes both zip writer and the zip file after its done with the writing.
		zipWriter.Close()
		archive.Close()

		// Reads the zip file as an slice of bytes.
		zipFile, err := os.ReadFile(zipFileName)
		if err != nil {
			return nil, fmt.Errorf("error reading zip file: %v", err)
		}

		return zipFile, nil
	case documentType:
		switch subType {
		case DOCX:
			var (
				stdout bytes.Buffer
				stderr bytes.Buffer
			)

			pdfFilename := filepath.Join("/tmp", p.filename)
			docxFileName := fmt.Sprintf(
				"%s.docx",
				strings.TrimSuffix(p.filename, filepath.Ext(p.filename)),
			)
			tmpDocxFileName := filepath.Join("/tmp", fmt.Sprintf(
				"%s.docx",
				strings.TrimSuffix(p.filename, filepath.Ext(p.filename)),
			))

			// Parses the file name of the Zip file.
			zipFileName := filepath.Join("/tmp", fmt.Sprintf(
				"%s.zip",
				strings.TrimSuffix(p.filename, filepath.Ext(p.filename)),
			))

			pdfFile, err := os.Create(pdfFilename)
			if err != nil {
				return nil, fmt.Errorf(
					"error creating file to store the incoming pdf locally %s: %w",
					p.filename,
					err,
				)
			}
			defer pdfFile.Close()

			if _, err := pdfFile.Write(fileBytes); err != nil {
				return nil, fmt.Errorf(
					"error storing the incoming pdf file %s: %w",
					p.filename,
					err,
				)
			}

			tmpDocxFile, err := os.Create(tmpDocxFileName)
			if err != nil {
				return nil, fmt.Errorf(
					"error at creating the temporary docx file to store the docx content: %w",
					err,
				)
			}

			cmdStr := "libreoffice --headless --infilter='writer_pdf_import' --convert-to %s --outdir %s %q"
			cmd := exec.Command(
				"bash",
				"-c",
				fmt.Sprintf(cmdStr, `docx:"MS Word 2007 XML"`, "/tmp", pdfFilename),
			)

			cmd.Stdout = &stdout
			cmd.Stderr = &stderr

			if err := cmd.Run(); err != nil {
				return nil, fmt.Errorf(
					"error converting pdf to docx using libreoffice: %w",
					err,
				)
			}

			if stderr.String() != "" {
				return nil, fmt.Errorf(
					"error converting pdf to docx calling libreoffice: %s",
					stderr.String(),
				)
			}

			log.Println(stdout.String())

			tmpDocxFile.Close()

			tmpDocxFile, err = os.Open(tmpDocxFileName)
			if err != nil {
				return nil, fmt.Errorf(
					"error at opening the docx file: %w",
					err,
				)
			}
			defer tmpDocxFile.Close()

			// Creates the zip file that will be returned.
			archive, err := os.Create(zipFileName)
			if err != nil {
				return nil, fmt.Errorf(
					"error at creating the zip file to store the docx file: %w",
					err,
				)
			}

			// Creates a Zip Writer to add files later on.
			zipWriter := zip.NewWriter(archive)

			w1, err := zipWriter.Create(docxFileName)
			if err != nil {
				return nil, fmt.Errorf(
					"eror at creating a zip file: %w",
					err,
				)
			}

			if _, err := io.Copy(w1, tmpDocxFile); err != nil {
				return nil, fmt.Errorf(
					"error at writing the docx file content to the zip writer: %w",
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

			return zipFile, nil
		}
	}

	return nil, errors.New("not implemented")
}

// DocumentType returns the type of ducument of Pdf.
func (p *Pdf) DocumentType() string {
	return PDF
}
