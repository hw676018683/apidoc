package router

import (
	"encoding/json"
	"fmt"
	"log"
)

func Example_RouterTree() {
	r := New(``)
	g1 := r.Group(`/group1`).Title(`分组1`)
	g1.Get(`/users`).Doc(`用户`, ``, ``, nil, nil)
	g1.Group(`/child`)

	printJson(r)

	// Output:
}

func printJson(v interface{}) {
	data, err := json.MarshalIndent(v, ``, `  `)
	if err != nil {
		log.Panic(err)
	}
	fmt.Println(string(data))
}
