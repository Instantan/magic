package magic

import (
	"log"

	"github.com/panjf2000/ants"
)

var routinePool, _ = ants.NewPool(1000, ants.WithNonblocking(true), ants.WithPanicHandler(func(v interface{}) {
	log.Printf("Cought panic: %v", v)
}))

func submitTask(task func()) {
	routinePool.Submit(task)
}
