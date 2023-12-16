package array

import (
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"log"
	"testing"
)

type item struct {
	Id     uuid.UUID
	Name   string
	Number int
}

func TestArray_Push(t *testing.T) {
	arr := New[int]()
	arr.Push(1)
	arr.Push(2)
	arr.Push(3)
	if arr.Length() != 3 {
		t.FailNow()
	}
}

func TestArray_1(t *testing.T) {
	arr := New[int]()
	arr.Push(1)
	arr.Push(2)
	arr.Push(3)

	for i := 0; i < 3; i++ {
		val := arr.Lookup(i)
		if val != i+1 {
			t.FailNow()
		}
	}
}

func TestArray_2(t *testing.T) {
	expected := []string{"a", "b", "c"}

	arr := New[string]()
	arr.Push("a")
	arr.Push("b")
	arr.Push("c")

	for i := 0; i < 3; i++ {
		val := arr.Lookup(i)
		if val != expected[i] {
			t.FailNow()
		}
	}
}

func TestArray_Merge(t *testing.T) {
	arr1 := New[int]()
	arr2 := New[int]()

	arr1.Push(1)
	arr1.Push(2)
	arr1.Push(3)

	arr2.Push(4)
	arr2.Push(5)
	arr2.Push(6)

	mergedArray := arr1.Merge(arr2)

	for i := 0; i < 6; i++ {
		n := mergedArray.Lookup(i)
		t.Log(n)
		if n != i+1 {
			t.FailNow()
		}
	}
}

func TestArray_3(t *testing.T) {
	expected := []string{"c", "b", "a"}

	arr := New[string]()
	arr.Push("a")
	arr.Push("b")
	arr.Push("c")
	arr.Reverse()

	for i := 0; i < 3; i++ {
		val := arr.Lookup(i)
		if val != expected[i] {
			t.FailNow()
		}
	}
}

func TestArray_4(t *testing.T) {
	arr := New[int]()

	arr.Push(1)
	arr.Push(2)
	arr.Push(3)
	value := arr.Pop()
	if value != 3 {
		t.FailNow()
	}
	value = arr.Pop()
	if value != 2 {
		t.FailNow()
	}
	value = arr.Pop()
	if value != 1 {
		t.FailNow()
	}
	count := arr.Length()
	if count != 0 {
		t.FailNow()
	}
}

func TestArray_Copy(t *testing.T) {
	arr := New[int]()
	arr.Push(1)
	arr.Push(2)
	arr.Push(3)
	arrCopy := arr.Copy()

	for i := 0; i < 3; i++ {
		x := arr.Lookup(i)
		y := arrCopy.Lookup(i)
		log.Println(x, y)
		if x != y {
			t.FailNow()
		}
	}
}

func TestArray_Copy2(t *testing.T) {
	arr := New[int]()
	arr.Push(1)
	arr.Push(2)
	arr.Push(3)
	arrCopy := arr.Copy()

	for i := 0; i < arr.length; i++ {
		arrCopy.Set(i, i+3)
	}

	for i := 0; i < arr.length; i++ {
		x := arr.Lookup(i)
		y := arrCopy.Lookup(i)
		log.Println(x, y)
	}
}

func TestArray_Copy3(t *testing.T) {
	arr := New[item]()

	item1 := item{
		Id:     uuid.New(),
		Name:   "original 1",
		Number: 1,
	}

	item2 := item{
		Id:     uuid.New(),
		Name:   "original 2",
		Number: 2,
	}

	item3 := item{
		Id:     uuid.New(),
		Name:   "original 3",
		Number: 3,
	}

	arr.Push(item1)
	arr.Push(item2)
	arr.Push(item3)

	arrCopy := arr.Copy()

	copy1 := arrCopy.Lookup(0)
	copy2 := arrCopy.Lookup(1)
	copy3 := arrCopy.Lookup(2)

	copy1.Id = uuid.New()
	copy1.Name = "copy 1"
	copy1.Number = 4

	copy2.Id = uuid.New()
	copy2.Name = "copy 2"
	copy2.Number = 5

	copy3.Id = uuid.New()
	copy3.Name = "copy 3"
	copy3.Number = 6

	// need to set to update the array
	arrCopy.Set(0, copy1)
	arrCopy.Set(1, copy2)
	arrCopy.Set(2, copy3)

	for i := 0; i < arr.Length(); i++ {
		log.Println(arr.Lookup(i).Id, arr.Lookup(i).Name, arr.Lookup(i).Number)
	}

	for i := 0; i < arrCopy.Length(); i++ {
		log.Println(arrCopy.Lookup(i).Id, arrCopy.Lookup(i).Name, arrCopy.Lookup(i).Number)
	}
}

func setItem(items *Array[item]) {
	for i := 0; i < items.Length(); i++ {
		item := items.Lookup(i)
		item.Id = uuid.New()
		item.Name = fmt.Sprintf("setting %d", i)
		item.Number = i + 3
		items.Set(i, item)
	}
}

func TestArray_Copy4(t *testing.T) {
	arr := New[item]()

	item1 := item{
		Id:     uuid.New(),
		Name:   "original 1",
		Number: 1,
	}

	item2 := item{
		Id:     uuid.New(),
		Name:   "original 2",
		Number: 2,
	}

	item3 := item{
		Id:     uuid.New(),
		Name:   "original 3",
		Number: 3,
	}

	arr.Push(item1)
	arr.Push(item2)
	arr.Push(item3)

	arrCopy := arr.Copy()

	for i := 0; i < 3; i++ {
		setItem(arrCopy)
	}

	arr = arr.Merge(arrCopy)

	for i := 0; i < arr.Length(); i++ {
		log.Println(arr.Lookup(i).Id, arr.Lookup(i).Name, arr.Lookup(i).Number)
	}
}

func TestArray_Set(t *testing.T) {
	arr := New[int]()
	arr.Push(1)
	arr.Push(2)
	arr.Push(3)

	arr.Set(0, 4)
	arr.Set(1, 5)
	arr.Set(2, 6)

	for i := 0; i < 3; i++ {
		x := arr.Lookup(i)
		log.Println(x)
		if x != (i + 4) {
			t.FailNow()
		}
	}

}

func TestArray_UnmarshalJSON(t *testing.T) {
	jsonMsg := []byte(`[0,1,2,3,4,5,6,7,8,9]`)
	arr := New[int]()

	err := json.Unmarshal(jsonMsg, &arr)
	if err != nil {
		fmt.Println(err)
		t.FailNow()
	}

	fmt.Printf("%#v\n", arr)
}

func TestArray_MarshalJSON(t *testing.T) {
	arr := New[int]()
	for i := 0; i < 10; i++ {
		arr.Push(i)
	}

	b, err := json.Marshal(arr)
	if err != nil {
		t.FailNow()
	}
	fmt.Println(string(b))
}

func BenchmarkArray_Push(b *testing.B) {
	arr := New[int]()
	for n := 0; n < b.N; n++ {
		arr.Push(n)
	}
}

func BenchmarkArray_Copy(b *testing.B) {
	arr := New[int]()
	arr.Push(1)
	arr.Push(2)
	arr.Push(3)
	for n := 0; n < b.N; n++ {
		arr.Copy()
	}
}
