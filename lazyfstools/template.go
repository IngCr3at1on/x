package lazyfstools

import (
	"html/template"
	"io/fs"
	"path/filepath"
	"strings"

	"github.com/spf13/afero"
)

const defaultSuffix = `.tpl`

type (
	TemplateOption interface{}
)

func RenderTemplates(afs afero.Fs, path string, target string, data any, _ ...TemplateOption) error {
	// FIXME: get suffix from options, also add better directory functionality.
	return Walk(afs, path, applyTemplate(target, defaultSuffix, data), SkipDirs())
}

func applyTemplate(target, suffix string, data any) ProcessFunc {
	return func(afs afero.Fs, path string, info fs.FileInfo) error {
		if strings.HasSuffix(path, suffix) {
			byt, err := afero.ReadFile(afs, path)
			if err != nil {
				return err
			}

			name := strings.TrimSuffix(path, suffix)
			name = filepath.Base(name)
			t := template.New(name)
			t, err = t.Parse(string(byt))
			if err != nil {
				return err
			}

			name = filepath.Join(target, name)
			f, err := afs.Create(name)
			if err != nil {
				return err
			}
			defer func() {
				_ = f.Close()
			}()

			return t.Execute(f, data)
		}

		return nil
	}
}
