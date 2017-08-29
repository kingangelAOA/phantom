package phantom

import "time"

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
	runTime              int64
	runNum               int64
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
		for !p.finished {
			<-p.done
			if !p.finished {
				p.ConcurrentGoroutines <- struct{}{}
			}
		}
	}()

	go func() {
		nowTime := time.Now().UnixNano() / int64(time.Millisecond)
		runNum := int64(1)
		for !p.finished {
			<-p.ConcurrentGoroutines
			go func() {
				if err := p.task(); err != nil {
					p.Result <- err
				} else {
					MuRun.Lock()
					p.runTime = time.Now().UnixNano()/int64(time.Millisecond) - nowTime
					p.runNum = runNum
					MuRun.Unlock()
					runNum++
				}
				if !p.finished {
					p.done <- true
				}
			}()
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
