package main


import (
	"fmt"
	"os"
	"github.com/kb1vc/radiopath/nedmap"
)


func main() {
	for _, infile := range os.Args[1:] {
		mi, err := nedmap.GetInfo(infile)
		if err != nil {
			panic(err)
		}
		fmt.Println(mi)
	}
}

