package main

import (
	"time"

	"github.com/Instantan/magic"
)

type IndexPageData struct {
	Name       magic.Get[string]   `json:"name"`
	HelloWorld magic.Get[[]string] `json:"helloWorld"`
	Counter    struct {
		Count magic.Get[int]  `json:"count"`
		Even  magic.Get[bool] `json:"even"`
	} `json:"counter"`
}

func IndexPage(ctx magic.PageContext) any {
	getName, _ := magic.Signal("Felix")
	getCount, setCount := magic.Signal(1)
	getHelloWorld, _ := magic.Signal([]string{"H", "e", "l", "l", "o"})

	ctx.Cleanup(func() {

	})

	go func() {
		for {
			time.Sleep(time.Second * 1)
			setCount(getCount() + 1)
		}
	}()

	even := magic.Computed(func() bool {
		return getCount()%2 == 0
	}, getCount)

	magic.Effect(func() {
		if even() {
			println("is even")
		} else {
			println("is odd")
		}
	}, even)

	return IndexPageData{
		Name:       getName,
		HelloWorld: getHelloWorld,
		Counter: struct {
			Count magic.Get[int]  `json:"count"`
			Even  magic.Get[bool] `json:"even"`
		}{
			Count: getCount,
			Even:  even,
		},
	}

}
