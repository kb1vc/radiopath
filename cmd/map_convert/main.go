package main


import (
	"fmt"
	"time"
	"os"
	"github.com/kb1vc/radiopath/nedmap"
)



func main() {
	rawE, cmpfname, cerr := nedmap.ConvertFile(os.Args[1], os.Args[2])
	if cerr != nil { panic(cerr) }

	// now check the map
	zreadStart := time.Now()
	mcc,_ := nedmap.ReadZCompressedMap(cmpfname)
	zreadElapsed := time.Since(zreadStart)
	fmt.Printf("Compressed readtime: %s\n", zreadElapsed)

	var min, max float32
	min = 1e9
	max = 0.0
	
	for i := range rawE.Elevation {
		for j := range rawE.Elevation[i] {
			diff := rawE.Elevation[i][j] - mcc.Elevation[i][j]
			if diff < 0.0 { diff = -1.0 * diff }
			if min > diff { min = diff }
			if max < diff { max = diff }
		}
	}

	fmt.Printf("%s to %s min/max diff over all samples [%f %f]\n",
		os.Args[2], cmpfname,
		min, max)
}

