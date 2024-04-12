package documents

import (
	"archive/zip"
	"bytes"
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"slices"
	"strings"

	"github.com/tealeg/xlsx/v3"
)

// Csv struct implements the File and Document interface from the file package.
type Csv struct {
	filename            string
	compatibleFormats   map[string][]string
	compatibleMIMETypes map[string][]string
}

// NewCsv returns a pointer to Csv.
func NewCsv(filename string) *Csv {
	c := Csv{
		filename: filename,
		compatibleFormats: map[string][]string{
			"Document": {
				XLSX,
			},
		},
		compatibleMIMETypes: map[string][]string{
			"Document": {
				XLSX,
			},
		},
	}

	return &c
}

// SupportedFormats returns a map witht the compatible formats that CSv is
// compatible to be converted to.
func (c *Csv) SupportedFormats() map[string][]string {
	return c.compatibleFormats
}

// SupportedMIMETypes returns a map witht the compatible MIME types that Docx is
// compatible to be converted to.
func (c *Csv) SupportedMIMETypes() map[string][]string {
	return c.compatibleMIMETypes
}

func (c *Csv) ConvertTo(fileType, subType string, file io.Reader) (io.Reader, error) {
	compatibleFormats, ok := c.SupportedFormats()[fileType]
	if !ok {
		return nil, fmt.Errorf("file type not supported: %s", fileType)
	}

	if !slices.Contains(compatibleFormats, subType) {
		return nil, fmt.Errorf("sub-type not supported: %s", subType)
	}

	switch strings.ToLower(fileType) {
	case documentType:
		switch subType {
		case XLSX:
			xlsxFilename := fmt.Sprintf(
				"%s.xlsx",
				strings.TrimSuffix(c.filename, filepath.Ext(c.filename)),
			)

			xlsxPath := filepath.Join("/tmp", xlsxFilename)

			// Parses the file name of the Zip file.
			zipFileName := filepath.Join("/tmp", fmt.Sprintf(
				"%s.zip",
				strings.TrimSuffix(c.filename, filepath.Ext(c.filename)),
			))

			reader := csv.NewReader(file)
			xlsxFile := xlsx.NewFile()
			sheet, err := xlsxFile.AddSheet(strings.TrimSuffix(c.filename, filepath.Ext(c.filename)))
			if err != nil {
				return nil, fmt.Errorf("error creating a xlsx sheet %w", err)
			}

			for {
				fields, err := reader.Read()
				if err == io.EOF {
					break
				}

				row := sheet.AddRow()
				for _, field := range fields {
					cell := row.AddCell()
					cell.Value = field
				}
			}

			xlsxFile.Save(xlsxPath)

			tmpCsvFile, err := os.Open(xlsxPath)
			if err != nil {
				return nil, fmt.Errorf(
					"error at opening the pdf file: %w",
					err,
				)
			}
			defer tmpCsvFile.Close()

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

			w1, err := zipWriter.Create(xlsxFilename)
			if err != nil {
				return nil, fmt.Errorf(
					"eror at creating a zip file: %w",
					err,
				)
			}

			if _, err := io.Copy(w1, tmpCsvFile); err != nil {
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

func (c *Csv) DocumentType() string {
	return CSV
}
