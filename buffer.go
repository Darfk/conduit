package conduit

type Buffer struct {
	input  chan Job
	output chan Job
	grow   int
	shrink int
	cancel chan struct{}
	done   chan struct{}
	buf    []Job
}

func NewBuffer(grow, shrink int) (b *Buffer) {
	b = &Buffer{
		input:  make(chan Job),
		grow:   grow,
		shrink: shrink,
		cancel: make(chan struct{}),
		done:   make(chan struct{}),
		buf:    make([]Job, grow),
	}
	return
}

func (b *Buffer) main() {
	var (
		open bool = true
		ipt  int  = 0
		opt  int  = 0
	)

	for open {
		if opt == ipt {
			select {
			case <-b.cancel:
				open = false
			case job := <-b.input:
				b.buf[ipt] = job
				ipt++
			}
		} else {
			select {
			case <-b.cancel:
				open = false
			case job := <-b.input:
				b.buf[ipt] = job
				ipt++
			case b.output <- b.buf[opt]:
				opt++
			}
		}

		if ipt == len(b.buf) {
			nbf := make([]Job, len(b.buf)+b.grow)
			copy(nbf, b.buf)
			b.buf = nbf
		}

		if opt == b.shrink {
			b.buf = b.buf[b.shrink:]
			opt -= b.shrink
			ipt -= b.shrink
		}

	}
	close(b.done)
}
