package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"

	"github.com/ingcr3at1on/x/pjs"
	"github.com/ingcr3at1on/x/pjs/internal/env"
	"github.com/ingcr3at1on/x/sigctx"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgconn/stmtcache"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/spf13/cobra"
)

const (
	dsnFK   = `dsn`
	typesFK = `show-types`

	envQueriesDir = `PJS_CUSTOM_QUERIES`
)

var (
	dsnF   *string
	typesF *bool

	root = &cobra.Command{
		Use:           "pjs",
		SilenceErrors: true,
		SilenceUsage:  true,
		PreRunE:       cobra.NoArgs,
		RunE:          transpose,
	}
)

func init() {
	flags := root.PersistentFlags()
	dsnF = flags.StringP(dsnFK, `d`, ``, `a dsn string`)
	typesF = flags.BoolP(typesFK, `t`, false, `show types`)

	if err := initImpl(); err != nil {
		log.Fatal(err)
	}
}

func initImpl() error {
	queriesPath, ok := os.LookupEnv(envQueriesDir)
	if !ok {
		home, err := os.UserHomeDir()
		if err != nil {
			return err
		}
		queriesPath = filepath.Join(home, ".pjs", "queries")
	}

	if err := loadCustomQueries(root, queriesPath); err != nil {
		return err
	}
	return nil
}

func main() {
	if err := sigctx.StartWith(root.ExecuteContext); err != nil {
		log.Fatal(err)
	}
}

type output map[string]interface{}

func transpose(cmd *cobra.Command, _ []string) error {
	byt, err := io.ReadAll(os.Stdin)
	if err != nil {
		return err
	}

	return transposeImpl(cmd.Context(), string(byt))
}

func transposeImpl(ctx context.Context, query string) error {
	var dsn string
	if *dsnF != `` {
		dsn = *dsnF
	} else {
		dsn = env.ParseAlternateSettings()
	}

	pc, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return fmt.Errorf("pgxpool.ParseConfig -- %w", err)
	}
	pc.ConnConfig.BuildStatementCache = func(conn *pgconn.PgConn) stmtcache.Cache {
		return stmtcache.New(conn, stmtcache.ModeDescribe, 1024)
	}

	pool, err := pgxpool.ConnectConfig(ctx, pc)
	if err != nil {
		return fmt.Errorf("pgxpool.ConnectConfig -- %w", err)
	}
	defer pool.Close()

	if err = pjs.Transpose(ctx, pool, query, func(receivers []pjs.Receiver) error {
		out := make(output)
		for _, rec := range receivers {
			key := rec.Name
			if *typesF {
				key = fmt.Sprintf("%s (%T, %d)", key, rec.Val, rec.DataTypeOID)
			}
			out[key] = addToOutput(rec)
		}

		byt, err := json.Marshal(out)
		if err != nil {
			return fmt.Errorf("json.Marshal -- %w", err)
		}

		_, err = fmt.Println(string(byt))
		if err != nil {
			return fmt.Errorf("fmt.Println -- %w", err)
		}
		return nil
	}); err != nil {
		return fmt.Errorf("pjs.Transpose, %w", err)
	}

	return nil
}
