package main

import (
	"fmt"

	"github.com/mmlt/testr"
)

// E is an custom error.
type E struct {
	str string
}

func (e E) Error() string {
	return e.str
}

// T emulates GO's test package Log() method.
type T struct{}

func (t T) Log(args ...interface{}) {
	fmt.Println(args)
}

func main() {
	testr.SetVerbosity(1)
	log := testr.New(T{})
	log = log.WithName("MyName").WithValues("user", "you")
	log.Info("hello", "val1", 1, "val2", map[string]int{"k": 1})
	log.V(1).Info("you should see this")
	log.V(3).Info("you should NOT see this")
	log.Error(nil, "uh oh", "trouble", true, "reasons", []float64{0.1, 0.11, 3.14})
	log.Error(E{"an error occurred"}, "goodbye", "code", -1)
}
