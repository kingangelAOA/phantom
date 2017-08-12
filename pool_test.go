package phantom

import (
	"fmt"
	"testing"
	"time"
)

func Download(url string) error {
	time.Sleep(1 * time.Second)
	fmt.Println("Download " + url)
	return nil
}

func DownloadFinish() {
	fmt.Println("Download finsh")
}

func TestInit(t *testing.T) {
	var p Pool
	url := []string{"11111", "22222", "33333", "444444", "55555", "66666", "77777", "88888", "999999"}
	p.SetFinishCallback(DownloadFinish)
	p.Init(uint16(9), uint16(9))
	// if false {
	// 	t.Error("sdfasf")
	// }

	for i := range url {
		u := url[i]
		p.addTask(func() error {
			return Download(u)
		})
	}
	p.start()
	p.stop()
}
