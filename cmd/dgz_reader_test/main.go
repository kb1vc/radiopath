package main

import (
	"fmt"
	"os"
	"github.com/kb1vc/radiopath/nedmap"
	"net/http"
	"log"
	"time"
)

import _ "net/http/pprof"


func doWork(i int) (int) {
	mcc,_ := nedmap.ReadZCompressedMap(os.Args[1])
	i++
	if i > 3600 { i = 1 }
	fmt.Println(mcc.Elevation[i][i])
	return i
}

func main() {
	go func() {
		log.Println(http.ListenAndServe("localhost:6060", nil))
	}()

	i := 1
	for {
		zreadStart := time.Now()		
		i = doWork(i)
		zreadElapsed := time.Since(zreadStart)
		fmt.Printf("Compressed readtime: %s\n", zreadElapsed)		
	}
}
