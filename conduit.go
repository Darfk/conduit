package conduit

import (
	"fmt"
	"sync"
)

const (
	PoolSize = iota + 1
	GrowBy
	ShrinkBy
)

type stage struct {
	inc      chan Job
	cancel   chan struct{}
	poolSize int
	growBy   int
	shrinkBy int
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

	stage.poolSize = 1
	stage.growBy = 16
	stage.shrinkBy = 16

	for i, _ := range options {
		if options[i].key == PoolSize {
			if val, ok := options[i].val.(int); ok && val >= 1 {
				stage.poolSize = val
			} else {
				panic(fmt.Errorf("PoolSize option expects a positive integer, got (%T)%q",
					options[i].val, options[i].val))
			}
		} else if options[i].key == GrowBy {
			if val, ok := options[i].val.(int); ok && val >= 1 {
				stage.growBy = val
			} else {
				panic(fmt.Errorf("GrowBy option expects a positive integer, got (%T)%q",
					options[i].val, options[i].val))
			}
		} else if options[i].key == ShrinkBy {
			if val, ok := options[i].val.(int); ok && val >= 1 {
				stage.shrinkBy = val
			} else {
				panic(fmt.Errorf("ShrinkBy option expects a positive integer, got (%T)%q",
					options[i].val, options[i].val))
			}
		} else {
			panic(fmt.Errorf("option (%T)%q:(%T)%q is not a suitable option for a stage",
				options[i].key, options[i].key, options[i].val, options[i].val))
		}
	}

	stage.inc = make(chan Job, 1)

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
			grw int = stage.growBy
			// distance the output pointer gets away from the start before resize
			shk int = stage.shrinkBy

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
				if opt == shk {
					buf = buf[shk:]
					opt -= shk
					ipt -= shk
				}
				if opn {
					if ipt == len(buf) {
						nbf := make([]Job, len(buf)+grw)
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
		close(stage.inc)
	}
	net.sw.Wait()
}
