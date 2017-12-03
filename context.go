package conduit

import ()

type Context struct {
	node  *Node
	input interface{}
}

func (q *Context) Out(port int, data interface{}) {
	q.node.outputs[port] <- data
}

func (q *Context) In() interface{} {
	return q.input
}
