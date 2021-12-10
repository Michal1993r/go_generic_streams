package stream

import (
	"fmt"
	"testing"
	"time"
)

func TestStream(t *testing.T) {
	var array []string
	for i := 0; i < 50; i++ {
		array = append(array, fmt.Sprintf("test%v", i))
	}

	s := FromIterable(array).
		Filter(func(s string) bool {
			return s != "test3"
		})

	slice := Map(s, func(t string) int {
		return len(t)
	}).
		Filter(func(s int) bool {
			return s > 5
		}).
		ToSlice()

	if len(slice) != 40 {
		t.Errorf("some elements were skipped, missing %d elements\n", 40 - len(slice))
	}
}

func TestStream_is_lazy(t *testing.T) {

	array := []string{"test", "test3", "test4"}
	wasCalled := false

	s := FromIterable(array).
		Filter(func(s string) bool {
			return s != "test3"
		})

	Map(s, func(t string) int {
		wasCalled = true
		return len(t)
	}).
		Filter(func(s int) bool {
			time.Sleep(1 * time.Second)
			return s > 4
		})

	if wasCalled {
		t.Error("stream is no longer lazy")
	}
}

func BenchmarkStream(b *testing.B) {
	var array []string
	for i := 0; i < 100000; i++ {
		array = append(array, fmt.Sprintf("test%v", i))
	}

	for i := 0; i < b.N; i++ {

		s := FromIterable(array).
			Filter(func(s string) bool {
				return s != "test3"
			})

		Map(s, func(t string) int {
			return len(t)
		}).
			Filter(func(s int) bool {
				return s > 6
			}).
			ForEach(func(t int) {
				fmt.Sprintln(t)
			})
	}
}

func BenchmarkForLoop(b *testing.B) {
	var array []string
	for i := 0; i < 100000; i++ {
		array = append(array, fmt.Sprintf("test%v", i))
	}

	for i := 0; i < b.N; i++ {
		for _, s := range array {
			if s != "test3" {
				l := len(s)
				if l > 6 {
					fmt.Sprintln(l)
				}
			}
		}
	}
}
