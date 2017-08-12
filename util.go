package phantom

import (
	"fmt"
	"runtime"
	"strconv"
	"strings"
)

//Strings string slice
type Strings []string

//Ints int slice
type Ints []int

//Ints64 ints64 slice
type Ints64 []int64

//Contains ints Contains
func (sl Ints) Contains(v int) bool {
	for _, vv := range sl {
		if vv == v {
			return true
		}
	}
	return false
}

//Contains int64 Contains
func (sl Ints64) Contains(v int64) bool {
	for _, vv := range sl {
		if vv == v {
			return true
		}
	}
	return false
}

//Contains string Contains
func (sl Strings) Contains(v string) bool {
	for _, vv := range sl {
		if vv == v {
			return true
		}
	}
	return false
}

//GoID GoroutineId
func GoID() int {
	var buf [64]byte
	n := runtime.Stack(buf[:], false)
	idField := strings.Fields(strings.TrimPrefix(string(buf[:n]), "goroutine "))[0]
	id, err := strconv.Atoi(idField)
	if err != nil {
		panic(fmt.Sprintf("cannot get goroutine id: %v", err))
	}
	return id
}
