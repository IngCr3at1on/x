package env_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/ingcr3at1on/x/lazyfstools"
	env "github.com/ingcr3at1on/x/with-env"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoadFrom(t *testing.T) {
	afs := afero.NewMemMapFs()
	testCase(t, afs, "env")
}

func testCase(t *testing.T, afs afero.Fs, abs string) {
	err := afs.MkdirAll(abs, 0777)
	require.NoError(t, err)

	err = lazyfstools.WriteFile(afs, filepath.Join("env", "test.env"), func(f afero.File) error {
		_, err := f.WriteString(`FOO=BAR`)
		return err
	})
	require.NoError(t, err)

	err = env.LoadFrom(afs, "env", "env/test")
	require.NoError(t, err)

	v, ok := os.LookupEnv("FOO")
	assert.True(t, ok)
	assert.Equal(t, "BAR", v)
}
