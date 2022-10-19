package lazyfstools_test

import (
	"encoding/json"
	"io/fs"
	"path/filepath"
	"testing"

	"github.com/ingcr3at1on/x/lazyfstools"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWalk(t *testing.T) {
	afs := afero.NewMemMapFs()
	abs := "test"
	err := afs.MkdirAll(filepath.Join(abs, "dir"), 0777)
	require.NoError(t, err)

	err = lazyfstools.WriteFile(afs, filepath.Join(abs, "test.json"), func(f afero.File) error {
		_, err := f.WriteString(`{"str":"foobar"}`)
		return err
	})
	require.NoError(t, err)

	err = lazyfstools.WriteFile(afs, filepath.Join(abs, "dir", "test.json"), func(f afero.File) error {
		_, err := f.WriteString(`{"str":"barfoo"}`)
		return err
	})
	require.NoError(t, err)

	type st struct {
		Str string `json:"str"`
	}
	var sts []st

	err = lazyfstools.Walk(afs, abs,
		func(afs afero.Fs, path string, info fs.FileInfo) error {
			f, err := afs.Open(path)
			if err != nil {
				return err
			}

			var _st st
			if err = json.NewDecoder(f).Decode(&_st); err != nil {
				return err
			}

			sts = append(sts, _st)
			return nil
		}, nil)
	require.NoError(t, err)

	assert.Len(t, sts, 2)
	assert.Equal(t, "barfoo", sts[0].Str)
	assert.Equal(t, "foobar", sts[1].Str)
}
