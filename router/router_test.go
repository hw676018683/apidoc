package router

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/lovego/goa"
)

func ExampleRouterTree() {
	router := goa.New()
	r := New(router.Group(`/root`), `root`)
	g1 := r.Group(`/group1`)
	g2 := r.Group(`/group2`)
	g1.GetX(`/users`, func(c *goa.Context) {})
	g2.Post(`/book`, func(c *goa.Context) {})

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