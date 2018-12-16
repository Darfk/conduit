package conduit

import ()

type stage struct {
	route  int
	pool   *Pool
	buffer *Buffer
}

type Network struct {
	input  chan Job
	output chan Job
	cancel chan struct{}
	done   chan struct{}
	stages map[int]*stage
}

func NewNetwork(cfg Config) (net *Network, err error) {
	if err = CheckConfig(cfg); err != nil {
		return
	}

	net = &Network{
		input:  make(chan Job),
		done:   make(chan struct{}),
		cancel: make(chan struct{}),
		stages: make(map[int]*stage),
	}

	for _, stageConfig := range cfg.Stages {
		buffer := NewBuffer(stageConfig.Grow, stageConfig.Shrink)
		pool := NewPool(stageConfig.Size)

		buffer.output = pool.input
		pool.output = net.input

		net.stages[stageConfig.Route] = &stage{
			route:  stageConfig.Route,
			pool:   pool,
			buffer: buffer,
		}
	}

	return
}

func (net *Network) main() {

	for _, stage := range net.stages {
		go stage.buffer.main()
		go stage.pool.main()
	}

	var (
		open bool = true
		job  Job  = nil
	)

	for open {
		if job == nil {
			select {
			case <-net.cancel:
				open = false
			case job = <-net.input:
			}
		} else {
			route := job.Route()
			select {
			case <-net.cancel:
				open = false
			case net.stages[route].buffer.input <- job:
				job = nil
			}
		}
	}

	for _, stage := range net.stages {
		close(stage.buffer.cancel)
		close(stage.pool.cancel)
	}

	for _, stage := range net.stages {
		<-stage.buffer.done
		<-stage.pool.done
	}

	close(net.done)

}
