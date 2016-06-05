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

// Routines to convert grid to latlon or latlon to grid
package location

import (
	"strings"
	"errors"
	"fmt"
	"math"
)

var positions = []string  { "first", "second", "third", "fourth", "fifth", "sixth", "seventh", "eighth" }

func checkgrid(gs string, pos int, min byte, max byte) error {
	for i := 0; i < 2; i++ {
		if (gs[pos + i] < min) || (gs[pos + i] > max) {
			erstr := fmt.Sprintf("Bad grid specification %s character was %c; must be in range [%c, %c] inclusive.", gs, gs[pos+i], min, max)
			return errors.New(erstr)
		}
	}
	return nil
}

func FromGrid(gs string) (LatLon, error) {
	gs = strings.ToUpper(gs)
	
	var ret LatLon
	
	if er := checkgrid(gs, 0, 'A', 'R'); er != nil {
		return ret, er
	}

	ret.lon = 20.0 * float64(gs[0] - 'A') - 180.0
	ret.lat = 10.0 * float64(gs[1] - 'A') - 90.0

	if len(gs) < 4 {
		ret.lon += 10.0;
		ret.lat += 5.0;
		return ret, nil
	}

	if er := checkgrid(gs, 2, '0', '9'); er != nil {
		return ret, er
	}


	ret.lon += 2.0 * float64(gs[2] - '0')
	ret.lat += 1.0 * float64(gs[3] - '0')

	if len(gs) < 6 {
		ret.lon += 1.0
		ret.lat += 0.5
		return ret, nil
	}
		
	if er := checkgrid(gs, 4, 'A', 'X'); er != nil {
		return ret, er
	}


	ret.lon += (2.0/24.0) * float64(gs[4] - 'A')
	ret.lat += (1.0/24.0) * float64(gs[5] - 'A')

	if len(gs) < 8 {
		ret.lon += (1.0 / 24.0)
		ret.lat += (0.5 / 24.0)
		return ret, nil
	}

	if er := checkgrid(gs, 6, '0', '9'); er != nil {
		return ret, er
	}

	ret.lon += (0.2 / 24.0) * (float64(gs[6] - '0') + 0.5)
	ret.lat += (0.1 / 24.0) * (float64(gs[7] - '0') + 0.5)

	
	return ret, nil
}

func ToGrid(ll LatLon) (string, error) {
	var ret string
	ret = "      "

	if (ll.lon < -180.0) || (ll.lon > 180.0) || 
		(ll.lat < -90.0) || (ll.lat > 90.0) {
		erstr := fmt.Sprintf("Lat/Lon coordinate is out of range %f lat %f lon", ll.lat, ll.lon)
		return ret, errors.New(erstr)
	}

	// find the lon chars. 
	lon := ll.lon + 180
	offset := math.Floor(lon / 20.0)
	lon = lon - 20.0 * offset
	ret += string('A' + int(offset))

	lat := ll.lat + 90.0
	offset = math.Floor(lat / 10.0)
	lat = lat - 10.0 * offset
	ret += string('A' + int(offset))

	offset = math.Floor(lon / 2.0)
	lon = lon - 2.0 * offset
	ret += string('0' + int(offset))

	offset = math.Floor(lat / 1.0)
	lat = lat - 1.0 * offset
	ret += string('0' + int(offset))
	
	offset = math.Floor(lon / (2.0 / 24.0))
	ret += string('a' + int(offset))

	offset = math.Floor(lat / (1.0 / 24.0))
	ret += string('a' + int(offset))

	return ret, nil
	
}
