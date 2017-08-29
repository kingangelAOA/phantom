package phantom

import (
	"errors"
	"fmt"
	"io/ioutil"
	"time"
)

var log Log

//Run run one Scene
func (s *Scene) Run() {
	var p Pool
	log.Scene = s.Name
	s.choiceModel(p)
}

func (s *Scene) choiceModel(p Pool) {
	if s.RunConfig.Type == RunTypeTime {
		s.controlByTime(p)
	}
}

func (s *Scene) controlByTime(p Pool) error {
	r := s.RunConfig
	num := r.ThreadNum
	runTime := r.Time
	endFlag := time.After(time.Duration(runTime) * time.Second)
	if runTime == 0 {
		return errors.New("time can not 0")
	}
	p.Init(num, num)
	log.ThreadNum = num
	logIns := initLogIn(s.Interfaces)
	for _, v := range logIns {
		log.Interfaces = append(log.Interfaces, v)
	}
	p.addTask(func() error {
		if err := runInterfaces(s.Interfaces, logIns); err != nil {
			return err
		}
		return nil
	})
	p.start()
	log.RealTimeData(&p)
	for {
		select {
		case <-endFlag:
			break
		case err := <-p.Result:
			fmt.Println(err)
		}
	}
}

func runInterfaces(ins []Interface, logIns map[string]*In) error {
	cache := NewCache()
	for _, in := range ins {
		logIn, _ := logIns[in.Name]
		if in.TestData != nil {
			if err := in.TestData.updateCache(cache); err != nil {
				return err
			}
		}
		// fmt.Println(cache)
		if err := in.DataPrepare(cache); err != nil {
			return err
		}
		res, err := in.Request()
		if err != nil {
			return err
		}
		MuConsume.Lock()
		logIn.Consumes = append(logIn.Consumes, in.Consuming)
		MuConsume.Unlock()
		if res.StatusCode == 200 {
			body, err := ioutil.ReadAll(res.Body)
			if err != nil {
				return err
			}
			for _, assert := range in.Asserts {
				if err := assert.Do(body); err != nil {
					return err
				}
			}
			for _, store := range in.Stores {
				err := store.Save(body, cache)
				if err != nil {
					return err
				}
			}
		} else {
			body, err := ioutil.ReadAll(res.Body)
			if err != nil {
				return err
			}
			return fmt.Errorf("request error, status: %d, body: %s", res.StatusCode, string(body))
		}
	}
	return nil
}

func initLogIn(ins []Interface) map[string]*In {
	logIns := map[string]*In{}
	for _, in := range ins {
		logIn := &In{
			Name:     in.Name,
			Consumes: []float64{},
		}
		logIns[in.Name] = logIn
	}
	return logIns
}
