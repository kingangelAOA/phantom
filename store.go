package phantom

import (
	"bytes"
	"encoding/json"
	"fmt"

	JSONPath "k8s.io/client-go/util/jsonpath"
)

//Save save data
func (s *Store) Save(originData string, cache *Cache) error {
	if s.Type == StoreTypeJSONPath {
		var pointsData interface{}
		if err := json.Unmarshal([]byte(originData), &pointsData); err != nil {
			return err
		}
		j := JSONPath.New("")
		if err := j.Parse(s.JSONPath); err != nil {
			return err
		}
		buf := new(bytes.Buffer)
		if err := j.Execute(buf, &pointsData); err != nil {
			return fmt.Errorf("store JSONPath: %s did not found, please check path", s.JSONPath)
		}
		cache.Data[s.ToKey] = buf.String()
	} else if s.Type == StoreTypeCommon {
		cache.Data[s.ToKey] = s.Value
	} else {
		return fmt.Errorf("store type \"%s\" did not surport, please check config json or params", s.Type)
	}
	return nil
}
