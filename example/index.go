package main

import (
	"time"

	"github.com/Instantan/magic"
)

type IndexPageData struct {
	Name    magic.Get[string] `json:"name"`
	Counter struct {
		Count magic.Get[int]  `json:"count"`
		Even  magic.Get[bool] `json:"even"`
	} `json:"counter"`
}

func IndexPage(ctx magic.PageContext) any {
	getName, _ := magic.Signal("Felix")
	getCount, setCount := magic.Signal(1)

	go func() {
		for {
			time.Sleep(time.Second)
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
		Name: getName,
		Counter: struct {
			Count magic.Get[int]  `json:"count"`
			Even  magic.Get[bool] `json:"even"`
		}{
			Count: getCount,
			Even:  even,
		},
	}
}
