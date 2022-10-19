package lazyfstools_test

import (
	"encoding/json"
	"path/filepath"
	"testing"

	"github.com/ingcr3at1on/x/lazyfstools"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWriteFile(t *testing.T) {
	afs := afero.NewMemMapFs()
	abs := "test"
	err := afs.Mkdir(abs, 0777)
	require.NoError(t, err)

	err = lazyfstools.WriteFile(afs, filepath.Join(abs, "test.json"), func(f afero.File) error {
		_, err := f.WriteString(`{"str": "foobar"}`)
		return err
	})
	require.NoError(t, err)

	f, err := afs.Open(filepath.Join(abs, "test.json"))
	require.NoError(t, err)

	var st struct {
		Str string `json:"str"`
	}

	err = json.NewDecoder(f).Decode(&st)
	require.NoError(t, err)

	assert.Equal(t, "foobar", st.Str)
}
