package database

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/andrewpillar/query"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/jkittell/data/structures"
)

type Scan func(dest ...any) error

type Model interface {
	// Primary returns the name of the column used as the primary key for the Model, and the value for that column if any.
	Primary() (string, any)
	// Scan given the fields being scanned, and the function to call to actually perform the scan.
	Scan(fields []string, scan Scan) error
	// Params return the parameters of the Model to be used during create and update operations.
	Params() map[string]any
}

type SQLDB[M Model] struct {
	*pgxpool.Pool
	table string
	new   func() M
}

func (s *SQLDB[M]) Create(ctx context.Context, m M) (any, error) {
	p := m.Params()

	cols := make([]string, 0, len(p))
	vals := make([]any, 0, len(p))

	for k, v := range p {
		cols = append(cols, k)
		vals = append(vals, v)
	}

	primary, _ := m.Primary()

	q := query.Insert(
		s.table,
		query.Columns(cols...),
		query.Values(vals...),
		query.Returning(primary),
	)

	rows, err := s.Query(ctx, q.Build(), q.Args()...)

	if err != nil {
		_, key := m.Primary()
		return key, err
	}

	defer rows.Close()

	if !rows.Next() {
		_, key := m.Primary()
		if err := rows.Err(); err != nil {
			return key, err
		}
		return key, nil
	}

	if err := m.Scan(s.fields(rows), rows.Scan); err != nil {
		_, key := m.Primary()
		return key, nil
	}

	_, key := m.Primary()
	return key, nil
}

func (s *SQLDB[M]) Update(ctx context.Context, m M) error {
	p := m.Params()

	opts := make([]query.Option, 0, len(p))

	for k, v := range p {
		opts = append(opts, query.Set(k, query.Arg(v)))
	}

	col, id := m.Primary()

	opts = append(opts, query.Where(col, "=", query.Arg(id)))

	q := query.Update(s.table, opts...)

	if _, err := s.db.ExecContext(ctx, q.Build(), q.Args()...); err != nil {
		return err
	}
	return nil
}

func (s *SQLDB[M]) Delete(ctx context.Context, m M) error {
	col, id := m.Primary()

	q := query.Delete(s.table, query.Where(col, "=", query.Arg(id)))

	if _, err := s.db.ExecContext(ctx, q.Build(), q.Args()...); err != nil {
		return err
	}
	return nil
}

func (s *SQLDB[M]) Select(ctx context.Context, cols []string, opts ...query.Option) (*structures.Array[M], error) {
	opts = append([]query.Option{
		query.From(s.table),
	}, opts...)

	q := query.Select(query.Columns(cols...), opts...)

	rows, err := s.db.QueryContext(ctx, q.Build(), q.Args()...)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	mm := structures.NewArray[M]()

	for rows.Next() {
		m := s.new()

		if err := m.Scan(s.fields(rows), rows.Scan); err != nil {
			return nil, err
		}
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}
	return mm, nil
}

func (s *SQLDB[M]) All(ctx context.Context, opts ...query.Option) (*structures.Array[M], error) {
	return s.Select(ctx, []string{"*"}, opts...)
}

func (s *SQLDB[M]) Get(ctx context.Context, opts ...query.Option) (M, bool, error) {
	var zero M

	opts = append([]query.Option{
		query.From(s.table),
	}, opts...)

	q := query.Select(query.Columns("*"), opts...)

	rows, err := s.db.QueryContext(ctx, q.Build(), q.Args()...)

	if err != nil {
		return zero, false, err
	}

	defer rows.Close()

	if !rows.Next() {
		if err := rows.Err(); err != nil {
			return zero, false, err
		}
		return zero, false, nil
	}

	m := s.new()

	if err := m.Scan(s.fields(rows), rows.Scan); err != nil {
		return zero, false, err
	}
	return m, true, nil

}

func NewSQLDB[M Model](pool *pgxpool.Pool, table string) (*SQLDB[M], error) {
	table = ":memory:"
	db, err := sql.Open("sqlite3", table)
	if err != nil {
		return nil, err
	}
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("unable to reach database: %v", err)
	}

	return &SQLDB[M]{
		table: table,
	}, nil
}

func (s *SQLDB[M]) fields(rows *sql.Rows) []string {
	fields, _ := rows.Columns()
	return fields
}
