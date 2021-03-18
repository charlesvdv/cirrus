package cirrus_test

import (
	"testing"

	cirrus "github.com/charlesvdv/cirrus/backend"
	"github.com/stretchr/testify/require"
)

func TestPathParse(t *testing.T) {
	pathsTable := []struct {
		path     string
		elements []string
	}{
		{"", []string{}},
		{"/", []string{}},
		{"foo", []string{"foo"}},
		{"foo/bar", []string{"foo", "bar"}},
		{"foo//bar", []string{"foo", "bar"}},
		{"/foo///bar//baz", []string{"foo", "bar", "baz"}},
		{"/foo /bar", []string{"foo ", "bar"}},
		{"/foo \\/ bar/ baz", []string{"foo / bar", " baz"}},
		{"foo\\d bar", nil},
	}

	for _, tt := range pathsTable {
		t.Run(tt.path, func(t *testing.T) {
			path, err := cirrus.ParsePath(tt.path)
			if tt.elements == nil {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.elements, path.Elements())
			}
		})
	}
}
