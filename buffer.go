package conduit

type Buffer struct {
	input  chan Job
	output chan Job
	grow   int
	shrink int
	cancel chan struct{}
	done   chan struct{}
}

func vis(buf []Job, ipt, opt int, opn bool) string {
	s := ""
	if opn {
		s = "O"
	} else {
		s = "X"
	}

	t := "["

	for i := 0; i < len(buf); i++ {
		if ipt > i && opt <= i {
			t += "-"
		} else {
			t += "_"
		}

		if ipt == i && opt == i {
			s += "*"
		} else if ipt == i {
			s += "i"
		} else if opt == i {
			s += "o"
		} else {
			s += " "
		}
	}
	t += "]"

	return s + "\n" + t
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
		open   bool     = true
		ipt    int      = 0
		opt    int      = 0
		buf    []Job    = make([]Job, b.grow)
		output chan Job = nil
	)

	for open {
		select {
		case <-b.cancel:
			open = false
		case job := <-b.input:
			if ipt == len(buf) {
				nbf := make([]Job, len(buf)+b.grow)
				copy(nbf, buf)
				buf = nbf
			}
			buf[ipt] = job
			ipt++
			output = b.output
		case output <- buf[opt]:
			opt++
			if opt == b.shrink {
				buf = buf[b.shrink:]
				opt -= b.shrink
				ipt -= b.shrink
			}
			if opt == ipt {
				output = nil
			}
		}
	}
	close(b.done)
}
