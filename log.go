package phantom

import (
	"fmt"
	"sort"
	"sync"
	"time"
)

//MuConsume lock consumes
var MuConsume sync.Mutex

//Log performance Log
type Log struct {
	Scene      string
	Interfaces []*In
	ThreadNum  uint16
}

//In inteface log
type In struct {
	Name     string
	Consumes []int
}

//RealTimeData print real time log
func (l *Log) RealTimeData(p *Pool) {
	go func() {
		for {
			time.Sleep(5 * time.Second)
			fmt.Printf("scene: %s; threadNum: %d  *************\n", l.Scene, l.ThreadNum)
			for _, in := range l.Interfaces {
				fmt.Printf("name: %s, QPS: %d, average: %d, 95 line: %d, 99 line: %d, highest: %d\n", in.Name, len(in.Consumes)/int(p.RunTime), in.averageData(), in.lineData(0.95), in.lineData(0.99), in.highestConsume())
			}
		}
	}()
}

func (i *In) lineData(f float32) int {
	MuConsume.Lock()
	defer MuConsume.Unlock()
	sort.Ints(i.Consumes)
	index := int(float32(len(i.Consumes)) * f)
	return i.Consumes[index]
}

func (i *In) highestConsume() int {
	MuConsume.Lock()
	defer MuConsume.Unlock()
	sort.Ints(i.Consumes)
	length := len(i.Consumes)
	return i.Consumes[length-1]
}

func (i *In) averageData() int {
	MuConsume.Lock()
	defer MuConsume.Unlock()
	sum := 0
	for _, v := range i.Consumes {
		sum += v
	}
	return sum / len(i.Consumes)
}
