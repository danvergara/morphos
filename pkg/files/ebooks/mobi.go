package ebooks

import (
	"archive/zip"
	"bufio"
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

type Mobi struct {
	filename            string
	compatibleFormats   map[string][]string
	compatibleMIMETypes map[string][]string
}

func NewMobi(filename string) Mobi {
	m := Mobi{
		filename: filename,
		compatibleFormats: map[string][]string{
			"Ebook": {
				EPUB,
			},
		},
		compatibleMIMETypes: map[string][]string{
			"Ebook": {
				EPUB,
			},
		},
	}

	return m
}

// SupportedFormats returns a map witht the compatible formats that MOBI is
// compatible to be converted to.
func (m Mobi) SupportedFormats() map[string][]string {
	return m.compatibleFormats
}

// SupportedMIMETypes returns a map witht the compatible MIME types that MOBI is
// compatible to be converted to.
func (m Mobi) SupportedMIMETypes() map[string][]string {
	return m.compatibleMIMETypes
}

func (m Mobi) ConvertTo(fileType, subtype string, file io.Reader) (io.Reader, error) {
	// These are guard clauses that check if the target file type is valid.
	compatibleFormats, ok := m.SupportedFormats()[fileType]
	if !ok {
		return nil, fmt.Errorf("ConvertTo: file type not supported: %s", fileType)
	}

	if !slices.Contains(compatibleFormats, subtype) {
		return nil, fmt.Errorf("ConvertTo: file sub-type not supported: %s", subtype)
	}

	buf := new(bytes.Buffer)
	if _, err := buf.ReadFrom(file); err != nil {
		return nil, fmt.Errorf(
			"error getting the content of the pdf file in form of slice of bytes: %w",
			err,
		)
	}

	fileBytes := buf.Bytes()

	switch strings.ToLower(fileType) {
	case ebookType:
		switch subtype {
		case EPUB:
			// Create a temporary empty file where the input is gonna be stored.
			tmpInputMobi, err := os.CreateTemp("", fmt.Sprintf("*.%s", MOBI))
			if err != nil {
				return nil, fmt.Errorf("error creating temporary pdf file: %w", err)
			}
			defer os.Remove(tmpInputMobi.Name())

			// Write the content of the input pdf into the temporary file.
			if _, err = tmpInputMobi.Write(fileBytes); err != nil {
				return nil, fmt.Errorf(
					"error writting the input reader to the temporary pdf file",
				)
			}

			if err := tmpInputMobi.Close(); err != nil {
				return nil, err
			}

			epubFileName := fmt.Sprintf(
				"%s.%s",
				strings.TrimSuffix(tmpInputMobi.Name(), filepath.Ext(tmpInputMobi.Name())),
				EPUB,
			)

			// Parses the file name of the Zip file.
			zipFileName := fmt.Sprintf(
				"%s.zip",
				strings.TrimSuffix(m.filename, filepath.Ext(m.filename)),
			)

			// Parse the output file name.
			outputEpubFilename := fmt.Sprintf(
				"%s.%s",
				strings.TrimSuffix(m.filename, filepath.Ext(m.filename)),
				EPUB,
			)

			cmd := exec.Command("ebook-convert", tmpInputMobi.Name(), epubFileName)

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
			cf, err := os.Open(epubFileName)
			if err != nil {
				return nil, err
			}
			defer os.Remove(cf.Name())

			// Creates the zip file that will be returned.
			archive, err := os.CreateTemp("", zipFileName)
			if err != nil {
				return nil, fmt.Errorf(
					"error at creating the zip file to store the epub file: %w",
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
					"error at writing the epub file content to the zip writer: %w",
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

	return nil, errors.New("file format not implemented")
}

// EbookType returns the Ebook type which is MOBI in this case.
func (m Mobi) EbookType() string {
	return EPUB
}
