package pjs

import (
	"context"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v4/pgxpool"
)

// Receiver is a value receiver when reading postgres data.
type Receiver struct {
	Name        string
	DataTypeOID uint32
	Val         interface{}
}

// Transpose reads arbitrary data out of pgx.Rows.
func Transpose(
	ctx context.Context,
	pool *pgxpool.Pool,
	sql string,
	transposeF func(receivers []Receiver) error,
) error {
	sql = strings.TrimSpace(sql)

	batch, queries, err := scanQueries(sql)
	if err != nil {
		return err
	}

	res := pool.SendBatch(ctx, &batch)
	defer res.Close()

	for n := range queries {
		// Wrap in a function to allow defer to run.
		if err = func() error {
			rows, err := res.Query()
			if err != nil {
				return fmt.Errorf("pool.Query[%d] -- %w", n, err)
			}
			defer rows.Close()

			descs := rows.FieldDescriptions()
			receivers := func() ([]Receiver, []interface{}) {
				l := len(descs)
				recs := make([]Receiver, l)
				vals := make([]interface{}, l)
				for n, desc := range rows.FieldDescriptions() {
					r := Receiver{
						Name:        string(desc.Name),
						DataTypeOID: desc.DataTypeOID,
						Val:         getOb(desc.DataTypeOID),
					}
					recs[n] = r
					vals[n] = r.Val
				}
				return recs, vals
			}

			for rows.Next() {
				recievers, vals := receivers()
				if err := rows.Scan(vals...); err != nil {
					return err
				}

				if err := transposeF(recievers); err != nil {
					return err
				}
			}

			if err = rows.Err(); err != nil {
				return fmt.Errorf("rows.Err() -- %w", err)
			}

			return nil
		}(); err != nil {
			return err
		}
	}
	return nil
}
