package phantom

import (
	"fmt"
	"sort"
	"sync"
	"time"
)

//MuConsume lock consumes
var MuConsume sync.Mutex

var MuRun sync.Mutex

//Log performance Log
type Log struct {
	Scene      string
	Interfaces []*In
	ThreadNum  uint16
}

//In inteface log
type In struct {
	Name     string
	Consumes []float64
}

//RealTimeData print real time log
func (l *Log) RealTimeData(p *Pool) {
	go func() {
		for {
			time.Sleep(2 * time.Second)
			fmt.Printf("scene: %s; threadNum: %d  *************\n", l.Scene, l.ThreadNum)
			for _, in := range l.Interfaces {
				average := in.averageData()
				min, max := in.minMaxConsume()
				runNum := p.runNum
				runTime := p.runTime
				if runNum > 0 && runTime > 0 && (p.runTime/1000) > 0 {
					fmt.Printf("name: %s, QPS: %d, average: %f, 95 line: %f, 99 line: %f, min: %f, max: %f, err: %d\n", in.Name, runNum/(p.runTime/1000), average, in.lineData(0.95), in.lineData(0.99), min, max, len(p.Result)/int(runNum))
					in.Consumes = in.Consumes[:0]
				}
			}
		}
	}()
}

func (i *In) lineData(f float32) float64 {
	MuConsume.Lock()
	defer MuConsume.Unlock()
	sort.Float64s(i.Consumes)
	index := int(float32(len(i.Consumes)) * f)
	return i.Consumes[index]
}

func (i *In) minMaxConsume() (float64, float64) {
	MuConsume.Lock()
	defer MuConsume.Unlock()
	sort.Float64s(i.Consumes)
	length := len(i.Consumes)
	return i.Consumes[0], i.Consumes[length-1]
}

func (i *In) averageData() float64 {
	MuConsume.Lock()
	defer MuConsume.Unlock()
	sum := float64(0)
	num := int64(len(i.Consumes))
	for _, v := range i.Consumes {
		sum += v
	}
	return sum / float64(num)
}
