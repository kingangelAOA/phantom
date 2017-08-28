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
			time.Sleep(1 * time.Second)
			fmt.Printf("scene: %s; threadNum: %d  *************\n", l.Scene, l.ThreadNum)
			for _, in := range l.Interfaces {
				average, runNum, runTime := in.averageData()
				fmt.Println("run:", runTime, runNum, (float64(runTime) / 1000))
				if runNum > 0 && runTime > 0 {
					fmt.Printf("name: %s, QPS: %f, average: %f, 95 line: %f, 99 line: %f, highest: %f, err: %d\n", in.Name, float64(runNum)/(float64(runTime)/1000), average, in.lineData(0.95), in.lineData(0.99), in.highestConsume(), len(p.Result)/int(runNum))
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

func (i *In) highestConsume() float64 {
	MuConsume.Lock()
	defer MuConsume.Unlock()
	sort.Float64s(i.Consumes)
	length := len(i.Consumes)
	return i.Consumes[length-1]
}

func (i *In) averageData() (float64, int, float64) {
	MuConsume.Lock()
	defer MuConsume.Unlock()
	sum := float64(0)
	num := len(i.Consumes)
	for _, v := range i.Consumes {
		sum += v
	}
	// fmt.Println("consumes:", sum, i.Consumes)
	return sum / float64(num), num, sum
}
