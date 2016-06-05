// Routines to calculate bearing and distance from one latlon to another
package location

import (
	"math"
)

const clarkeAL = 6378206.4   /*Clarke 1866 ellipsoid*/
const clarkeBL = 6356583.8
const deg2rad  = (math.Pi / 180.0)
const pi2 = (2.0 * math.Pi)

/*Determines bearing and distance based on*/
/*algorithm compensating for earth's shape*/
/* Modified and heavily plagiarized from code that is */
/*Copyright 1993 Michael R. Owen, W9IP & */
/*               Paul Wade, N1BWT */
/*Released to the public domain Feb 23, 1993*/
/* 
** 
** Translated to c via p2c  with appropriate
** post-translation hacks by Matt Reilly (KB1VC). 
*/
/* 
**
** Recoding to be compatible with the DEM/DEM library
** functions by Matt Reilly KB1VC  June, 1996 
*/

/*
** Translated to go by Matt Reilly KB1VC June 2016 
*/

/*NOTE: NORTH latitude, EAST longitude are positive;   */
/*      South & West should be input as negative numbers
**   This is in keeping with the convention from the DEM 
** files from the USGS.  Note that the original N1BWT code
** assumed that negative longitudes were east.
*/

// return (bearing, reverse bearing, distance) in degrees and km respectively
func (fr LatLon) Bearing(to LatLon) (float64, float64, float64) {
	if (math.Abs(fr.lat - to.lat) < 1.0e-10) &&
		(math.Abs(fr.lon - to.lon) < 1.0e-10) {
		return 0.0, 0.0, 0.0
	}

	boa := clarkeBL / clarkeAL
	f := 1.0 - boa
	p1r := fr.lat * deg2rad
	p2r := to.lat * deg2rad
	l1r := fr.lon * deg2rad
	l2r := to.lon * deg2rad
	
	dlr := l1r - l2r

	t1r := math.Atan(boa * math.Tan(p1r))
	t2r := math.Atan(boa * math.Tan(p2r))
	tm := (t1r + t2r) / 2.0
	dtm := (t2r - t1r) / 2.0

	stm := math.Sin(tm)
	ctm := math.Cos(tm)

	sdtm := math.Sin(dtm)
	cdtm := math.Cos(dtm)

	kl := stm * cdtm
	kk := sdtm * ctm

	sdlmr := math.Sin(dlr / 2.0)
	l := sdtm * sdtm + sdlmr * sdlmr * (cdtm * cdtm - stm * stm)
	cd := 1.0 - 2.0 * l
	dl := math.Acos(cd)
	sd := math.Sin(dl)
	t := dl / sd

	u := 2.0 * kl * kl / (1.0 - l)
	v := 2.0 * kk * kk / l
	d := 4.0 * t * t
	x := u + v
	e := -2.0 * cd
	y := u - v
	a := -d * e
	ff64 := f * f / 64.0

	dist := clarkeAL * sd * (t - f / 4.0 * (t * x - y) +
		ff64 * (x * (a + (t - (a + e) / 2.0) * x) +
		y * (e * y - 2.0 * d) + d * x * y )) / 1000.0

	
	tdlpm := math.Tan((dlr - (e * (4.0 - x) + 2.0 * y) * (f / 2.0 * t +
		ff64 * (32.0 * t + (a - 20.0 * t) * x -
		2.0 * (d + 2.0) * y)) / 4.0 * math.Tan(dlr)) / 2.0);

	hapbr := atan2(sdtm, ctm * tdlpm); 
	hambr := atan2(cdtm, stm * tdlpm); 

	a1m2 := pi2 + hambr - hapbr;
	a2m1 := pi2 - hambr - hapbr;


	
	for a1m2 < 0.0 {
		a1m2 += pi2
	}
	for a1m2 >= pi2 {
		a1m2 -= pi2
	}
	for a2m1 < 0.0 {
		a2m1 += pi2
	}
	for a2m1 >= pi2 {
		a2m1 -= pi2
	}

	az := 0.0
	if a1m2 != 0.0 {
		az = 360.0 - a1m2 / deg2rad
	}
	revaz := 0.0
	if a2m1 != 0.0 {
		revaz = 360.0 - a2m1 / deg2rad
	}

	return az, revaz, dist
}

func (fr LatLon) Onpath(az float64, dist float64) LatLon {

	const eps = 5e-14;
	const a = 6378206.4; /* (meters) */
	const  f = 1.0/298.25722210088; 

	/* *** SOLUTION OF THE GEODETIC DIRECT PROBLEM AFTER T.VINCENTY */
	/* *** MODIFIED RAINSFORD'S METHOD WITH HELMERT'S ELLIPTICAL TERMS */
	/* *** EFFECTIVE IN ANY AZIMUTH AND AT ANY DISTANCE SHORT OF ANTIPODAL */

	/* *** A IS THE SEMI-MAJOR AXIS OF THE REFERENCE ELLIPSOID */
	/* *** F IS THE FLATTENING OF THE REFERENCE ELLIPSOID */
	/* *** LATITUDES AND LONGITUDES IN RADIANS POSITIVE NORTH AND EAST */
	/* *** AZIMUTHS IN RADIANS CLOCKWISE FROM NORTH */
	/* *** GEODESIC DISTANCE S ASSUMED IN UNITS OF SEMI-MAJOR AXIS A */

	/* *** PROGRAMMED FOR CDC-6600 BY LCDR L.PFEIFER NGS ROCKVILLE MD 20FEB75 
*/
	/* *** MODIFIED FOR SYSTEM 360 BY JOHN G GERGEN NGS ROCKVILLE MD 750608 */
	/* *** HACKED TO HELL AND GONE BY F2C (July-17-96 version) and by 
           *** Matt Reilly (KB1VC) to make it fit with the dem-gridlib routines.
           *** Feb 24, 1997. */

	/* translated to "go" in June of 2016 by Matt Reilly (kb1vc) */

	var sy, cy, cz float64
	var e, c, d, x, y float64
	var s, faz, glat1, glon1 float64
	var glat2, glon2 float64
	var r, tu, sf, cf, rbaz float64
	var cu, su, sa, c2a float64
	
	/* distance is in Km... */ 
	s = dist * 1000.0
	faz = az * math.Pi / 180.0

	glat1 = fr.lat * math.Pi / 180.0
	glon1 = fr.lon * math.Pi / 180.0 
	
	
	r = 1.0 - f
	
	tu = r * math.Sin(glat1) / math.Cos(glat1)
	
	sf = math.Sin(faz)
	cf = math.Cos(faz)
	rbaz = 0.0
	if (cf != 0.0) {
		rbaz = math.Atan2(tu, cf) * 2.0
	}
	cu = 1.0 / math.Sqrt(tu * tu + 1.0)
	su = tu * cu
	sa = cu * sf
	c2a = -sa * sa + 1.0
	x = math.Sqrt((1.0 / r / r - 1.0) * c2a + 1.0) + 1.0

	x = (x - 2.0) / x
	c = 1.0 - x
	c = (x * x / 4.0 + 1) / c
	d = (x * .375 * x - 1.0) * x
	tu = s / r / a / c
	y = tu

	flag := 1
	for flag > 0 {
		sy = math.Sin(y)
		cy = math.Cos(y)
		cz = math.Cos(rbaz + y)
		e = cz * cz * 2.0 - 1.0
		c = y
		x = e * cy
		y = e + e - 1.0
		y = (((sy * sy * 4.0 - 3.0) * y * cz * d / 6.0 + x) * d / 4.0 - cz) * sy * d + tu
		if (math.Abs(y - c) <= eps) {
			flag = 0
		}
	}
	rbaz = cu * cy * cf - su * sy
	c = r * math.Sqrt(sa * sa + rbaz * rbaz)
	d = su * cy + cu * sy * cf
	glat2 = math.Atan2(d, c)
	c = cu * cy - su * sy * cf
	x = math.Atan2(sy * sf, c)
	c = ((c2a * -3.0 + 4.0) * f + 4.0) * c2a * 
		f / 16.0
	d = ((e * cy * c + cz) * sy * c + y) * sa
	glon2 = glon1 + x - (1.0 - c) * d * f


	var ret LatLon
	
	ret.lat = glat2 * 180.0 / math.Pi
	ret.lon = glon2 * 180.0 / math.Pi 

	
	return ret 
}

// return arctan such that the return value is 
// in the range 0..2pi
func atan2(x float64, y float64) float64 {
retval := math.Atan2(x, y)
if retval < 0.0 {
	retval += 2.0 * math.Pi
}
return retval
}


