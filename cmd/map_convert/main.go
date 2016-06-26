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
	fmt.Printf("elevation[0][1] =  %f\n", me.Elevation[0][1])	

	ofd, oerr := os.Create(os.Args[3])
	if oerr != nil { panic(oerr) }
	me.WriteCompressedMap(ofd);

	ofd.Close()

	ofdi, ierr := os.Open(os.Args[3])
	if ierr != nil { panic(ierr) }
	nedmap.ReadCompressedMap(ofdi)

	
}

