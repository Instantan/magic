package magic

import (
	"github.com/panjf2000/ants"
)

var routinePool, _ = ants.NewPool(10000, ants.WithNonblocking(true))

func submitTask(task func()) {
	routinePool.Submit(task)
}
