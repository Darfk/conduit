package conduit

import (
	"sync"
)

type Network struct {
	nodes []*Node
	wg    sync.WaitGroup
}

func NewNetwork() *Network {
	net := new(Network)
	return net
}

func (q *Network) AddNode(fn NodeFunc) (n *Node) {
	n = newNode(fn)
	q.nodes = append(q.nodes, n)
	return
}

func (q *Network) Start() {
	for _, node := range q.nodes {

		node := node

		ctx := &Context{
			node: node,
		}

		if len(node.consolidate) > 0 {
			for _, c := range node.consolidate {
				node.wg.Add(1)
				go func(c chan interface{}) {
					for v := range c {
						node.input <- v
					}
					node.wg.Done()
				}(c)
			}

			go func() {
				node.wg.Wait()
				close(node.input)
			}()
		}

		q.wg.Add(int(node.pool))

		for i:= uint(0);i<node.pool;i++ {
			go func() {
				for ctx.input = range node.input {
					node.state |= StateBusy
					node.fn(ctx)
					node.state &= ^StateBusy
				}

				node.state |= StateTerminated
				node.shutdown()
				q.wg.Done()

			}()
		}

	}
}

func (q *Network) MapInput(n *Node, c chan interface{}) {
	n.addInput(c)
}

func (q *Network) MapOutput(n *Node, port int) chan interface{} {
	return n.addOutput(port)
}

func (q *Network) Wait() {
	q.wg.Wait()
}
