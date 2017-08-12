package phantom

import (
	"errors"
	"reflect"

	validator "gopkg.in/validator.v2"
)

func matchMethod(v interface{}, param string) error {
	Methods := Strings{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS"}
	st := reflect.ValueOf(v)
	if st.Kind() != reflect.String {
		return errors.New("method only validates strings")
	}
	if !Methods.Contains(st.String()) {
		return errors.New("value is not GET POST PUT PATCH DELETE HEAD OPTIONS")
	}
	return nil
}

func addValidationFunc() error {
	if err := validator.SetValidationFunc("matchMethod", matchMethod); err != nil {
		return err
	}
	return nil
}

// ValidateStruct Validate Struct
func ValidateStruct(t interface{}) error {
	if err := addValidationFunc(); err != nil {
		return err
	}
	if err := validator.Validate(t); err != nil {
		return err
	}
	return nil
}
