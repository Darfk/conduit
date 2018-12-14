package conduit

import ()

type Pool struct {
	input  chan Job
	output chan Job
	size   int
	cancel chan struct{}
	done   chan struct{}
}

func NewPool(size int) (p *Pool) {
	p = &Pool{
		size:   size,
		input:  make(chan Job),
		done:   make(chan struct{}),
		cancel: make(chan struct{}),
	}
	return
}

func (p *Pool) main() {

	var done chan struct{} = make(chan struct{})
	var cancel chan struct{} = make(chan struct{})

	for i := 0; i < p.size; i++ {
		i := i
		go func() {
			var (
				open bool = true
				opt  int  = 0
				jobs []Job
			)
			for open {
				if opt == len(jobs) {
					select {
					case <-cancel:
						open = false
					case job := <-p.input:
						println("pool", i, "working on", job)
						jobs = job.Do()
						opt = 0
					}
				} else {
					select {
					case <-cancel:
						open = false
					case p.output <- jobs[opt]:
						opt++
					}
				}
			}
			done <- struct{}{}
		}()
	}

	<-p.cancel

	close(cancel)

	for i := 0; i < p.size; i++ {
		<-done
	}

	close(p.done)
}
