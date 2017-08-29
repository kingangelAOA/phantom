package phantom

import "encoding/json"

//Scene 场景
type Scene struct {
	Name       string `validate:"nonzero"`
	Interfaces []Interface
	RunConfig  RunConfig
}

//Interface 场景中的接口
type Interface struct {
	Name      string `validate:"nonzero"`
	Hosts     []string
	Path      string `validate:"nonzero"`
	Headers   map[string]string
	Method    string `validate:"nonzero"`
	Body      string `validate:"nonzero"`
	Consuming float64
	Stores    []Store
	TestData  *TestData
	Asserts   []Assert
}

//Assert assert config
type Assert struct {
	Type     string
	JSONPath string
	Expect   string
}

//TestData test data
type TestData struct {
	Config    string `validate:"nonzero"`
	Separator string `validate:"nonzero"`
	Path      string `validate:"nonzero"`
	Type      string
	index     int
	Data      []map[string]string
}

// Store save some response data to global
type Store struct {
	Type     string `validate:"nonzero"`
	ToKey    string `validate:"nonzero"`
	JSONPath string `validate:"nonzero"`
	Value    string
}

//Cache cache interface response data
type Cache struct {
	Data map[string]interface{}
}

//RunConfig run config
type RunConfig struct {
	Type      uint8
	Time      uint16 `validate:"min=1"`
	ThreadNum uint16 `validate:"min=1"`
	UserNum   uint16
	Wait      bool
}

//NewCache init Cache
func NewCache() *Cache {
	var c Cache
	c.Data = make(map[string]interface{})
	return &c
}

//JSONToScenes json to struct
func JSONToScenes(b []byte) ([]*Scene, error) {
	var scenes []*Scene
	err := json.Unmarshal(b, &scenes)
	if err != nil {
		return nil, err
	}
	for _, scene := range scenes {
		if err := ValidateStruct(scene); err != nil {
			return nil, err
		}
		if err := ValidateStruct(scene.RunConfig); err != nil {
			return nil, err
		}
		for _, in := range scene.Interfaces {
			if err := ValidateStruct(in); err != nil {
				return nil, err
			}
			for _, store := range in.Stores {
				if err := ValidateStruct(store); err != nil {
					return nil, err
				}
			}
		}
	}
	for _, scene := range scenes {
		for _, in := range scene.Interfaces {
			if in.TestData != nil {
				if err := in.TestData.initTestData(); err != nil {
					return nil, err
				}
			}
		}
	}
	return scenes, nil
}
