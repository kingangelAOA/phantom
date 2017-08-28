package phantom

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"
)

//initTestData init test data
func (t *TestData) initTestData() error {
	if err := validateData(t.Config, t.Separator); err != nil {
		return err
	}
	file, err := os.Open(t.Path)
	defer file.Close()
	if err != nil {
		return err
	}
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lineStr := scanner.Text()
		data := map[string]string{}
		if err := validateData(lineStr, t.Separator); err != nil {
			return err
		}
		lineStrSlice := strings.Split(lineStr, t.Separator)
		configSlice := strings.Split(t.Config, t.Separator)
		for index, key := range configSlice {
			if index > len(lineStrSlice)-1 {
				continue
			}
			data[key] = lineStrSlice[index]
		}
		t.Data = append(t.Data, data)
	}
	return nil
}

func (t *TestData) updateCache(c *Cache) error {
	dataLen := len(t.Data)
	if t.Type == TestDataAlphabetical {
		if t.index == dataLen-1 {
			t.index = 0
		}
		for k, v := range t.Data[t.index] {
			if err := validCache(c, k); err != nil {
				return err
			}
			c.Data[k] = v
		}
		t.index++
	} else if t.Type == TestDataRandom {
		s1 := rand.NewSource(time.Now().UnixNano())
		r1 := rand.New(s1)
		randNum := r1.Intn(dataLen)
		data := t.Data[randNum]
		for k, v := range data {
			if err := validCache(c, k); err != nil {
				return err
			}
			c.Data[k] = v
		}
	} else {
		return fmt.Errorf("test data type error, '%s' is not exist", t.Type)
	}
	return nil
}

func validCache(c *Cache, k string) error {
	if _, ok := c.Data[k]; ok {
		return fmt.Errorf("test data update to cache error, '%s' is exist in cahce", k)
	}
	return nil
}

func validateData(config, separator string) error {
	if !strings.Contains(config, separator) {
		return fmt.Errorf("config: %s did not contains %s", config, separator)
	}
	return nil
}
