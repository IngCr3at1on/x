package internal

import (
	"bufio"
	"bytes"
	"io"
	"io/fs"
	"strings"

	"github.com/ingcr3at1on/x/lazyfstools"
	"github.com/spf13/afero"
)

// Walk a directory tree to acquire all index files.
func GetEnvFile(afs afero.Fs, abs, name string) (io.Reader, error) {
	returnF := func() io.Reader { return nil }
	err := lazyfstools.Walk(afs, abs,
		func(afs afero.Fs, path string, _ fs.FileInfo) error {
			f, err := afs.Open(path)
			if err != nil {
				return err
			}

			var buf bytes.Buffer
			aliases, err := scanAliases(io.TeeReader(f, &buf))
			if err != nil {
				return err
			}
			aliases = append(aliases, strings.TrimSuffix(path, ".env"))

			for _, alias := range aliases {
				if alias == name {
					returnF = func() io.Reader { return &buf }
					return nil
				}
			}

			return nil
		}, nil)
	if err != nil {
		return nil, err
	}

	return returnF(), nil
}

func scanAliases(r io.Reader) ([]string, error) {
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		t := scanner.Text()
		prefix := "# with-env aliases: "
		if strings.HasPrefix(t, prefix) {
			t = strings.TrimPrefix(t, prefix)
			return strings.Split(t, ","), nil
		}
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}

	// No aliases found.
	return nil, nil
}
