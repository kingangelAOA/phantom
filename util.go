package phantom

import (
	"bytes"
	"encoding/json"
	"fmt"
	"reflect"
	"runtime"
	"strconv"
	"strings"

	JSONPath "k8s.io/client-go/util/jsonpath"
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

//GetValueByJSONPath get value by jsonpath
func GetValueByJSONPath(path string, originData []byte) (string, error) {
	var pointsData interface{}
	if err := json.Unmarshal(originData, &pointsData); err != nil {
		return "", err
	}
	j := JSONPath.New("")
	if err := j.Parse(path); err != nil {
		return "", err
	}
	buf := new(bytes.Buffer)
	if err := j.Execute(buf, &pointsData); err != nil {
		return "", fmt.Errorf("store JSONPath: %s did not found, please check path", path)
	}
	return buf.String(), nil
}

// Any formats any value as a string.
func Any(value interface{}) string {
	return formatAtom(reflect.ValueOf(value))
}

// formatAtom formats a value without inspecting its internal structure.
func formatAtom(v reflect.Value) string {
	switch v.Kind() {
	case reflect.Invalid:
		return "invalid"
	case reflect.Int, reflect.Int8, reflect.Int16,
		reflect.Int32, reflect.Int64:
		return strconv.FormatInt(v.Int(), 10)
	case reflect.Uint, reflect.Uint8, reflect.Uint16,
		reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return strconv.FormatUint(v.Uint(), 10)
	// ...floating-point and complex cases omitted for brevity...
	case reflect.Bool:
		return strconv.FormatBool(v.Bool())
	case reflect.String:
		return strconv.Quote(v.String())
	case reflect.Chan, reflect.Func, reflect.Ptr, reflect.Slice, reflect.Map:
		return v.Type().String() + " 0x" +
			strconv.FormatUint(uint64(v.Pointer()), 16)
	default: // reflect.Array, reflect.Struct, reflect.Interface
		return v.Type().String() + " value"
	}
}
