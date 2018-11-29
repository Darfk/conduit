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
	sw     sync.WaitGroup
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
			}
		}

		if options[i].key == PoolSize {
			if val, ok := options[i].val.(int); ok && val >= 1 {
				stage.poolSize = val
			} else {
				panic(fmt.Errorf("PoolSize option expects a positive integer, got (%T)%q", options[i].val, options[i].val))
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

		var (
			// size to increase when the input pointer runs off the edge
			// TODO: make these configurable
			siz int = 10
			// distance the output pointer gets away from the start before resize
			drg int = 10

			// is the channel open?
			opn bool = true
			// output channel
			ouc = make(chan Job, 1)
			// output pointer
			opt int
			// input pointer
			ipt int
			// buffer
			buf []Job
		)

		go func() {
			for {
				if opt == drg {
					buf = buf[drg:]
					opt -= drg
					ipt -= drg
				}
				if opn {
					if ipt == len(buf) {
						nbf := make([]Job, len(buf)+siz)
						copy(nbf, buf)
						buf = nbf
					}
					if ipt == opt {
						select {
						case buf[ipt], opn = <-stage.inc:
							ipt++
						}
					} else if ipt > opt {
						select {
						case ouc <- buf[opt]:
							opt++
						case buf[ipt], opn = <-stage.inc:
							ipt++
						}
					}
				} else {
					if opt == ipt-1 {
						break
					}
					select {
					case ouc <- buf[opt]:
						opt++
					}
				}
			}
			close(ouc)
		}()

		for i := 0; i < stage.poolSize; i++ {
			net.sw.Add(1)
			go func() {
				for in := range ouc {
					// println(in)
					jobs := in.Do()
					net.route(jobs)
					net.wg.Done()
				}
				net.sw.Done()
			}()
		}

	}
}

func (net *Network) Wait() {
	// wait for all jobs to be out of the network
	net.wg.Wait()
	for _, stage := range net.stages {
		// println("close stage inc")
		close(stage.inc)
	}
	net.sw.Wait()
}
