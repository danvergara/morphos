package documents

import "errors"

	"github.com/unidoc/unipdf/v3/model"
	"github.com/unidoc/unipdf/v3/render"

	"github.com/danvergara/morphos/pkg/files/images"
)

// Pdf struct implements the File and Document interface from the file package.
type Pdf struct {
	filename          string
	compatibleFormats map[string][]string
}

func NewPdf(filename string) *Pdf {
	p := Pdf{
		filename: filename,
		compatibleFormats: map[string][]string{
			"Image": {
				images.GIF,
				images.JPG,
				images.JPEG,
				images.PNG,
				images.TIFF,
				images.WEBP,
			},
		},
	}

	return &p
}

func (p *Pdf) SupportedFormats() map[string][]string {
	return make(map[string][]string)
}

func (p *Pdf) ConvertTo(fileType, subType string, fileBytes []byte) ([]byte, error) {
	return nil, errors.New("not implemented")
}

func (p *Pdf) DocumentType() string {
	return PDF
}
