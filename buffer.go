package conduit

type Buffer struct {
	input  chan Job
	output chan Job
	grow   int
	shrink int
	cancel chan struct{}
	done   chan struct{}
}

func NewBuffer(grow, shrink int) (b *Buffer) {
	b = &Buffer{
		input:  make(chan Job),
		grow:   grow,
		shrink: shrink,
		cancel: make(chan struct{}),
		done:   make(chan struct{}),
	}
	return
}

func (b *Buffer) main() {
	var (
		open bool  = true
		ipt  int   = 0
		opt  int   = 0
		buf  []Job = make([]Job, b.grow)
	)

	for open {
		if opt == ipt {
			select {
			case <-b.cancel:
				open = false
			case job := <-b.input:
				buf[ipt] = job
				ipt++
			}
		} else {
			select {
			case <-b.cancel:
				open = false
			case job := <-b.input:
				buf[ipt] = job
				ipt++
			case b.output <- buf[opt]:
				opt++
			}
		}

		if ipt == len(buf) {
			nbf := make([]Job, len(buf)+b.grow)
			copy(nbf, buf)
			buf = nbf
		}

		if opt == b.shrink {
			buf = buf[b.shrink:]
			opt -= b.shrink
			ipt -= b.shrink
		}

	}
	close(b.done)
}
