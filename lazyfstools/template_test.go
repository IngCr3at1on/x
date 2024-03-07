package lazyfstools_test

import (
	"github.com/ingcr3at1on/x/lazyfstools"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io"
	"path/filepath"
	"testing"
)

func TestRenderTemplates(t *testing.T) {
	afs := afero.NewMemMapFs()
	abs := "test"

	err := afs.Mkdir(abs, 0755)
	require.NoError(t, err)

	err = lazyfstools.WriteFile(afs, filepath.Join(abs, "test1.txt.tpl"), func(f afero.File) error {
		_, err := f.WriteString(`Hello {{ .Value }}!`)
		return err
	})
	require.NoError(t, err)

	err = lazyfstools.WriteFile(afs, filepath.Join(abs, "test2.txt.tpl"), func(f afero.File) error {
		_, err := f.WriteString(`{{ .Value }}!!!`)
		return err
	})
	require.NoError(t, err)

	target := "target"
	lazyfstools.RenderTemplates(afs, abs, target, struct {
		Value string
	}{
		Value: "Templates",
	})

	validateTemplate(t, afs, filepath.Join(target, "test1.txt"), "Hello Templates!")
	validateTemplate(t, afs, filepath.Join(target, "test2.txt"), "Templates!!!")
}

func validateTemplate(t *testing.T, afs afero.Fs, path, val string) {
	f, err := afs.Open(path)
	require.NoError(t, err)

	byt, err := io.ReadAll(f)
	require.NoError(t, err)

	assert.Equal(t, val, string(byt))
}
