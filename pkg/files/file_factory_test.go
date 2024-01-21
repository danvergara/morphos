package files

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/danvergara/morphos/pkg/files/documents"
	"github.com/danvergara/morphos/pkg/files/images"
)

func TestImageFactory(t *testing.T) {
	imgF, err := BuildFactory(Img, "foo.png")
	require.NoError(t, err)

	imageFile, err := imgF.NewFile(images.PNG)
	require.NoError(t, err)

	png, ok := imageFile.(Image)
	if !ok {
		t.Fatal("struct assertion has failed")
	}

	t.Logf("Png image has type %s", png.ImageType())
}

func TestDocumentFactory(t *testing.T) {
	docF, err := BuildFactory(Doc, "foo.pdf")
	require.NoError(t, err)

	docFile, err := docF.NewFile(documents.PDF)
	require.NoError(t, err)

	pdf, ok := docFile.(Document)
	if !ok {
		t.Fatal("struct assertion has failed")
	}

	t.Logf("PDF document has type %s", pdf.DocumentType())
}
