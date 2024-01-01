package database

import (
	"context"
	"encoding/json"
	"github.com/andrewpillar/query"
	"github.com/brianvoe/gofakeit/v6"
	"github.com/jackc/pgx/v4/pgxpool"
	"go.mongodb.org/mongo-driver/bson"
	"log"
	"testing"
	"time"
)

type User struct {
	ID        int
	Email     string
	FirstName string
	LastName  string
	Active    bool
	Password  string
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (u *User) Primary() (string, any) {
	return "id", u.ID
}

func (u *User) Scan(fields []string, scan ScanFunc) error {
	return Scan(map[string]any{
		"id":         &u.ID,
		"email":      &u.Email,
		"first_name": &u.FirstName,
		"last_name":  &u.LastName,
		"password":   &u.Password,
		"active":     &u.Active,
		"created_at": &u.CreatedAt,
		"updated_at": &u.UpdatedAt,
	}, fields, scan)
}

func (u *User) Params() map[string]any {
	return map[string]any{
		"email":      u.Email,
		"first_name": u.FirstName,
		"last_name":  u.LastName,
		"password":   u.Password,
		"active":     u.Active,
		"created_at": u.CreatedAt,
		"updated_at": u.UpdatedAt,
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

	id, err := users.Create(context.TODO(), &User{
		ID:        0,
		Email:     gofakeit.Email(),
		FirstName: gofakeit.FirstName(),
		LastName:  gofakeit.LastName(),
		Password:  gofakeit.Password(true, true, true, true, false, 10),
		Active:    true,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	})

	if err != nil {
		t.Error(err)
	}

	if id != nil {
		t.Log("id: ", id)
	} else {
		t.Fatal("no id returned")
	}

	u, ok, err := users.Get(context.TODO(), query.Where("id", "=", query.Arg(id)))
	if err != nil {
		t.Error(err)
	}

	if !ok {
		t.Log("user not found")
	}

	data, _ := json.Marshal(u)
	t.Log(string(data))

	u.Password = gofakeit.Password(true, true, true, true, false, 10)

	if err = users.Update(context.TODO(), u); err != nil {
		if err != nil {
			t.Error(err)
		}
	}

	data, _ = json.Marshal(u)
	t.Log(string(data))

	list, err := users.All(context.TODO())
	if err != nil {
		t.Fatal(err)
	}
	for i := 0; i < list.Length(); i++ {
		u = list.Lookup(i)
		data, _ = json.Marshal(u)
		t.Log(string(data))
	}
}

func newUser(id int, email string) User {
	return User{
		ID:        id,
		Email:     email,
		FirstName: gofakeit.FirstName(),
		LastName:  gofakeit.LastName(),
		Password:  gofakeit.Password(true, true, true, true, false, 10),
		Active:    true,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}

func TestMongoDB(t *testing.T) {
	mongo, err := NewMongoDB("test.env")
	if err != nil {
		log.Fatal(err)
	}
	//disconnect when done
	defer mongo.Client.Disconnect(context.Background())

	email := gofakeit.Email()

	//create a new user
	user := newUser(gofakeit.Number(0, 100000), gofakeit.Email())
	err = mongo.Insert(context.TODO(), "users", user)
	if err != nil {
		log.Fatal(err)
	}

	//create a filter by email
	filter := bson.D{{Key: "email", Value: email}}
	var results []User

	err = mongo.Search(context.TODO(), "users", filter, &results)
	if err != nil {
		log.Fatal(err)
	}

	for _, u := range results {
		t.Logf("%+v", u)
	}

	// get all
	var all []User
	_ = mongo.All(context.TODO(), "users", &all)
	for _, u := range all {
		t.Logf("%+v", u)
	}
}

func TestMongoDB_Drop(t *testing.T) {
	mongo, err := NewMongoDB("test.env")
	if err != nil {
		log.Fatal(err)
	}
	//disconnect when done
	defer mongo.Client.Disconnect(context.Background())

	for i := 0; i < 10; i++ {
		user := newUser(gofakeit.Number(0, 100000), gofakeit.Email())
		err = mongo.Insert(context.TODO(), "users", user)
		if err != nil {
			log.Fatal(err)
		}
	}

	countBefore, err := mongo.Count(context.TODO(), "users")
	if err != nil {
		log.Fatal(err)
	}
	t.Log(countBefore)

	err = mongo.Drop(context.TODO(), "users")
	if err != nil {
		t.Fatal(err)
	}

	countAfter, err := mongo.Count(context.TODO(), "users")
	if err != nil {
		log.Fatal(err)
	}
	t.Log(countAfter)
	if countAfter > 0 {
		t.Fatal(countAfter)
	}

}
