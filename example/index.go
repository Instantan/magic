package main

import (
	"time"

	"github.com/Instantan/magic"
)

type IndexPageData struct {
	Name   magic.Get[string] `json:"name"`
	Nested struct {
		Count magic.Get[int] `json:"count"`
	} `json:"nested"`
}

func IndexPage(ctx magic.PageContext) any {
	getName, _ := magic.CreateSignal("Felix")
	getCount, setCount := magic.CreateSignal(1)

	go func() {
		for {
			time.Sleep(time.Second)
			setCount(getCount() + 1)
		}
	}()

	return IndexPageData{
		Name: getName,
		Nested: struct {
			Count magic.Get[int] `json:"count"`
		}{
			Count: getCount,
		},
	}
}
