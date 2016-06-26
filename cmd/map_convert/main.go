package main


import (
	"fmt"
	"os"
	"github.com/kb1vc/radiopath/nedmap"
)


func main() {
	
	infd, err := os.Open(os.Args[1])
	if err != nil { panic(err) }			
	mi, err := nedmap.GetInfo(infd)
	infd.Close()
	if err != nil {
		panic(err)
	}
	fmt.Println(mi)

	evfd, eerr := os.Open(os.Args[2])
	if eerr != nil { panic(eerr) }
	me, merr := nedmap.GetFloatMap(evfd, mi)
	evfd.Close()
	if merr != nil {
		panic(merr)
	}

	fmt.Printf("elevation[0][0] =  %f\n", me.Elevation[0][0])
	fmt.Printf("elevation[253][131] =  %f\n", me.Elevation[253][131])	

	ofd, oerr := os.Create(os.Args[3])
	if oerr != nil { panic(oerr) }
	me.WriteCompressedMap(ofd);
	ofd.Close()

	ofdi, ierr := os.Open(os.Args[3])
	if ierr != nil { panic(ierr) }
	mc,_ := nedmap.ReadCompressedMap(ofdi)
	ofdi.Close()

	fmt.Printf("elevation[0][0] =  %f\n", mc.Elevation[0][0])
	fmt.Printf("elevation[253][131] =  %f\n", mc.Elevation[253][131])	


	zfname := os.Args[3] + ".gz"

	me.WriteZCompressedMap(zfname)

	mcc,_ := nedmap.ReadZCompressedMap(zfname)

	var min, max float32
	min = 1e9
	max = 0.0
	histo := make([]int, 21)
	
	for i := range me.Elevation {
		for j := range me.Elevation[i] {
			diff := me.Elevation[i][j] - mcc.Elevation[i][j]
			idx := int(diff * 10.0) + 10
			if idx < 0 { idx = 0 }
			if idx > 20 { idx = 20 }
			histo[idx]++
			
			if diff < 0.0 { diff = -1.0 * diff }
			if min > diff { min = diff }
			if max < diff { max = diff }
		}
	}

	fmt.Printf("min/max diff over all samples [%f %f]\n", min, max)

	for idx, v := range histo {
		fmt.Printf("%d %d\n", idx, v)
	}
}

