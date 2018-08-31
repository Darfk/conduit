package conduit

import (
	"fmt"
	"sync"
)

const (
	InputBuffer = iota
	PoolSize
)

type stage struct {
	inc         chan Job
	cancel      chan struct{}
	poolSize    int
	inputBuffer int
}

type Job interface {
	Do() []Job
	Port() int
}

type Network struct {
	stages map[int]*stage
	wg     sync.WaitGroup
}

func NewNetwork(options ...option) (net *Network) {
	n := &Network{
		stages: make(map[int]*stage),
	}

	net = n

	return
}

type option struct {
	key int
	val interface{}
}

func Option(key int, val interface{}) option {
	return option{
		key, val,
	}
}

func (net *Network) AddStage(port int, options ...option) {
	if _, exists := net.stages[port]; exists {
		panic(fmt.Errorf("cannot create a stage on port %d, a stage already exists", port))
		return
	}

	stage := &stage{}

	stage.inputBuffer = 0
	stage.poolSize = 1

	for i, _ := range options {
		if options[i].key == InputBuffer {
			if val, ok := options[i].val.(int); ok && val >= 0 {
				stage.inputBuffer = val
			} else {
				panic(fmt.Errorf("InputBuffer option expects a positive integer or zero, got (%T)%q", options[i].val, options[i].val))
				return
			}
		}

		if options[i].key == PoolSize {
			if val, ok := options[i].val.(int); ok && val >= 1 {
				stage.poolSize = val
			} else {
				panic(fmt.Errorf("PoolSize option expects a positive integer, got (%T)%q", options[i].val, options[i].val))
				return
			}
		}
	}

	stage.inc = make(chan Job, stage.inputBuffer)

	net.stages[port] = stage
}

func (net *Network) route(jobs []Job) {
	net.wg.Add(len(jobs))
	for _, job := range jobs {
		port := job.Port()
		if stage, exists := net.stages[port]; exists {
			stage.inc <- job
			continue
		}
		panic(fmt.Errorf("could not route job %q: no stage exists at port %d", job, port))
	}
}

func (net *Network) AddJobs(jobs ...Job) {
	net.route(jobs)
}

func (net *Network) Start() {
	for i, _ := range net.stages {
		stage := net.stages[i]
		pool := sync.WaitGroup{}
		pool.Add(stage.poolSize)
		for i := 0; i < stage.poolSize; i++ {
			go func() {
				defer pool.Done()
				for {

					in, open := <-stage.inc

					if !open {
						return
					}
					jobs := in.Do()
					net.route(jobs)
					net.wg.Done()
				}
			}()
		}
	}
}

func (net *Network) Wait() {
	net.wg.Wait()
}
