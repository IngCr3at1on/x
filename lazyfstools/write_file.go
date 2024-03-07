package lazyfstools

import (
	"fmt"

	"github.com/spf13/afero"
)

// WriteFile is a quick and dirty way to write/sync/close a file.
func WriteFile(afs afero.Fs, path string, writers ...func(f afero.File) error) error {
	f, err := afs.Create(path)
	if err != nil {
		return fmt.Errorf("afs.Create: %w", err)
	}

	for _, writer := range writers {
		if err = writer(f); err != nil {
			return err
		}
	}

	if err = f.Sync(); err != nil {
		return err
	}

	return f.Close()
}
