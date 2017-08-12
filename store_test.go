package phantom

import (
	"fmt"
	"io/ioutil"
	"testing"
)

func TestSave(t *testing.T) {
	data, err := ioutil.ReadFile("example.json")

	if err != nil {
		fmt.Print(err)
	}
	scenes, err := JSONToScenes(data)
	if err != nil {
		fmt.Print(err)
	}
	for _, scene := range scenes {
		for _, in := range scene.Interfaces {
			cache := NewCache()
			for _, store := range in.Stores {
				if err := store.Save(`{"message": "test", "password": "111111", "status": "teststatus"}`, cache); err != nil {
					t.Error(err)
				}
			}
			fmt.Println(cache)
		}
	}
}
