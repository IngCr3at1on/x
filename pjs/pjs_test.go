package pjs

import (
	"context"
	"testing"

	"github.com/ingcr3at1on/x/pjs/internal/env"
	"github.com/jackc/pgtype"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func initTable(ctx context.Context, tb testing.TB, pool *pgxpool.Pool) {
	const sql = `
	CREATE TABLE IF NOT EXISTS pjs_test (
		s varchar PRIMARY KEY,
		n integer NOT NULL,
		b bool NOT NULL default FALSE
	);
	`
	_, err := pool.Exec(ctx, sql)
	require.NoError(tb, err)
}

func initPool(ctx context.Context, tb testing.TB) *pgxpool.Pool {
	dsn := env.ParseAlternateSettings()

	pc, err := pgxpool.ParseConfig(dsn)
	require.NoError(tb, err)

	pool, err := pgxpool.ConnectConfig(ctx, pc)
	require.NoError(tb, err)

	initTable(ctx, tb, pool)
	return pool
}

func simplePre(ctx context.Context, tt *testing.T, pool *pgxpool.Pool) {
	const sql = `
	INSERT INTO pjs_test (
		s, n
	) VALUES (
		$1, $2
	);`
	_, err := pool.Exec(ctx, sql, `foo`, 42)
	require.NoError(tt, err)
}

func TestTranspose(t *testing.T) {
	ctx := context.Background()
	pool := initPool(ctx, t)
	defer pool.Close()

	for name, tc := range map[string]struct {
		pre      func(ctx context.Context, tt *testing.T, pool *pgxpool.Pool)
		sql      string
		expected []Receiver
	}{
		"ok": {
			pre: simplePre,
			sql: `SELECT * FROM pjs_test;`,
			expected: []Receiver{
				{
					Name:        "s",
					DataTypeOID: pgtype.VarcharOID,
					Val: func() interface{} {
						return &pgtype.Varchar{
							String: "foo",
							Status: pgtype.Present,
						}
					}(),
				},
				{
					Name:        "n",
					DataTypeOID: pgtype.Int4OID,
					Val: func() interface{} {
						return &pgtype.Int4{
							Int:    42,
							Status: pgtype.Present,
						}
					}(),
				},
				{
					Name:        "b",
					DataTypeOID: pgtype.BoolOID,
					Val: func() interface{} {
						return &pgtype.Bool{
							Bool:   false,
							Status: pgtype.Present,
						}
					}(),
				},
			},
		},
	} {
		t.Run(name, func(tt *testing.T) {
			defer func() {
				_, err := pool.Exec(ctx, `TRUNCATE pjs_test`)
				require.NoError(tt, err)
			}()

			tc.pre(ctx, tt, pool)

			err := Transpose(
				context.Background(),
				pool,
				tc.sql,
				func(receivers []Receiver) error {
					assert.Equal(tt, tc.expected, receivers)
					return nil
				},
			)
			require.NoError(tt, err)
		})
	}
}
