package conduit

import (
	"testing"
)

const (
	BadPort = 9999
)

type BadPortJob struct{}

func (job *BadPortJob) Port() int { return BadPort }
func (job *BadPortJob) Do() []Job { return nil }

func TestUnroutable(t *testing.T) {
	func() {
		defer func() {
			if r := recover(); r == nil {
				t.Errorf("network did not panic when a job could not be routed")
			}
		}()

		// create a new network which will panic for unroutable jobs
		net := NewNetwork()

		// add a stage which listens to port 0
		net.AddStage(0)

		// start the network
		net.Start()

		// add a job that should fail to route
		net.AddJobs(&BadPortJob{})

		// don't wait for the network to finish because it never will
		net.Wait()
	}()
}

func TestPanics(t *testing.T) {
	net := NewNetwork()

	func() {
		defer func() {
			if r := recover(); r == nil {
				t.Errorf("adding a duplicate stage did not cause a panic")
			}
		}()

		net.AddStage(0)
		net.AddStage(0)

	}()

	{
		tests := []struct {
			opt option
			msg string
		}{
			{Option(PoolSize, 0),
				"setting the PoolSize option to a number < 1 did not cause a panic"},
			{Option(GrowBy, 0),
				"setting the GrowBy option to a number < 1 did not cause a panic"},
			{Option(ShrinkBy, 0),
				"setting the ShrinkBy option to a number < 1 did not cause a panic"},
		}
		for _, test := range tests {
			func() {
				defer func() {
					if r := recover(); r == nil {
						t.Errorf(test.msg)
					}
				}()
				net.AddStage(2, test.opt)
			}()
		}
	}

	func() {
		defer func() {
			if r := recover(); r == nil {
				t.Errorf("using an invalid stage option did not cause a panic")
			}
		}()

		invalidKey := 0

		net.AddStage(0, Option(invalidKey, 0))
	}()

}

// the point of this job is to duplicate itself and
// die after 4 iterations
type CoolJob struct {
	n int
}

func (job *CoolJob) Port() int { return 0 }
func (job *CoolJob) Do() []Job {
	job.n++
	if job.n > 4 {
		return nil
	}
	return []Job{
		&CoolJob{job.n},
		&CoolJob{job.n},
	}
}

func TestOperation(t *testing.T) {
	net := NewNetwork()
	net.AddStage(0)
	net.AddJobs(&CoolJob{})
	net.Start()
	net.Wait()
}

func TestOperationNetworkStartFlipped(t *testing.T) {
	net := NewNetwork()
	net.AddStage(0)
	net.Start()
	net.AddJobs(&CoolJob{})
	net.Wait()
}

func TestInfiniteChan(t *testing.T) {
	net := NewNetwork()
	net.AddStage(0)
	net.Start()
	net.AddJobs(&CoolJob{})
	net.AddJobs(&CoolJob{})
	net.AddJobs(&CoolJob{})
	net.AddJobs(&CoolJob{})
	net.AddJobs(&CoolJob{})
	net.AddJobs(&CoolJob{})
	net.Wait()
}
