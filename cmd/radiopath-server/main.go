/*
Copyright (c) 2012, Matthew H. Reilly (kb1vc)
All rights reserved.

Redistribution and use in source and binary forms, with or without
modification, are permitted provided that the following conditions are
met:

    Redistributions of source code must retain the above copyright
    notice, this list of conditions and the following disclaimer.
    Redistributions in binary form must reproduce the above copyright
    notice, this list of conditions and the following disclaimer in
    the documentation and/or other materials provided with the
    distribution.

THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS
"AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT
LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR
A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT
HOLDER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL,
SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT
LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE,
DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY
THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
(INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
*/




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
	
	fmt.Println(wach.OnPath(brg, dist))
}
