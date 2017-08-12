package phantom

import (
	"errors"
	"io/ioutil"
	"time"
)

var log Log

//Run run one Scene
func (s *Scene) Run() {
	var p Pool
	log.Scene = s.Name
	choiceModel(s, p)
}

func choiceModel(s *Scene, p Pool) {
	if s.RunConfig.Type == RunTypeTime {
		controlByTime(s, p)
	}
}

func controlByTime(s *Scene, p Pool) error {
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
		if err := runInterfaces(s, logIns); err != nil {
			return err
		}
		return nil
	})
	p.start()
	log.RealTimeData(&p)
	<-endFlag
	p.stop()
	return nil
}

func initLogIn(ins []Interface) map[string]*In {
	logIns := map[string]*In{}
	for _, in := range ins {
		logIn := &In{
			Name:     in.Name,
			Consumes: []int{},
		}
		logIns[in.Name] = logIn
	}
	return logIns
}

func runInterfaces(s *Scene, logIns map[string]*In) error {
	cache := NewCache()
	for _, in := range s.Interfaces {
		logIn, _ := logIns[in.Name]
		err := in.DataPrepare(cache)
		if err != nil {
			return err
		}
		res, err := in.Request()
		MuConsume.Lock()
		logIn.Consumes = append(logIn.Consumes, int(in.Consuming))
		MuConsume.Unlock()
		if err != nil {
			return err
		}
		if res.StatusCode == 200 {
			body, err := ioutil.ReadAll(res.Body)
			if err != nil {
				return err
			}
			for _, store := range in.Stores {
				err := store.Save(string(body), cache)
				if err != nil {
					return err
				}
			}
		}
	}
	return nil
}
