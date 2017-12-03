package conduit

import (
	"testing"
	"time"
)

const (
	portOne = iota
	portTwo
)

var NoOp NodeFunc = func(ctx *Context) {
}

var Passthrough NodeFunc = func(ctx *Context) {
	ctx.Out(portOne, ctx.In())
}

var AddOne NodeFunc = func(ctx *Context) {
	ctx.Out(portOne, ctx.In().(int)+1)
}

func Sleeper(d time.Duration) NodeFunc {
	return func(ctx *Context) {
		time.Sleep(d)
	}
}

func TestNoOp(t *testing.T) {
	net := NewNetwork()

	node := net.AddNode(NoOp)

	input := make(chan interface{})
	net.MapInput(node, input)

	net.Start()
	input <- 1
	close(input)
	net.Wait()
}

func TestOp(t *testing.T) {
	var val int = 1

	net := NewNetwork()

	addOne := net.AddNode(func(ctx *Context) {
		ctx.Out(portOne, ctx.In().(int)+1)
	})

	input := make(chan interface{})
	net.MapInput(addOne, input)

	output := net.MapOutput(addOne, portOne)

	net.Start()

	input <- val
	close(input)

	if <-output != val+1 {
		t.Fail()
	}

	net.Wait()
}

func TestConnections(t *testing.T) {
	var val int = 1

	net := NewNetwork()

	addOne1 := net.AddNode(AddOne)
	addOne2 := net.AddNode(AddOne)
	addOne3 := net.AddNode(AddOne)

	addOne1.Connect(addOne2, portOne)
	addOne2.Connect(addOne3, portOne)

	input := make(chan interface{})
	net.MapInput(addOne1, input)

	output := net.MapOutput(addOne3, portOne)

	net.Start()

	input <- val
	close(input)

	if <-output != val+3 {
		t.Fail()
	}

	net.Wait()
}

func TestConsolidated(t *testing.T) {
	var val int = 1

	net := NewNetwork()

	node1 := net.AddNode(Passthrough)
	node2 := net.AddNode(Passthrough)
	node3 := net.AddNode(Passthrough)

	input1 := make(chan interface{})
	net.MapInput(node1, input1)

	input2 := make(chan interface{})
	net.MapInput(node2, input2)

	// 1 -> 3 <- 2
	node1.Connect(node3, portOne)
	node2.Connect(node3, portOne)

	output := net.MapOutput(node3, portOne)

	net.Start()

	input1 <- val
	close(input1)

	input2 <- val
	close(input2)

	if <-output != val {
		t.Fail()
	}
	if <-output != val {
		t.Fail()
	}

	net.Wait()
}

func TestPool(t *testing.T) {
	var poolSize uint = 4

	net := NewNetwork()

	delta := 100 * time.Millisecond
	node := net.AddNode(Sleeper(delta))
	node.Pool(poolSize)

	input := make(chan interface{})
	net.MapInput(node, input)

	go func() {
		for i := uint(0); i < poolSize; i++ {
			input <- 1
		}
		close(input)
	}()

	start := time.Now()
	net.Start()
	net.Wait()
	elapsed := time.Now().Sub(start)

	if elapsed >= time.Duration(poolSize) * delta {
		t.Errorf("elapsed should be < %s, is %s", time.Duration(poolSize) * delta, elapsed)
	}

}

func TestPanics(t *testing.T) {

	net := NewNetwork()
	node1 := net.AddNode(NoOp)
	node2 := net.AddNode(NoOp)

	func() {
		defer func() {
			e := recover()
			if e == nil {
				t.Error("expected panic but recover returned nil")
			}else if e != ErrOutputConnected {
				t.Errorf("expected recover to return %s but returned %s", ErrOutputConnected, e)
			}
		}()

		node1.Connect(node2, portOne)
		node1.Connect(node2, portOne)
	}()

	func() {
		defer func() {
			e := recover()
			if e == nil {
				t.Error("expected panic but recover returned nil")
			}else if e != ErrPoolSizeZero {
				t.Errorf("expected recover to return %s but returned %s", ErrPoolSizeZero, e)
			}
		}()

		node1.Pool(0)
	}()

}
