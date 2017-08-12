package phantom

import (
	"fmt"
	"time"
)

//Pool pool object
type Pool struct {
	GoroutineNum         uint16
	ConcurrentGoroutines chan struct{}
	Total                uint16
	Result               chan error
	FinishCallback       func()
	done                 chan bool
	waitForAllJobs       chan bool
	task                 func() error
	finished             bool
	RunNum               uint64
	RunTime              int64
}

//Init  init pool
func (p *Pool) Init(GoroutineNum uint16, total uint16) {
	p.Total = total
	p.GoroutineNum = GoroutineNum
	p.ConcurrentGoroutines = make(chan struct{}, p.GoroutineNum)
	p.Result = make(chan error, total)
	p.done = make(chan bool)
	p.waitForAllJobs = make(chan bool)
	p.initConcurrentGoroutines()
	p.finished = false
}

//AddTask add task
func (p *Pool) addTask(task func() error) {
	p.task = task
}

func (p *Pool) initConcurrentGoroutines() {
	for i := 0; i < int(p.GoroutineNum); i++ {
		p.ConcurrentGoroutines <- struct{}{}
	}
}

//Start start pool
func (p *Pool) start() {
	go func() {
		for {
			<-p.done
			if !p.finished {
				p.ConcurrentGoroutines <- struct{}{}
			}
		}
	}()

	go func() {
		beginTime := time.Now().Unix()
		index := uint64(1)
		for !p.finished {
			<-p.ConcurrentGoroutines
			go func() {
				if err := p.task(); err != nil {
					fmt.Println(err)
				}
				if !p.finished {
					p.done <- true
				}
				p.RunNum = index
				p.RunTime = time.Now().Unix() - beginTime
			}()
			index++
		}
	}()
}

//Stop step pool
func (p *Pool) stop() {
	p.finished = true
	close(p.ConcurrentGoroutines)
	close(p.done)
}

//SetFinishCallback callback when finish
func (p *Pool) SetFinishCallback(fun func()) {
	p.FinishCallback = fun
}
