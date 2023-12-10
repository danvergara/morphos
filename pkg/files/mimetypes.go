package files

import (
	"fmt"
	"strings"
)

// TypeAndSupType returns a the type and the sub-type of a
// given mimetype.
// e.g. image/png
// type: image
// subtype: png
func TypeAndSupType(mimetype string) (string, string, error) {
	types := strings.Split(mimetype, "/")

	if len(types) != 2 {
		return "", "", fmt.Errorf("%s not valid", mimetype)
	}

	return types[0], types[1], nil
}
