package database

import (
	"context"
	"encoding/json"
	"github.com/andrewpillar/query"
	"github.com/brianvoe/gofakeit/v6"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v4/pgxpool"
	"go.mongodb.org/mongo-driver/bson"
	"log"
	"testing"
	"time"
)

type User struct {
	ID        uuid.UUID `bson:"_id"`
	Email     string
	FirstName string
	LastName  string
	Active    bool
	Password  string
	CreatedAt time.Time `bson:"created_at,omitempty"`
	UpdatedAt time.Time `bson:"updated_at,omitempty"`
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

	users, err := NewPostgresDB[*User](".env", "users", func() *User {
		return &User{}
	})
	if err != nil {
		t.Fatal(err)
	}

	id, err := users.Create(context.TODO(), &User{
		ID:        uuid.New(),
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

func newUser(id uuid.UUID, email string) User {
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
	mongo, err := NewMongoDB[User](".env", "users")
	if err != nil {
		log.Fatal(err)
	}
	//disconnect when done
	defer mongo.Client.Disconnect(context.Background())

	email := gofakeit.Email()

	//create a new user
	user := newUser(uuid.New(), gofakeit.Email())
	err = mongo.Insert(context.TODO(), user)
	if err != nil {
		log.Fatal(err)
	}

	//create a filter by email
	filter := bson.D{{Key: "email", Value: email}}

	results, err := mongo.Search(context.TODO(), filter, nil)
	if err != nil {
		log.Fatal(err)
	}

	for _, u := range results {
		t.Logf("%+v", u)
	}

	// get all
	all, err := mongo.All(context.TODO(), nil)
	for _, u := range all {
		t.Logf("%+v", u)
	}
}

func TestMongoDB_FindByID(t *testing.T) {
	mongo, err := NewMongoDB[User](".env", "users")
	if err != nil {
		log.Fatal(err)
	}
	//disconnect when done
	defer mongo.Client.Disconnect(context.Background())

	email := gofakeit.Email()
	//create a new user
	user := newUser(uuid.New(), email)
	err = mongo.Insert(context.TODO(), user)
	if err != nil {
		log.Fatal(err)
	}

	//create a filter by email
	filter := bson.D{{Key: "email", Value: email}}

	t.Log(filter)
	results, err := mongo.Search(context.TODO(), filter, nil)
	if err != nil {
		t.Fatal(err)
	}

	res := results[0]
	t.Log(res)
	u, err := mongo.FindByID(context.TODO(), res.ID)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(u)
}

func TestMongoDB_Update(t *testing.T) {

}

func TestMongoDB_Drop(t *testing.T) {
	mongo, err := NewMongoDB[User](".env", "users")
	if err != nil {
		log.Fatal(err)
	}
	//disconnect when done
	defer mongo.Client.Disconnect(context.Background())

	for i := 0; i < 10; i++ {
		user := newUser(uuid.New(), gofakeit.Email())
		err = mongo.Insert(context.TODO(), user)
		if err != nil {
			log.Fatal(err)
		}
	}

	countBefore, err := mongo.Count(context.TODO())
	if err != nil {
		log.Fatal(err)
	}
	t.Log(countBefore)

	err = mongo.Drop(context.TODO())
	if err != nil {
		t.Fatal(err)
	}

	countAfter, err := mongo.Count(context.TODO())
	if err != nil {
		log.Fatal(err)
	}
	t.Log(countAfter)
	if countAfter > 0 {
		t.Fatal(countAfter)
	}

}
