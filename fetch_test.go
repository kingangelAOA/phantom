package phantom

import (
	"fmt"
	"io/ioutil"
	"testing"
)

func TestRequest(t *testing.T) {
	data, err := ioutil.ReadFile("example.json")

	if err != nil {
		t.Error(err)
	}
	scenes, err := JSONToScenes(data)
	if err != nil {
		t.Error(err)
	}

	for _, scene := range scenes {
		ins := scene.Interfaces
		cache := NewCache()
		for _, in := range ins {
			err := in.DataPrepare(cache)
			if err != nil {
				t.Error(err)
			}
			res, err := in.Request()
			if err != nil {
				t.Error(err)
			}
			if res.StatusCode == 200 {
				body, err := ioutil.ReadAll(res.Body)
				if err != nil {
					t.Error(err)
				}
				for _, store := range in.Stores {
					err := store.Save(body, cache)
					if err != nil {
						t.Error(err)
					}
				}
			}
		}
		fmt.Println(cache)
	}
}
