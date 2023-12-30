package database

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/andrewpillar/query"
	"github.com/brianvoe/gofakeit/v6"
	"github.com/jackc/pgx/v4/pgxpool"
	"log"
	"testing"
	"time"
)

type User struct {
	ID        int64
	Email     string
	Username  string
	Password  []byte
	CreatedAt time.Time
}

func (u *User) Primary() (string, any) {
	return "id", u.ID
}

func (u *User) Scan(fields []string, scan ScanFunc) error {
	return Scan(map[string]any{
		"id":         &u.ID,
		"email":      &u.Email,
		"username":   &u.Username,
		"password":   &u.Password,
		"created_at": &u.CreatedAt,
	}, fields, scan)
}

func (u *User) Params() map[string]any {
	return map[string]any{
		"email":      u.Email,
		"username":   u.Username,
		"password":   u.Password,
		"created_at": u.CreatedAt,
	}
}

func TestPosgresDB(t *testing.T) {
	// get the database connection URL.
	// usually, this is taken as an environment variable as in below commented out code
	// databaseURL = os.Getenv("DATABASE_URL")

	// for the time being, let's hard code it as follows.
	// ensure to change values as needed.
	databaseURL := "postgres://postgres:changeme@localhost:5432/postgres"

	// this returns connection pool
	pool, err := pgxpool.Connect(context.Background(), databaseURL)

	if err != nil {
		log.Fatalf("Unable to connect to database: %v\n", err)
	}
	defer pool.Close()

	users := NewPostgresDB[*User](pool, "users", func() *User {
		return &User{}
	})

	username := gofakeit.Username()
	id, err := users.Create(context.TODO(), &User{
		ID:        0,
		Email:     gofakeit.Email(),
		Username:  username,
		Password:  []byte(gofakeit.Password(true, true, true, true, false, 10)),
		CreatedAt: time.Now(),
	})

	if err != nil {
		t.Error(err)
	}

	if id != nil {
		t.Log("id: ", id)
	} else {
		t.Fatal("no id returned")
	}

	u, ok, err := users.Get(context.TODO(), query.Where("username", "=", query.Arg(username)))
	if err != nil {
		t.Error(err)
	}

	if !ok {
		fmt.Println("user not found")
	}

	data, _ := json.Marshal(u)
	t.Log(string(data))

	u.Password = []byte(gofakeit.Password(true, true, true, true, false, 10))

	if err = users.Update(context.TODO(), u); err != nil {
		if err != nil {
			t.Error(err)
		}
	}

	data, _ = json.Marshal(u)
	t.Log(string(data))
}
