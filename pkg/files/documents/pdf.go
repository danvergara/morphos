package documents

import (
	"archive/zip"
	"bytes"
	"errors"
	"fmt"
	"image"
	"image/gif"
	"image/jpeg"
	"image/png"
	"io"
	"os"
	"path/filepath"
	"slices"
	"strings"

	"github.com/chai2010/webp"
	"github.com/unidoc/unioffice/document"
	"github.com/unidoc/unipdf/v3/extractor"
	"github.com/unidoc/unipdf/v3/model"
	"github.com/unidoc/unipdf/v3/render"
	"golang.org/x/image/bmp"
	"golang.org/x/image/tiff"

	"github.com/danvergara/morphos/pkg/files/images"
)

// Pdf struct implements the File and Document interface from the file package.
type Pdf struct {
	filename          string
	compatibleFormats map[string][]string
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
	}

	return &p
}

// SupportedFormats returns a map witht he compatible formats that Pds is
// compatible to be converted to.
func (p *Pdf) SupportedFormats() map[string][]string {
	return p.compatibleFormats
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
		pdfReader, err := model.NewPdfReader(bytes.NewReader(fileBytes))
		if err != nil {
			return nil, fmt.Errorf("ConvertTo: error at opening the input pdf: %w", err)
		}

		// Get the number of pages from the pdf file.
		pages, err := pdfReader.GetNumPages()
		if err != nil {
			return nil, fmt.Errorf(
				"ConvertTo: error at getting the number of pages from the input pdf: %w",
				err,
			)
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

		// Creats a Zip Writer to add files later on.
		zipWriter := zip.NewWriter(archive)

		device := render.NewImageDevice()
		// Set the image width. The height will be calculated accordingly.
		device.OutputWidth = 2048

		for pageNum := 1; pageNum <= pages; pageNum++ {
			// Parses the file name image.
			imgFileName := fmt.Sprintf(
				"%s_%d.%s",
				strings.TrimSuffix(p.filename, filepath.Ext(p.filename)),
				pageNum,
				subType,
			)

			tmpImgFileMame := filepath.Join("/tmp", imgFileName)

			// Converts the current pdf page to an image.Image.
			img, err := convertPDFPageToImage(pdfReader, device, pageNum)
			if err != nil {
				return nil, fmt.Errorf(
					"ConvertTo: error at converting the pdf page number %d to image: %w",
					pageNum,
					err,
				)
			}

			// Saves the image on disk.
			imgFile, err := os.Create(tmpImgFileMame)
			if err != nil {
				return nil, fmt.Errorf(
					"ConvertTo: error at storing the pdf image from the page #%d: %w",
					pageNum,
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
						pageNum,
						err,
					)
				}
			case images.JPG, images.JPEG:
				err = jpeg.Encode(imgFile, img, nil)
				if err != nil {
					return nil, fmt.Errorf(
						"ConvertTo: error at encoding the pdf page %d as jpeg: %w",
						pageNum,
						err,
					)
				}
			case images.GIF:
				err = gif.Encode(imgFile, img, nil)
				if err != nil {
					return nil, fmt.Errorf(
						"ConvertTo: error at encoding the pdf page %d as gif: %w",
						pageNum,
						err,
					)
				}
			case images.WEBP:
				err = webp.Encode(imgFile, img, nil)
				if err != nil {
					return nil, fmt.Errorf(
						"ConvertTo: error at encoding the pdf page %d as webp: %w",
						pageNum,
						err,
					)
				}
			case images.TIFF:
				err = tiff.Encode(imgFile, img, nil)
				if err != nil {
					return nil, fmt.Errorf(
						"ConvertTo: error at encoding the pdf page %d as tiff: %w",
						pageNum,
						err,
					)
				}
			case images.BMP:
				err = bmp.Encode(imgFile, img)
				if err != nil {
					return nil, fmt.Errorf(
						"ConvertTo: error at encoding the pdf page %d as bmp: %w",
						pageNum,
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
					pageNum,
					err,
				)
			}
			defer imgFile.Close()

			// Adds the image to the zip file.
			w1, err := zipWriter.Create(imgFileName)
			if err != nil {
				return nil, fmt.Errorf(
					"ConvertTo: error at creating a zip writer to store the page #%d: %w",
					pageNum,
					err,
				)
			}

			if _, err := io.Copy(w1, imgFile); err != nil {
				return nil, fmt.Errorf(
					"ConvertTo: error at copying the content of the page #%d to the zipwriter: %w",
					pageNum,
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
			pdfReader, err := model.NewPdfReader(bytes.NewReader(fileBytes))
			if err != nil {
				return nil, fmt.Errorf("ConvertTo: error at opening the input pdf: %w", err)
			}

			// Get the number of pages from the pdf file.
			pages, err := pdfReader.GetNumPages()
			if err != nil {
				return nil, fmt.Errorf(
					"ConvertTo: error at getting the number of pages from the input pdf: %w",
					err,
				)
			}

			doc := document.New()
			for pageNum := 1; pageNum <= pages; pageNum++ {
				page, err := pdfReader.GetPage(pageNum)
				if err != nil {
					if err != nil {
						return nil, fmt.Errorf(
							"ConvertTo: error at getting the current page %d: %w",
							pageNum,
							err,
						)
					}
				}

				extractor, err := extractor.New(page)
				if err != nil {
					return nil, fmt.Errorf(
						"ConvertTo: error at extracting the content of the current page %d: %w",
						pageNum,
						err,
					)
				}

				text, err := extractor.ExtractText()
				if err != nil {
					return nil, fmt.Errorf(
						"ConvertTo: error at extracting the content of the current page %d: %w",
						pageNum,
						err,
					)
				}

				para := doc.AddParagraph()
				run := para.AddRun()
				run.AddText(text)
			}
		}
	}

	return nil, errors.New("not implemented")
}

// convertPDFPageToImage converts the pdf page to an image.
// The functions receives the pdf Reader, the Image Device and the page number.
// Returns a image.Image or an error if something goes wrong.
func convertPDFPageToImage(
	pdfReader *model.PdfReader,
	device *render.ImageDevice,
	pageNum int,
) (image.Image, error) {
	// Get the page based on the given page number.
	page, err := pdfReader.GetPage(pageNum)
	if err != nil {
		return nil, fmt.Errorf(
			"error at getting a page given a page number %d: %w",
			pageNum,
			err,
		)
	}

	// Render returns an image.Image given a page.
	img, err := device.Render(page)
	if err != nil {
		return nil, fmt.Errorf(
			"error at converting the pdf page number %d to image: %w",
			pageNum,
			err,
		)
	}

	return img, nil
}

// DocumentType returns the type of ducument of Pdf.
func (p *Pdf) DocumentType() string {
	return PDF
}
