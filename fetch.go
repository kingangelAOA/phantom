package phantom

import (
	"bytes"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"regexp"
	"strings"
	"time"
)

//Fetch fetch interface
type Fetch interface {
	request() interface{}
}

//DataPrepare interface data prepare
func (i *Interface) DataPrepare(cache *Cache) error {
	urlValue, err := parseDataByCache(i.Path, cache)
	if err != nil {
		return err
	}
	i.Path = urlValue
	bodyValue, err := parseDataByCache(i.Body, cache)
	if err != nil {
		return err
	}
	i.Body = bodyValue
	return nil
}

func (i *Interface) getContentByType() (io.Reader, error) {
	contentType, ok := i.Headers["Content-Type"]
	if !ok {
		return nil, fmt.Errorf("headers has not Content-Type")
	}
	if strings.Contains(contentType, "application/json") {
		return bytes.NewBuffer([]byte(i.Body)), nil
	}
	return nil, fmt.Errorf("body type %s did not surport", contentType)
}

//Request http request
func (i *Interface) Request() (*http.Response, error) {
	client := &http.Client{}
	body, err := i.getContentByType()
	if err != nil {
		return nil, err
	}
	dataLen := len(i.Hosts)
	s1 := rand.NewSource(time.Now().UnixNano())
	r1 := rand.New(s1)
	randNum := r1.Intn(dataLen)
	host := i.Hosts[randNum]
	req, err := http.NewRequest(i.Method, host+i.Path, body)
	if err != nil {
		return nil, err
	}
	for k, v := range i.Headers {
		req.Header.Set(k, v)
	}
	nowTime := time.Now().UnixNano() / int64(time.Millisecond)
	response, err := client.Do(req)
	i.Consuming = float64(time.Now().UnixNano()/int64(time.Millisecond) - nowTime)
	if err != nil {
		return nil, err
	}
	return response, nil
}

func parseDataByCache(data string, cache *Cache) (string, error) {
	rgx := regexp.MustCompile(`\$\{(.*?)\}`)
	matchs := rgx.FindAllStringSubmatch(data, -1)
	for _, match := range matchs {
		matchContent := match[0]
		key := match[1]
		value, ok := cache.Data[key]
		if !ok {
			return "", fmt.Errorf("cache has not %s when data prepare", key)
		}
		cacheResult, ok := value.(string)
		if !ok {
			return "", fmt.Errorf("convert %s (type interface {}) to type string error", value)
		}
		data = strings.Replace(data, matchContent, cacheResult, -1)
	}
	return data, nil
}
