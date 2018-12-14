package conduit

import (
	"testing"
	"time"
)

type TestPoolJob struct {
	n int
}

func (j *TestPoolJob) Do() []Job {
	time.Sleep(time.Second)
	return []Job{j}
}
func (j *TestPoolJob) Route() int { return 0 }

func TestPool(t *testing.T) {
	p := NewPool(10)

	p.output = make(chan Job)
	go p.main()

	go func() {
		for i := 0; i < 20; i++ {
			p.input <- &TestPoolJob{1}
		}
	}()

	for i := 0; i < 10; i++ {
		<-p.output
	}

	close(p.cancel)
	<-p.done
}
