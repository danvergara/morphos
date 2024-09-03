package documents

import (
	"archive/zip"
	"bufio"
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
func NewPdf(filename string) Pdf {
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
			"Ebook": {
				EPUB,
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
				EpubMimeType,
			},
			"Ebook": {
				EpubMimeType,
			},
		},
	}

	return p
}

// SupportedFormats returns a map witht the compatible formats that Pdf is
// compatible to be converted to.
func (p Pdf) SupportedFormats() map[string][]string {
	return p.compatibleFormats
}

// SupportedMIMETypes returns a map witht the compatible MIME types that Pdf is
// compatible to be converted to.
func (p Pdf) SupportedMIMETypes() map[string][]string {
	return p.compatibleMIMETypes
}

// ConvertTo converts the current PDF file to another given format.
// This method receives the file type, the sub-type and the file as an slice of bytes.
// Returns the converted file as an slice of bytes, if something wrong happens, an error is returned.
func (p Pdf) ConvertTo(fileType, subType string, file io.Reader) (io.Reader, error) {
	// These are guard clauses that check if the target file type is valid.
	compatibleFormats, ok := p.SupportedFormats()[fileType]
	if !ok {
		return nil, fmt.Errorf("ConvertTo: file type not supported: %s", fileType)
	}

	if !slices.Contains(compatibleFormats, subType) {
		return nil, fmt.Errorf("ConvertTo: file sub-type not supported: %s", subType)
	}

	buf := new(bytes.Buffer)
	if _, err := buf.ReadFrom(file); err != nil {
		return nil, fmt.Errorf(
			"error getting the content of the pdf file in form of slice of bytes: %w",
			err,
		)
	}
	fileBytes := buf.Bytes()

	// If the file type is valid, figures out how to go ahead.
	switch strings.ToLower(fileType) {
	case imageType:
		// Creates a PDF Reader based on the pdf file.
		doc, err := fitz.NewFromMemory(fileBytes)
		if err != nil {
			return nil, fmt.Errorf("ConvertTo: error at opening the input pdf: %w", err)
		}

		// Parses the file name of the Zip file.
		zipFileName := fmt.Sprintf(
			"%s.zip",
			strings.TrimSuffix(p.filename, filepath.Ext(p.filename)),
		)

		// Creates the zip file that will be returned.
		archive, err := os.CreateTemp("", zipFileName)
		if err != nil {
			return nil, fmt.Errorf(
				"ConvertTo: error at creating the zip file to store the images: %w",
				err,
			)
		}
		defer os.Remove(archive.Name())

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
			imgFile, err := os.Create(fmt.Sprintf("/tmp/%s", imgFileName))
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
			imgFile, err = os.Open(imgFile.Name())
			if err != nil {
				return nil, fmt.Errorf(
					"ConvertTo: error at storing the pdf image from the page #%d: %w",
					n,
					err,
				)
			}

			// Adds the image to the zip file.
			w1, err := zipWriter.Create(filepath.Base(imgFile.Name()))
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

			imgFile.Close()
			os.Remove(imgFile.Name())
		}

		// Closes both zip writer and the zip file after its done with the writing.
		zipWriter.Close()
		archive.Close()

		// Reads the zip file as an slice of bytes.
		zipFile, err := os.ReadFile(archive.Name())
		if err != nil {
			return nil, fmt.Errorf("error reading zip file: %v", err)
		}

		return bytes.NewReader(zipFile), nil
	case documentType:
		switch subType {
		case DOCX:
			var (
				stdout bytes.Buffer
				stderr bytes.Buffer
			)

			docxFileName := fmt.Sprintf(
				"%s.docx",
				strings.TrimSuffix(p.filename, filepath.Ext(p.filename)),
			)

			// Parses the file name of the Zip file.
			zipFileName := fmt.Sprintf(
				"%s.zip",
				strings.TrimSuffix(p.filename, filepath.Ext(p.filename)),
			)

			pdfFile, err := os.CreateTemp("", p.filename)
			if err != nil {
				return nil, fmt.Errorf(
					"error creating file to store the incoming pdf locally %s: %w",
					p.filename,
					err,
				)
			}
			defer os.Remove(pdfFile.Name())

			if _, err := pdfFile.Write(fileBytes); err != nil {
				return nil, fmt.Errorf(
					"error storing the incoming pdf file %s: %w",
					p.filename,
					err,
				)
			}

			tmpDocxFile, err := os.CreateTemp("", docxFileName)
			if err != nil {
				return nil, fmt.Errorf(
					"error at creating the temporary docx file to store the docx content: %w",
					err,
				)
			}
			defer os.Remove(tmpDocxFile.Name())

			cmdStr := "libreoffice --headless --infilter='writer_pdf_import' --convert-to %s --outdir %s %q"
			cmd := exec.Command(
				"bash",
				"-c",
				fmt.Sprintf(cmdStr, `docx:"MS Word 2007 XML"`, "/tmp", pdfFile.Name()),
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

			tmpDocxFile, err = os.Open(tmpDocxFile.Name())
			if err != nil {
				return nil, fmt.Errorf(
					"error at opening the docx file: %w",
					err,
				)
			}
			defer tmpDocxFile.Close()

			// Creates the zip file that will be returned.
			archive, err := os.CreateTemp("", zipFileName)
			if err != nil {
				return nil, fmt.Errorf(
					"error at creating the zip file to store the docx file: %w",
					err,
				)
			}
			defer os.Remove(archive.Name())

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
			zipFile, err := os.ReadFile(archive.Name())
			if err != nil {
				return nil, fmt.Errorf("error reading zip file: %v", err)
			}

			return bytes.NewReader(zipFile), nil
		}
	case ebookType:
		switch subType {
		case EPUB:
			// Create a temporary empty file where the input pdf is gonna be stored.
			tmpInputPDF, err := os.Create(
				fmt.Sprintf(
					"/tmp/%s.%s",
					strings.TrimSuffix(p.filename, filepath.Ext(p.filename)),
					PDF,
				),
			)
			if err != nil {
				return nil, fmt.Errorf("error creating temporary pdf file: %w", err)
			}
			defer os.Remove(tmpInputPDF.Name())

			// Write the content of the input pdf into the temporary file.
			if _, err = tmpInputPDF.Write(fileBytes); err != nil {
				return nil, fmt.Errorf(
					"error writting the input reader to the temporary pdf file",
				)
			}

			if err := tmpInputPDF.Close(); err != nil {
				return nil, err
			}

			epubName := fmt.Sprintf(
				"%s.epub",
				strings.TrimSuffix(tmpInputPDF.Name(), filepath.Ext(tmpInputPDF.Name())),
			)

			cmd := exec.Command("ebook-convert", tmpInputPDF.Name(), epubName)

			// Capture stdout.
			stdout, err := cmd.StdoutPipe()
			if err != nil {
				return nil, err
			}

			// Capture stderr.
			stderr, err := cmd.StderrPipe()
			if err != nil {
				return nil, err
			}

			// Start the command.
			if err := cmd.Start(); err != nil {
				return nil, err
			}
			// Create readers to read stdout and stderr.
			stdoutScanner := bufio.NewScanner(stdout)
			stderrScanner := bufio.NewScanner(stderr)

			// Read stdout line by line.
			go func() {
				for stdoutScanner.Scan() {
					log.Println("STDOUT:", stdoutScanner.Text())
				}
			}()

			// Read stderr line by line.
			go func() {
				for stderrScanner.Scan() {
					log.Println("STDERR:", stderrScanner.Text())
				}
			}()

			// Wait for the command to finish.
			if err := cmd.Wait(); err != nil {
				return nil, err
			}

			// Open the converted file to get the bytes out of it,
			// and then turning them into a io.Reader.
			cf, err := os.Open(epubName)
			if err != nil {
				return nil, err
			}
			defer os.Remove(cf.Name())

			// Parse the file name of the Zip file.
			zipFileName := fmt.Sprintf(
				"%s.zip",
				strings.TrimSuffix(p.filename, filepath.Ext(p.filename)),
			)

			// Parse the output file name.
			outputEpubFilename := fmt.Sprintf(
				"%s.%s",
				strings.TrimSuffix(p.filename, filepath.Ext(p.filename)),
				EPUB,
			)

			// Creates the zip file that will be returned.
			archive, err := os.CreateTemp("", zipFileName)
			if err != nil {
				return nil, fmt.Errorf(
					"error at creating the zip file to store the pdf file: %w",
					err,
				)
			}
			defer os.Remove(archive.Name())

			// Creates a Zip Writer to add files later on.
			zipWriter := zip.NewWriter(archive)

			// Adds the image to the zip file.
			w1, err := zipWriter.Create(outputEpubFilename)
			if err != nil {
				return nil, fmt.Errorf(
					"error creating the zip writer: %w",
					err,
				)
			}

			if _, err := io.Copy(w1, cf); err != nil {
				return nil, fmt.Errorf(
					"error at writing the docx file content to the zip writer: %w",
					err,
				)
			}

			// Closes both zip writer and the zip file after its done with the writing.
			zipWriter.Close()
			archive.Close()

			// Reads the zip file as an slice of bytes.
			zipFile, err := os.ReadFile(archive.Name())
			if err != nil {
				return nil, fmt.Errorf("error reading zip file: %v", err)
			}

			return bytes.NewReader(zipFile), nil
		}
	}

	return nil, errors.New("not implemented")
}

// DocumentType returns the type of ducument of Pdf.
func (p Pdf) DocumentType() string {
	return PDF
}
