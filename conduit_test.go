package conduit

import (
	"testing"
)

type TestConduitJob struct {
	n int
}

func (j *TestConduitJob) Do() []Job {
	if j.n == 1 {
		return []Job{&TestConduitJob{2}}
	}
	return nil
}

func (j *TestConduitJob) Route() int {
	return j.n
}

func TestConduit(t *testing.T) {

	cfg := Config{
		Stages: []Stage{
			Stage{
				Route: 1, Grow: 1, Shrink: 1, Size: 1,
			},
			Stage{
				Route: 2, Grow: 32, Shrink: 16, Size: 4,
			},
		},
	}

	net, err := NewNetwork(cfg)
	if err != nil {
		t.Fatal(err)
	}

	net.Start()

	// put 10 jobs into the network
	for i := 0; i < 10; i++ {
		net.Push([]Job{&TestConduitJob{1}})
	}

	net.Stop()
}
