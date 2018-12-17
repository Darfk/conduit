package main

import (
	conduit ".."
	"time"
)

const (
	RouteGreet = iota
	RoutePrint
)

type GreetJob struct {
	name string
}

func (j *GreetJob) Route() int { return RouteGreet }
func (j *GreetJob) Do() []conduit.Job {
	if j.name == "Andy" || j.name == "Nathan" {
		// say hi to our friends
		return []conduit.Job{&PrintJob{"Hey " + j.name + "!"}}
	} else if j.name == "Kirk" {
		// give Kirk the cold shoulder
		return nil
	}
	// we haven't met them before
	return []conduit.Job{
		&PrintJob{"Hi " + j.name},
		&PrintJob{"Welcome to conduit " + j.name + "."},
	}
}

type PrintJob struct {
	text string
}

func (j *PrintJob) Route() int { return RoutePrint }
func (j *PrintJob) Do() []conduit.Job {
	// printing takes a while, we will use a worker pool
	time.Sleep(time.Millisecond * 250)
	println(j.text)
	return nil
}

func main() {
	net, err := conduit.NewNetwork(conduit.Config{
		Stages: []conduit.Stage{
			// set up a stage to operate on greetings only, we should only need
			// a worker pool of size 1 because greetings happen very fast
			conduit.Stage{Route: RouteGreet, Grow: 8, Shrink: 8, Size: 1},

			// since printing takes a while we will allocate 8 workers
			conduit.Stage{Route: RoutePrint, Grow: 8, Shrink: 8, Size: 8},
		},
	})

	if err != nil {
		println(err.Error())
		return
	}

	// start the network so that it can accept work
	net.Start()

	// push all the jobs into the network
	net.Push([]conduit.Job{
		&PrintJob{"This is conduit"},
		&GreetJob{"Andy"},
		&GreetJob{"Nathan"},
		&GreetJob{"Kirk"},
		&GreetJob{"Rebecca"},
		&GreetJob{"Miles"},
		&GreetJob{"Emily"},
		&GreetJob{"Steven"},
	})

	time.Sleep(time.Second)

	// stop and wait for the network to completely drain and close
	// this will not complete the jobs in progress
	net.Stop()
}
