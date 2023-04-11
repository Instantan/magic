package main

import (
	"github.com/Instantan/magic"
)

type IndexPageData struct {
	Name       magic.ReactiveValue[string]   `json:"name"`
	HelloWorld magic.ReactiveValue[[]string] `json:"helloWorld"`
	Counter    struct {
		Count magic.ReactiveValue[int]  `json:"count"`
		Even  magic.ReactiveValue[bool] `json:"even"`
	} `json:"counter"`
}

func IndexPage(ctx magic.PageContext) any {
	name := magic.Value("Felix")
	count := magic.Value(1)
	helloWorld := magic.Value([]string{"H", "e", "l", "l", "o"})

	return IndexPageData{
		Name:       name,
		HelloWorld: helloWorld,
		Counter: struct {
			Count magic.ReactiveValue[int]  `json:"count"`
			Even  magic.ReactiveValue[bool] `json:"even"`
		}{
			Count: count,
		},
	}

}
