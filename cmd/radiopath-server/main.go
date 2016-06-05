// server that accepts requests with from/to location
// and produces a "slice of the earth" radio path model
package main

import (
	"github.com/kb1vc/radiopath/location"
	"fmt"
)

func main() {
	fmt.Println("Hello World")
	ll,err := location.FromGrid("FN42bl")
	if err == nil {
		fmt.Println(ll)
	} else {
		fmt.Println(err)
	}

	fmt.Println(location.ToGrid(ll))

	wach, _ := location.FromGrid("FM19ng")
	wash, _ := location.FromGrid("FN44ig")
	fmt.Println(wash)

	fmt.Println(wach.Bearing(wash))
	fmt.Println(wash.Bearing(wach))

	brg, _, dist := wach.Bearing(wash)
	
	fmt.Println(wach.Onpath(brg, dist))
}
