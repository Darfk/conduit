//go:generate conduit
package main

import (
	"fmt"
	"log"
	"math/rand"
	"time"
)

func init() {
	rand.Seed(1)
}

type divop struct {
	n, d int
	err  error
	ans  int
}

// conduit pool
func divide(op *divop) *divop {
	time.Sleep(100 * time.Millisecond)
	if op.d == 0 {
		op.err = fmt.Errorf("divide by zero")
	} else {
		op.ans = op.n / op.d
	}
	return op
}

// conduit
func op() *divop {
	return &divop{n: rand.Intn(15), d: rand.Intn(4)}
}

// conduit
func printResult(in *divop) {
	log.Printf("%-d / %-d = %-d %v", in.n, in.d, in.ans, in.err)
}

func main() {
	cancel := make(chan struct{})
	defer close(cancel)

	ops := opSource(cancel)
	res := divideStagePool(ops, cancel, 4)
	printResultSink(res, cancel)

	time.Sleep(2 * time.Second)
}
