package files

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTypeAndSupType(t *testing.T) {
	type expected struct {
		fileType string
		subType  string
		hasErr   bool
	}

	var tests = []struct {
		name     string
		mimetype string
		expected expected
	}{}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			fileType, subType, err := TypeAndSupType(tc.mimetype)
			if tc.expected.hasErr {
				require.Error(t, err)
			}

			require.NoError(t, err)
			require.Equal(t, tc.expected.fileType, fileType)
			require.Equal(t, tc.expected.subType, subType)
		})
	}
}
