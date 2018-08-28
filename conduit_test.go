package conduit

import (
	"log"
	"testing"
)

const (
	Factorial = iota
	Print
)

type BadPortJob struct{}

func (job *BadPortJob) Port() int { return 9999 }
func (job *BadPortJob) Do() []Job { return nil }

type PrintJob struct {
	a int
}

func (job *PrintJob) Port() int { return Print }
func (job *PrintJob) Do() []Job {
	log.Println(job.a)
	return nil
}

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

	func() {
		defer func() {
			if r := recover(); r == nil {
				t.Log(r)
				t.Errorf("setting the InputBuffer option to a negative number did not cause a panic")
			}
		}()

		net.AddStage(1, Option(InputBuffer, -1))
	}()

	func() {
		defer func() {
			if r := recover(); r == nil {
				t.Errorf("setting the PoolSize option to a negative number did not cause a panic")
			}
		}()
		net.AddStage(2, Option(PoolSize, 0))
	}()
}
