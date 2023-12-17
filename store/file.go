package store

import (
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"log"
	"math/rand"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"
)

type FileRepository[T any] struct {
	dir string
	mu  *sync.Mutex
}

func NewFileRepository[T any]() *FileRepository[T] {
	dir := path.Join(os.TempDir(), uuid.New().String())
	err := os.Mkdir(dir, 0755) // 0755 sets permissions for the directory
	if err != nil && !os.IsExist(err) {
		log.Fatal(err)
	}
	log.Println("new file repository: ", dir)
	return &FileRepository[T]{dir, &sync.Mutex{}}
}

func (r *FileRepository[T]) getDir() string {
	return r.dir
}

func (r *FileRepository[T]) getFilename(id int) string {
	return path.Join(r.dir, fmt.Sprintf("%d.json", id))
}

func (r *FileRepository[T]) Create(entity T) (int, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	id := rand.New(rand.NewSource(time.Now().UnixNano())).Int()
	filename := r.getFilename(id)
	err := saveJSON(filename, entity)
	return id, err
}

func (r *FileRepository[T]) Update(id int, entity T) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	return saveJSON(r.getFilename(id), entity)
}

func (r *FileRepository[T]) Get(id int) (T, error) {
	return loadJSON[T](r.getFilename(id))
}

func (r *FileRepository[T]) GetAll() map[int]T {
	entities := map[int]T{}
	files, err := os.ReadDir(r.dir)
	if err != nil {
		log.Println(err)
		return entities
	}

	for _, file := range files {
		if strings.HasSuffix(file.Name(), ".json") {
			filename := path.Join(r.dir, file.Name())
			idStr := file.Name()[:len(file.Name())-len(filepath.Ext(file.Name()))]
			id, err := strconv.Atoi(idStr)
			if err != nil {
				log.Println(err)
				continue
			}
			entity, err := loadJSON[T](filename)
			if err != nil {
				log.Println(err)
				continue
			}
			entities[id] = entity
		}
	}
	return entities
}

func (r *FileRepository[T]) Delete(id int) error {
	err := os.Remove(r.getFilename(id))
	if err != nil {
		return err
	}
	return nil
}

func loadJSON[T any](filename string) (T, error) {
	entity := *new(T)
	fi, err := os.Open(filename)
	if err != nil {
		return entity, err
	}
	defer fi.Close()

	decoder := json.NewDecoder(fi)
	err = decoder.Decode(&entity)
	if err != nil {
		return entity, err
	}

	return entity, nil
}

func saveJSON[T any](filename string, entity T) error {
	fi, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer fi.Close()

	encoder := json.NewEncoder(fi)
	err = encoder.Encode(entity)
	if err != nil {
		return err
	}
	return nil
}
