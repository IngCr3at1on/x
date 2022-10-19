package lazyfstools

import (
	"errors"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/spf13/afero"
)

type ProcessFunc = func(afs afero.Fs, path string, info fs.FileInfo) error

// Walk a directory tree and execute the appropriate function.
func Walk(afs afero.Fs, abs string, fileProcessor ProcessFunc, dirProcessor ProcessFunc) error {
	return afero.Walk(afs, abs, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			if !errors.Is(err, os.ErrNotExist) {
				return err
			}
			return nil
		}

		if info.Name() == filepath.Base(abs) {
			return nil
		}

		if info.IsDir() {
			if dirProcessor != nil {
				return dirProcessor(afs, path, info)
			}

			return nil
		}

		return fileProcessor(afs, path, info)
	})
}
