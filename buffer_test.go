package conduit

import (
	"testing"
)

type TestBufferJob struct {
}

func (j *TestBufferJob) Do() []Job  { return nil }
func (j *TestBufferJob) Route() int { return 0 }

func TestBuffer(t *testing.T) {

	b := NewBuffer()
	b.output = make(chan Job)

	go b.main()

	for i := 0; i < 20; i++ {
		b.input <- &TestBufferJob{}
	}

	for i := 0; i < 10; i++ {
		<-b.output
	}

	close(b.cancel)

	<-b.done
}
