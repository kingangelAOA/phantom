package phantom

import "fmt"

//Save save data
func (s *Store) Save(originData []byte, cache *Cache) error {
	if s.Type == StoreTypeJSONPath {
		result, err := GetValueByJSONPath(s.JSONPath, originData)
		if err != nil {
			return err
		}
		cache.Data[s.ToKey] = result
	} else if s.Type == StoreTypeCommon {
		cache.Data[s.ToKey] = s.Value
	} else {
		return fmt.Errorf("store type \"%s\" did not surport, please check config json or params", s.Type)
	}
	return nil
}
