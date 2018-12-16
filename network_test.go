package conduit

import (
	"testing"
)

type TestNetworkJob struct {
	n int
}

func (j *TestNetworkJob) Do() []Job {
	if j.n > 0 {
		j.n--
		return []Job{&TestNetworkJob{j.n}}
	}
	return nil
}

func (j *TestNetworkJob) Route() int {
	return j.n
}

func TestNetwork(t *testing.T) {

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
		t.Fail()
	}
	go net.main()

	for i := 0; i < 20; i++ {
		net.input <- &TestNetworkJob{2}
	}

	close(net.cancel)

	<-net.done
}
