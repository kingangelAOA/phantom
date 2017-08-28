package phantom

import "fmt"

//Do assert response
func (a *Assert) Do(r []byte) error {
	if a.Type == AssertJSON {
		result, err := GetValueByJSONPath(a.JSONPath, r)
		if err != nil {
			return err
		}
		if result != a.Expect {
			return fmt.Errorf("assert error; %s(%s) != %s", result, a.JSONPath, a.Expect)
		}
	}
	return nil
}
