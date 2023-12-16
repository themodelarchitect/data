package memory

import (
	"fmt"
	"strings"
	"testing"
)

type job struct {
	id      int
	message string
}

func (j *job) GetID() int {
	return j.id
}

func (i *job) SetID(id int) {
	i.id = id
}

func TestRepository_Create(t *testing.T) {
	db := New[*job]()

	for i := 0; i < 10; i++ {
		tm := &job{
			id:      0,
			message: fmt.Sprintf("%d", i),
		}
		_, err := db.Create(tm)
		if err != nil {
			t.Fatal(err)
		}
	}

	items := db.GetAll()

	if items.Length() != 10 {
		t.FailNow()
	}

}

func TestRepository_Update(t *testing.T) {
	db := New[*job]()

	for i := 0; i < 10; i++ {
		tm := &job{
			id:      0,
			message: fmt.Sprintf("%d", i),
		}
		_, err := db.Create(tm)
		if err != nil {
			t.Fatal(err)
		}
	}

	items := db.GetAll()
	for i := 0; i < items.Length(); i++ {
		j := items.Lookup(i)
		t.Log(j.id, j.message)
		err := db.Update(j.id, &job{
			id:      j.id,
			message: fmt.Sprintf("%d updated", j.id),
		})
		if err != nil {
			t.Fatal(err)
		}
	}

	items = db.GetAll()
	for i := 0; i < items.Length(); i++ {
		j := items.Lookup(i)
		t.Log(j.id, j.message)
		if !strings.Contains(j.message, "updated") {
			t.FailNow()
		}
	}
}

func TestRepository_Get(t *testing.T) {
	db := New[*job]()

	id, err := db.Create(&job{
		id:      1,
		message: fmt.Sprintf("%d", 1),
	})
	if err != nil {
		t.Fatal(err)
	}

	j, err := db.Get(id)
	if err != nil {
		t.Log(err)
	}
	if j.message != "1" {
		t.FailNow()
	}

	j, err = db.Get(100)
	if err != nil {
		t.Log(err)
	}
	t.Log(j)

}

func TestRepository_Delete(t *testing.T) {
	db := New[*job]()

	for i := 0; i < 10; i++ {
		tm := &job{
			id:      0,
			message: fmt.Sprintf("%d", i),
		}
		_, err := db.Create(tm)
		if err != nil {
			t.Fatal(err)
		}
	}

	items := db.GetAll()

	if items.Length() != 10 {
		t.FailNow()
	}

	for i := 0; i < items.Length(); i++ {
		j := items.Lookup(i)
		err := db.Delete(j.id)
		if err != nil {
			t.Fatal(err)
		}
	}

	items = db.GetAll()
	if items.Length() != 0 {
		t.FailNow()
	}
}
