package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/ingcr3at1on/x/pjs"
	"github.com/ingcr3at1on/x/sigctx"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgconn/stmtcache"
	"github.com/jackc/pgx/v4/pgxpool"
	flag "github.com/spf13/pflag"
)

const (
	dsnFK   = `dsn`
	typesFK = `show-types`
)

var (
	dsnF   *string
	typesF *bool
)

type output map[string]interface{}

func init() {
	dsnF = flag.StringP(dsnFK, `d`, ``, `a dsn string`)
	typesF = flag.BoolP(typesFK, `t`, false, `show types`)

	flag.Parse()
}

func main() {
	if err := sigctx.StartWith(func(ctx context.Context) error {
		byt, err := io.ReadAll(os.Stdin)
		if err != nil {
			return err
		}

		pc, err := pgxpool.ParseConfig(*dsnF)
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

		if err = pjs.Transpose(ctx, pool, string(byt), func(receivers []pjs.Receiver) error {
			out := make(output)
			for _, rec := range receivers {
				key := rec.Name
				if *typesF {
					key = fmt.Sprintf("%s (%T, %d)", key, rec.Val, rec.DataTypeOID)
				}
				out[key] = addToOutput(rec)
			}

			byt, err = json.Marshal(out)
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
	}); err != nil {
		log.Fatal(err)
	}
}
