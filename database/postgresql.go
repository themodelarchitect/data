package database

import (
	"context"
	"fmt"
	"github.com/andrewpillar/query"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/jkittell/data/structures"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
	"log"
)

type ScanFunc func(dest ...any) error

func Scan(structMap map[string]any, fields []string, scan ScanFunc) error {
	dest := make([]any, 0, len(fields))

	for _, fld := range fields {
		if p, ok := structMap[fld]; ok {
			dest = append(dest, p)
		}
	}
	return scan(dest...)
}

type Model interface {
	// Primary returns the name of the column used as the primary key for the Model, and the value for that column if any.
	Primary() (string, any)
	// Scan the given fields using the ScanFunc.
	Scan(fields []string, scan ScanFunc) error
	// Params return the parameters of the Model to be used during create and update operations.
	Params() map[string]any
}

type PosgresDB[M Model] struct {
	*pgxpool.Pool
	table string
	new   func() M
}

// NewPostgresDB takes a .env config file, table name and new func.
// .env file example
// POSTGRES_HOST=postgres
// POSTGRES_PORT=5432
// POSTGRES_USERNAME=postgres
// POSTGRES_PASSWORD=changeme
func NewPostgresDB[M Model](env, table string, new func() M) (PosgresDB[M], error) {
	var db PosgresDB[M]
	var databaseURL string
	viper.SetConfigFile(env)
	err := viper.ReadInConfig()
	if err != nil {
		return db, err
	}
	host := viper.Get("POSTGRES_HOST")
	port := viper.Get("POSTGRES_PORT")
	username := viper.Get("POSTGRES_USERNAME")
	password := viper.Get("POSTGRES_PASSWORD")

	host, ok := host.(string)
	if !ok {
		return db, errors.New("could not get postgres db host")
	}

	port, ok = port.(string)
	if !ok {
		return db, errors.New("could not get postgres db port")
	}

	username, ok = username.(string)
	if !ok {
		return db, errors.New("could not get postgres db username")
	}

	password, ok = password.(string)
	if !ok {
		return db, errors.New("could not get postgres db password")
	}

	//databaseURL := "postgres://postgres:changeme@localhost:5432/postgres"
	databaseURL = fmt.Sprintf("postgres://%s:%s@%s:%s/postgres", username, password, host, port)

	// this returns connection pool
	pool, err := pgxpool.Connect(context.Background(), databaseURL)

	if err != nil {
		log.Println(databaseURL)
		log.Fatalf("Unable to connect to database: %v\n", err)
	}

	db.Pool = pool
	db.table = table
	db.new = new

	return db, err
}

func (p *PosgresDB[M]) fields(rows pgx.Rows) []string {
	descriptions := rows.FieldDescriptions()
	fields := make([]string, 0, len(descriptions))

	for _, d := range descriptions {
		fields = append(fields, string(d.Name))
	}
	return fields
}

// Create a new entity M in the database and return the primary key.
func (p *PosgresDB[M]) Create(ctx context.Context, m M) (any, error) {
	var key any
	params := m.Params()

	cols := make([]string, 0, len(params))
	vals := make([]any, 0, len(params))

	for k, v := range params {
		cols = append(cols, k)
		vals = append(vals, v)
	}

	primary, _ := m.Primary()

	q := query.Insert(
		p.table,
		query.Columns(cols...),
		query.Values(vals...),
		query.Returning(primary),
	)

	rows, err := p.Query(ctx, q.Build(), q.Args()...)

	if err != nil {
		return key, err
	}

	defer rows.Close()

	if !rows.Next() {
		if err = rows.Err(); err != nil {
			return key, err
		}
	}

	if err = m.Scan(p.fields(rows), rows.Scan); err != nil {
		return key, err
	}

	_, key = m.Primary()
	return key, nil
}

func (p *PosgresDB[M]) Update(ctx context.Context, m M) error {
	params := m.Params()

	opts := make([]query.Option, 0, len(params))

	for k, v := range params {
		opts = append(opts, query.Set(k, query.Arg(v)))
	}

	col, id := m.Primary()

	opts = append(opts, query.Where(col, "=", query.Arg(id)))

	q := query.Update(p.table, opts...)

	if _, err := p.Exec(ctx, q.Build(), q.Args()...); err != nil {
		return err
	}
	return nil
}

func (p *PosgresDB[M]) Delete(ctx context.Context, m M) error {
	col, id := m.Primary()

	q := query.Delete(p.table, query.Where(col, "=", query.Arg(id)))

	if _, err := p.Exec(ctx, q.Build(), q.Args()...); err != nil {
		return err
	}
	return nil
}

func (p *PosgresDB[M]) Select(ctx context.Context, cols []string, opts ...query.Option) (*structures.Array[M], error) {
	opts = append([]query.Option{
		query.From(p.table),
	}, opts...)

	q := query.Select(query.Columns(cols...), opts...)

	rows, err := p.Query(ctx, q.Build(), q.Args()...)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	models := structures.NewArray[M]()

	for rows.Next() {
		m := p.new()
		if err = m.Scan(p.fields(rows), rows.Scan); err != nil {
			return nil, err
		}
		models.Push(m)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}
	return models, nil
}

func (p *PosgresDB[M]) All(ctx context.Context) (*structures.Array[M], error) {
	q := query.Select(query.Columns("*"), query.From(p.table))

	//log.Println(q.Build())

	rows, err := p.Query(ctx, q.Build(), q.Args()...)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	models := structures.NewArray[M]()

	for rows.Next() {
		m := p.new()
		if err = m.Scan(p.fields(rows), rows.Scan); err != nil {
			return nil, err
		}
		models.Push(m)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}
	return models, nil
}

func (p *PosgresDB[M]) Get(ctx context.Context, opts ...query.Option) (M, bool, error) {
	var zero M

	opts = append([]query.Option{
		query.From(p.table),
	}, opts...)

	q := query.Select(query.Columns("*"), opts...)

	rows, err := p.Query(ctx, q.Build(), q.Args()...)

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

	m := p.new()

	if err := m.Scan(p.fields(rows), rows.Scan); err != nil {
		return zero, false, err
	}
	return m, true, nil

}
