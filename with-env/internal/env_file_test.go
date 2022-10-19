package internal_test

import (
	"io"
	"path/filepath"
	"testing"

	"github.com/ingcr3at1on/x/lazyfstools"
	"github.com/ingcr3at1on/x/with-env/internal"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetEnvFile(t *testing.T) {
	afs := afero.NewMemMapFs()
	abs := "env"
	err := afs.MkdirAll(abs, 0777)
	require.NoError(t, err)

	type testCase struct {
		path    string
		val     string
		expects func(tt *testing.T, r io.Reader, tc testCase)
	}

	defaultExpect := func(tt *testing.T, r io.Reader, tc testCase) {
		require.NotNil(tt, r)

		byt, err := io.ReadAll(r)
		require.NoError(tt, err)
		assert.Equal(tt, tc.val, string(byt))
	}

	for name, tc := range map[string]testCase{
		"env/test": {
			path: "test.env",
			val:  "FOO=BAR",
		},
		"alias-test": {
			path: "alias_test.env",
			val: `# with-env aliases: alias-test
VAL=foobar`,
		},
		"multi-alias-test": {
			path: "multi_alias_test.env",
			val: `# with-env aliases: foo,multi-alias-test
VAL=WAT?`,
		},
		"not-found": {
			path: "not_found_test.env",
			val:  "FOO=BAR",
			expects: func(tt *testing.T, r io.Reader, tc testCase) {
				assert.Nil(tt, r)
			},
		},
	} {
		t.Run(name, func(tt *testing.T) {
			err = lazyfstools.WriteFile(afs, filepath.Join(abs, tc.path), func(f afero.File) error {
				_, err := f.WriteString(tc.val)
				return err
			})
			require.NoError(tt, err)

			r, err := internal.GetEnvFile(afs, abs, name)
			require.NoError(tt, err)

			if tc.expects == nil {
				tc.expects = defaultExpect
			}

			tc.expects(tt, r, tc)
		})
	}
}
