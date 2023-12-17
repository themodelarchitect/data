package store

import (
	"errors"
	"log"
	"strings"
	"testing"
	"time"
)

type job struct {
	Message   string
	CreatedAt string
	UpdatedAt string
}

func run[T job](repository Repository[job]) error {
	for i := 0; i < 10; i++ {
		now := time.Now().Format(time.RFC3339)
		j := job{
			Message:   "created",
			CreatedAt: now,
			UpdatedAt: now,
		}
		id, err := repository.Create(j)
		if err != nil {
			return err
		}

		j, err = repository.Get(id)
		if err != nil {
			return err
		}
		log.Println(j.Message, j.CreatedAt, j.UpdatedAt)
	}

	jobs := repository.GetAll()
	for k, v := range jobs {
		err := repository.Update(k, job{
			Message:   "updated",
			CreatedAt: v.CreatedAt,
			UpdatedAt: time.Now().Format(time.RFC3339),
		})
		if err != nil {
			return err
		}
	}

	jobs = repository.GetAll()
	for _, j := range jobs {
		log.Println(j.Message, j.CreatedAt, j.UpdatedAt)
		if !strings.Contains(j.Message, "updated") {
			return errors.New("not updated")
		}
	}

	jobs = repository.GetAll()
	for k := range jobs {
		err := repository.Delete(k)
		if err != nil {
			return err
		}
	}

	jobs = repository.GetAll()
	if len(jobs) != 0 {
		return errors.New("not all jobs deleted")
	}
	return nil
}

func TestRepository(t *testing.T) {
	err := run[job](NewInMemoryRepository[job]())
	if err != nil {
		t.FailNow()
	}
	err = run[job](NewFileRepository[job]())
	if err != nil {
		t.FailNow()
	}
}
