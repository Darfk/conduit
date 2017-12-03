package conduit

import (
	"sync"
)

const (
	StateTerminated = 1 << iota
	StateStarted
	StateBusy
)

type NodeFunc func(*Context)

type Node struct {
	fn          NodeFunc
	state       int
	consolidate []chan interface{}
	input       chan interface{}
	outputs     map[int]chan interface{}
	wg          sync.WaitGroup
	pool        uint
}

func newNode(fn NodeFunc) (n *Node) {
	n = &Node{
		fn:      fn,
		outputs: make(map[int]chan interface{}),
		pool:    1,
	}
	return
}

func (q *Node) addOutput(port int) chan interface{} {
	_, e := q.outputs[port]
	if e {
		panic(ErrOutputConnected)
	}

	q.outputs[port] = make(chan interface{})

	return q.outputs[port]
}

func (q *Node) addInput(c chan interface{}) {
	if q.input == nil {
		q.input = c
	} else {
		if len(q.consolidate) < 1 {
			q.consolidate = append(q.consolidate, q.input)
			q.input = make(chan interface{})
		}
		q.consolidate = append(q.consolidate, c)
	}
}

func (q *Node) shutdown() {
	for _, output := range q.outputs {
		close(output)
	}
}

func (q *Node) Connect(other *Node, port int) *Node {
	if q.state&StateStarted == StateStarted {
		panic(ErrNetworkStartedConnect)
	}
	other.addInput(q.addOutput(port))
	return q
}

func (q *Node) Pool(size uint) *Node {
	if size < 1 {
		panic(ErrPoolSizeZero)
	}
	if q.state&StateStarted == StateStarted {
		panic(ErrNetworkStartedPool)
	}
	q.pool = size
	return q
}
