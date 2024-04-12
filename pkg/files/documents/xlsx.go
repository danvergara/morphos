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

// Xlsx struct implements the File and Document interface from the file package.
type Xlsx struct {
	filename            string
	compatibleFormats   map[string][]string
	compatibleMIMETypes map[string][]string
}

// NewXlsx returns a pointer to Xlsx.
func NewXlsx(filename string) *Xlsx {
	x := Xlsx{
		filename: filename,
		compatibleFormats: map[string][]string{
			"Document": {
				CSV,
			},
		},
		compatibleMIMETypes: map[string][]string{
			"Document": {
				CSV,
			},
		},
	}

	return &x
}

// SupportedFormats returns a map witht the compatible formats that Xlsx is
// compatible to be converted to.
func (x *Xlsx) SupportedFormats() map[string][]string {
	return x.compatibleFormats
}

// SupportedMIMETypes returns a map witht the compatible MIME types that Docx is
// compatible to be converted to.
func (x *Xlsx) SupportedMIMETypes() map[string][]string {
	return x.compatibleMIMETypes
}

func (x *Xlsx) ConvertTo(fileType, subType string, file io.Reader) (io.Reader, error) {
	compatibleFormats, ok := x.SupportedFormats()[fileType]
	if !ok {
		return nil, fmt.Errorf("file type not supported: %s", fileType)
	}

	if !slices.Contains(compatibleFormats, subType) {
		return nil, fmt.Errorf("sub-type not supported: %s", subType)
	}

	buf := new(bytes.Buffer)
	if _, err := buf.ReadFrom(file); err != nil {
		return nil, fmt.Errorf(
			"error getting the content of the xlsx file in form of slice of bytes: %w",
			err,
		)
	}

	fileBytes := buf.Bytes()

	switch strings.ToLower(fileType) {
	case documentType:
		switch subType {
		case CSV:
			xlFile, err := xlsx.OpenBinary(fileBytes)
			if err != nil {
				return nil, fmt.Errorf("error trying to open the xlsx file based on bytes of file %w", err)
			}

			// Parses the file name of the Zip file.
			zipFileName := filepath.Join("/tmp", fmt.Sprintf(
				"%s.zip",
				strings.TrimSuffix(x.filename, filepath.Ext(x.filename)),
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

			for i, sheet := range xlFile.Sheets {
				csvFilename := fmt.Sprintf(
					"%s_%d.%s",
					strings.TrimSuffix(x.filename, filepath.Ext(x.filename)),
					i+1,
					subType,
				)

				tmpCsvFilename := filepath.Join("/tmp", csvFilename)

				// Saves the image on disk.
				csvFile, err := os.Create(tmpCsvFilename)
				if err != nil {
					return nil, fmt.Errorf(
						"error at storing the tmp csv file from the xlsx sheet #%d: %w",
						i+1,
						err,
					)
				}

				cw := csv.NewWriter(csvFile)

				var vals []string
				err = sheet.ForEachRow(func(row *xlsx.Row) error {
					if row != nil {
						vals = vals[:0]
						err := row.ForEachCell(func(cell *xlsx.Cell) error {
							str, err := cell.FormattedValue()
							if err != nil {
								return err
							}
							vals = append(vals, str)
							return nil
						})
						if err != nil {
							return err
						}
					}
					cw.Write(vals)
					return nil
				})
				if err != nil {
					return nil, fmt.Errorf("error at creating a csv based of a xlsx sheet %d %w", i+1, err)
				}

				cw.Flush()
				if cw.Error() != nil {
					return nil, fmt.Errorf("error at writing buffered data to a underlying csv %w", err)
				}

				csvFile.Close()

				// Saves the image on disk.
				csvFile, err = os.Open(tmpCsvFilename)
				if err != nil {
					return nil, fmt.Errorf(
						"error at opening the csv file based off a xlsx sheet #%d: %w",
						i+1,
						err,
					)
				}

				defer csvFile.Close()

				// Adds the image to the zip file.
				w1, err := zipWriter.Create(csvFilename)
				if err != nil {
					return nil, fmt.Errorf(
						"error at creating a zip writer to store the xlsx sheet #%d: %w",
						i+1,
						err,
					)
				}

				if _, err := io.Copy(w1, csvFile); err != nil {
					return nil, fmt.Errorf(
						"error at copying the content of the xlsx sheet #%d to the zipwriter: %w",
						i+1,
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

			return bytes.NewReader(zipFile), nil
		}
	}

	return nil, errors.New("not implemented")
}

func (x *Xlsx) DocumentType() string {
	return XLSX
}
