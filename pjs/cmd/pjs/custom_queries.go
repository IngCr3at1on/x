package main

import (
	"bytes"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/hashicorp/hcl/v2/hclsimple"
	"github.com/spf13/cobra"
)

type (
	customQueries struct {
		Queries []customQuery `hcl:"query,block"`

		parent *cobra.Command
	}

	customQuery struct {
		Use         string   `hcl:"use,label"`
		Description string   `hcl:"description,optional"`
		Template    string   `hcl:"template"`
		Args        []string `hcl:"args,optional"`

		t *template.Template
	}
)

func loadCustomQueries(cmd *cobra.Command, path string) error {
	abs, err := filepath.Abs(path)
	if err != nil {
		return fmt.Errorf("filepath.Abs(%s) -- %w", path, err)
	}

	if err = walk(cmd, abs); err != nil {
		return fmt.Errorf("filepath.Walk(%s) -- %w", abs, err)
	}

	return nil
}

func walk(cmd *cobra.Command, abs string) error {
	return filepath.WalkDir(abs, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			if !errors.Is(err, os.ErrNotExist) {
				return err
			}
		}

		if d.Name() == filepath.Base(abs) {
			return nil
		}

		if d.IsDir() {
			nc := cobra.Command{
				Use:           d.Name(),
				SilenceErrors: true,
				SilenceUsage:  true,
				PreRunE:       cobra.NoArgs,
			}
			cmd.AddCommand(&nc)
			if err = walk(&nc, path); err != nil {
				return fmt.Errorf("walk(%s) -- %w", path, err)
			}
		}

		if strings.HasSuffix(d.Name(), ".json") || strings.HasSuffix(d.Name(), ".hcl") {
			byt, err := os.ReadFile(path)
			if err != nil {
				return fmt.Errorf("os.ReadFile(%s) -- %w", path, err)
			}
			q := customQueries{
				parent: cmd,
			}
			if err = q.loadCustomQuery(d.Name(), byt); err != nil {
				return err
			}
		}

		return nil
	})
}

func (qs *customQueries) loadCustomQuery(filename string, byt []byte) error {
	if err := hclsimple.Decode(filename, byt, nil, qs); err != nil {
		return fmt.Errorf("hclsimple.Decode(%s) -- %w", filename, err)
	}

	for _, q := range qs.Queries {
		q.t = template.New(filename)
		var err error
		q.t, err = q.t.Parse(q.Template)
		if err != nil {
			return fmt.Errorf("error reading template from %s -- %w", filename, err)
		}

		qs.parent.AddCommand(&cobra.Command{
			Use:           q.Use,
			Long:          q.Description,
			SilenceErrors: true,
			PreRunE: func(cmd *cobra.Command, args []string) error {
				if len(args) != len(q.Args) {
					return cmd.Usage()
				}
				return nil
			},
			RunE: func(cmd *cobra.Command, args []string) error {
				m := make(map[string]string)
				for n, argKey := range q.Args {
					m[argKey] = args[n]
				}

				var buf bytes.Buffer
				if err := q.t.Execute(&buf, m); err != nil {
					return fmt.Errorf("q.t.Execute -- %w", err)
				}

				return transposeImpl(cmd.Context(), buf.String())
			},
		})
	}

	return nil
}
