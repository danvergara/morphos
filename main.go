package main

import (
	"bufio"
	"bytes"
	"fmt"
	"image/jpeg"
	"image/png"
	"io"
	"net/http"
	"os"
)

const uploadPath = "./upload"

func handleUploadFile(_ http.ResponseWriter, r *http.Request) {
	r.ParseMultipartForm(100)
	mForm := r.MultipartForm

	for k := range mForm.File {
		// k is the key of file part
		file, fileHeader, err := r.FormFile(k)
		if err != nil {
			fmt.Println("inovke FormFile error:", err)
			return
		}
		defer file.Close()
		fmt.Printf("the uploaded file: name[%s], size[%d], header[%#v]\n", fileHeader.Filename, fileHeader.Size, fileHeader.Header)

		f, ok := file.(*os.File)
		if !ok {
			fmt.Printf("not a file \n")
		}
		// Get the file size
		stat, err := f.Stat()
		if err != nil {
			fmt.Println(err)
			return
		}

		// Read the file into a byte slice
		bs := make([]byte, stat.Size())
		_, err = bufio.NewReader(file).Read(bs)
		if err != nil && err != io.EOF {
			fmt.Println(err)
			return
		}

		contentType := http.DetectContentType(bs)

		switch contentType {
		case "image/png":
			_, err = ToJpeg(bs)
			if err != nil {
				fmt.Println(err)
				return
			}
		case "image/jpeg":
			fmt.Println("jpeg converted")
			_, err = JpegToPng(bs)
			if err != nil {
				fmt.Println(err)
				return
			}
		}

		fmt.Printf("file %T uploaded ok\n", file)
	}
}

func main() {
	http.HandleFunc("/upload", handleUploadFile)
	http.ListenAndServe(":8080", nil)
}

// ToJpeg converts a PNG image to JPEG format
func ToJpeg(imageBytes []byte) ([]byte, error) {
	// Decode the PNG image bytes
	img, err := png.Decode(bytes.NewReader(imageBytes))

	if err != nil {
		return nil, err
	}

	buf := new(bytes.Buffer)

	// encode the image as a JPEG file
	if err := jpeg.Encode(buf, img, nil); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil

}

// JpegToPng converts a JPEG image to PNG format
func JpegToPng(imageBytes []byte) ([]byte, error) {
	img, err := jpeg.Decode(bytes.NewReader(imageBytes))

	if err != nil {
		return nil, err
	}

	buf := new(bytes.Buffer)

	if err := png.Encode(buf, img); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}
