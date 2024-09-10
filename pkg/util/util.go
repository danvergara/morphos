package util

import (
	"archive/zip"
	"bufio"
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// EbookConvert calls the ebook-convert binary from the Calibre project.
// It receives a inpunt format which is the format of the file passed to be converted,
// and a target format which is the the format that the file is going to be converted to.
// The function also receives the input file as an slice of bytes, which is the file that is
// going to be converted.
func EbookConvert(filename, inputFormat, outputFormat string, inputFile []byte) (io.Reader, error) {
	tmpInputFile, err := os.Create(
		fmt.Sprintf(
			"/tmp/%s.%s",
			strings.TrimSuffix(filename, filepath.Ext(filename)),
			inputFormat,
		),
	)
	if err != nil {
		return nil, fmt.Errorf("error creating temporary file: %w", err)
	}

	defer os.Remove(tmpInputFile.Name())

	// Write the content of the input file into the temporary file.
	if _, err = tmpInputFile.Write(inputFile); err != nil {
		return nil, fmt.Errorf(
			"error writting the input reader to the temporary file",
		)
	}
	if err := tmpInputFile.Close(); err != nil {
		return nil, err
	}

	// Parse the name of the output file.
	tmpOutputFileName := fmt.Sprintf(
		"%s.%s",
		strings.TrimSuffix(tmpInputFile.Name(), filepath.Ext(tmpInputFile.Name())),
		outputFormat,
	)

	// run the ebook-convert command with the input file and the name of the output file.
	cmd := exec.Command("ebook-convert", tmpInputFile.Name(), tmpOutputFileName)

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
	cf, err := os.Open(tmpOutputFileName)
	if err != nil {
		return nil, err
	}
	defer os.Remove(cf.Name())

	// Parse the file name of the Zip file.
	zipFileName := fmt.Sprintf(
		"%s.zip",
		strings.TrimSuffix(filename, filepath.Ext(filename)),
	)

	// Parse the output file name.
	outputFilename := fmt.Sprintf(
		"%s.%s",
		strings.TrimSuffix(filename, filepath.Ext(filename)),
		outputFormat,
	)

	// Creates the zip file that will be returned.
	archive, err := os.CreateTemp("", zipFileName)
	if err != nil {
		return nil, fmt.Errorf(
			"error at creating the zip file to store the file: %w",
			err,
		)
	}

	defer os.Remove(archive.Name())

	// Creates a Zip Writer to add files later on.
	zipWriter := zip.NewWriter(archive)

	// Adds the image to the zip file.
	w1, err := zipWriter.Create(outputFilename)
	if err != nil {
		return nil, fmt.Errorf(
			"error creating the zip writer: %w",
			err,
		)
	}

	// Copy the content of the converted file to the zip file.
	if _, err := io.Copy(w1, cf); err != nil {
		return nil, fmt.Errorf(
			"error at writing the file content to the zip writer: %w",
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
