/* 
 Elevation data for maps from the USGS National map database

 This package reads and writes compressed versions of the USGS NED
 elevation datasets and also reads "flt" format data files and XML
 metadata files common to the USGS NED datasets. It reformats the data
 into a compressed form that takes advantage of attributes of our
 natural world
   
    1. Most earth contours follow a gentle slope.  The elevation of 
       adjacent points on the earth rarely differs by more than 4 
       meters or so.  
    2. If we recode the scan rows in each NED file into 4 bit deltas
       (covering the range -7..+7 meters) with an escape (coded as 
       a nybble 0b1000) for bigger jumps, then we will realize a 
       nearly 8:1 compression rate over the 32 bit float elevations
    3. The second derivative of elevation is even more predictable
       than the first.  This allows normal compression algorithms like
       gzip to achieve nearly 2:1 compression rates, vs. 1.5:1 
       compression for the raw floating point input files. 

 The compressed elevation files are reduced in resolution: elevations
 are rounded to the nearest meter.  If your application requires better
 resolution, then don't use this package.

 Output files may be raw delta files, or gzipped delta files.  
 Raw files follow this format 

 The Header: 

 byte: escMarker
 int16: startOfFileMarker
 int16: fileFormatID
 float32: ll.Lat ll.Lon ur.Lat ur.Lon  -- the southwest and northeast corners for this map
 int16: rowcount colcount


 Each Row:
 nybble: escMarker
 int16: startOfLineMarker
 float32: starting elevation for this row
 followed by a string of 'colcount' delta/point specifications

 delta/point specification
 nybble: delta

 If delta != 0x8 then expand the nybble into an int in the range -7:7 inclusive
 add it to the previous elevation

 Otherwise,  the next four bytes (perhaps nybble aligned!) are the elevation of
 the current point. 

*/
package nedmap

import (
	"io"
	"os"
	"errors"
	"encoding/binary"
	"compress/gzip"
	"math"
	"bytes"
	"fmt"
)

// This fully describes a map. It contains the metadata (in the MD field -- see MapInfo)
// and a 2D array of elevation points.
//
// We can build a MapData object either from a raw floating point ned source file
// or from a compressed elevation file. 
type MapData struct {
	MD        MapInfo
	Elevation [][]float32
}


const startOfLineMarker int16 = 0x7fff
const startOfFileMarker int16 = 0x7ffe
const escMarker byte = 0x8
const currentFileVersion int16 = 0x0100

// Read a raw USGS input stream, given the metadata that defines its shape.
// In this case, the metadata probably came from an XML specification file.
func GetFloatMap(instr io.Reader, meta MapInfo) (* MapData, error) {
	m := &MapData{MD: meta, Elevation: make([][]float32, meta.rows)}
	for i := range m.Elevation {
		m.Elevation[i] = make([]float32, meta.cols)
		rerr := binary.Read(instr, binary.LittleEndian, m.Elevation[i])
		if rerr != nil {
			return m, rerr
		}
	}
	return m, nil
}

// Write a floating point map file in the USGS ned file format
// This is rarely used, as the file is often 15 times bigger than
// a compressed datafile.
func (m * MapData) WriteFloatMap(outstr io.Writer) (error) {
	for i := range m.Elevation {
		err := binary.Write(outstr, binary.LittleEndian, m.Elevation[i])
		if err != nil {
			return err
		}
	}
	return nil
}

const inbufSize int = 1024

type nybbleOutStream struct {
	wr io.Writer
	odd bool
	cur byte
}

type nybbleInStream struct {
	rd io.Reader
	odd bool
	inbuf [inbufSize]byte
	in_idx int
}

func (w * nybbleInStream) initIn() {
	w.in_idx = 0
	w.rd.Read(w.inbuf[:])
	w.odd = false
}

func (w * nybbleOutStream) terminateOut() {
	if w.odd {
		buf := []byte{w.cur}
		w.wr.Write(buf)
	}

	w.odd = false
	w.cur = 0
}

func (w * nybbleOutStream) putNybble(v byte) {
	if w.odd {
		w.cur = w.cur | ((v & 0xf) << 4)
		buf := []byte{w.cur}		
		w.wr.Write(buf)
		w.cur = 0
		w.odd = false
	} else {
		w.cur = v & 0xf
		w.odd = true
	}
}

func (w * nybbleInStream) getNybble() (byte) {
	var r byte
	v := w.inbuf[w.in_idx]

	if w.odd {
		r = (v >> 4) & 0xf
		w.odd = false
		w.in_idx++
		if w.in_idx == inbufSize {
			w.rd.Read(w.inbuf[:])
			w.in_idx = 0
		}
	} else {
		r = v & 0xf
		w.odd = true
	}
	return r
}

func (w * nybbleOutStream) putInt16(v int16) {
	w.putNybble(byte(v))
	w.putNybble(byte(v >> 4))	
	w.putNybble(byte(v >> 8))
	w.putNybble(byte(v >> 12))	
}


func (w * nybbleOutStream) put(v interface{}) {
	b := new(bytes.Buffer)
	err := binary.Write(b, binary.LittleEndian, v)
	if err != nil {
		panic(err)
	}

	vec := b.Bytes()
	
	for i := 0; i < len(vec); i++ {		
		vv := vec[i]
		w.putNybble(vv)
		w.putNybble(vv >> 4)		
	}
}

func (w * nybbleInStream) get(v * interface{}) {
	
}

func (w * nybbleInStream) getBytes(buf []byte) {
	for i := 0; i < len(buf); i++ {
		vl := w.getNybble()
		vh := w.getNybble()
		buf[i] = vl | (vh << 4)
	}
}

func (w * nybbleInStream) getInt16() (r int16) {
	buf := make([]byte, 2)
	w.getBytes(buf)
	r = int16(buf[0]) + int16(buf[1]) << 8
	return
}

func (w * nybbleInStream) getFloat32() (r float32) {
	buf := make([]byte, 4)
	w.getBytes(buf)	
	bb := bytes.NewBuffer(buf)
	binary.Read(bb, binary.LittleEndian, &r)
	return
}



func (w * nybbleOutStream) writeCompHeader(md * MapInfo) (error) {

	// we write lat/lon then rows/cols 
	// ll:(float64, float64) ur:(float64 float64) rows:16 int cols: int16
	w.put(escMarker)	
	w.put(startOfFileMarker)
	w.put(currentFileVersion)	
		
	llur := []float32{ float32(md.ll.Lat), float32(md.ll.Lon), float32(md.ur.Lat), float32(md.ur.Lon) }
	for _,v := range llur {
		w.put(v)
	}

	w.put(int16(md.rows))
	w.put(int16(md.cols))

	return nil
}


func (w * nybbleInStream) readCompHeader(md * MapInfo) (error) {
	// read the escMarker
	em := w.getNybble()
	w.getNybble()
	if em != escMarker {
		fmt.Printf("got bad nybble, expected esc got %02x\n", em)
	}
	
	// now the SOF marker
	sof := w.getInt16()
	if sof != startOfFileMarker {
		fmt.Printf("got bad startOfFileMarker, expected %04x got %04x\n", startOfFileMarker, sof)
	}

	// get the file version ID .. we don't care much
	fvid := w.getInt16()
	if fvid != currentFileVersion { 
		return errors.New(fmt.Sprintf("Got bad file version id: %04x  expected %04x\n", fvid, currentFileVersion))
	}
	
	md.ll.Lat = float64(w.getFloat32())
	md.ll.Lon = float64(w.getFloat32())
	md.ur.Lat = float64(w.getFloat32())
	md.ur.Lon = float64(w.getFloat32())

	md.rows = int(w.getInt16())
	md.cols = int(w.getInt16())
	
	return nil
}


func (w * nybbleOutStream) writeCompRowStart(el float32) (error) {
	w.putNybble(escMarker)
	w.put(startOfLineMarker)
	w.put(round(el))
	return nil
}

func round(v float32) (r float32) {
	v64 := float64(v)
	var r64 float64
	if v64 > 0.0 {
		r64 = math.Floor(v64 + 0.5)
	} else {
		r64 = math.Ceil(v64 - 0.5)
	}

	r = float32(r64)
	return r
}


func (w * nybbleOutStream) writeCompElevation(el, lel float32) (error) {
	diff := round(el) - round(lel)
	idiff := int(diff)
	if (idiff <= 7) && (idiff >= -7) {
		w.putNybble(byte(idiff & 0xf))
	} else {
		w.putNybble(escMarker)
		w.put(round(el))
	}

	return nil
}

func (m * MapData) WriteZCompressedMap(fname string) (error) {
	ofd, oerr := os.Create(fname)
	if oerr != nil {
		panic(oerr)
	}

	
	wr := gzip.NewWriter(ofd)
	wcerr := m.WriteCompressedMap(wr)
	wr.Close()
	return wcerr
}

func (m * MapData) WriteCompressedMap(outstr io.Writer) (error) {
	// create a nybble stream
	ns := nybbleOutStream{wr: outstr, odd: false, cur: 0}
	
	// first write the compressed header.
	ns.writeCompHeader(&m.MD)

	var last_ev float32
	for i := range m.Elevation {
		// row by row...
		for j := range m.Elevation[i] {
			if j == 0 {
				ns.writeCompRowStart(m.Elevation[i][0])
				last_ev = m.Elevation[i][0]
			} else {
				ev := m.Elevation[i][j]
				ns.writeCompElevation(ev, last_ev)
				last_ev = ev
			}
		}
		
	}

	ns.terminateOut()
	
	return nil
}

func (w * nybbleInStream) readCompElevation(last float32) (r float32) {
	// get a nybble.
	v := w.getNybble()
	if v != escMarker {
		// simple -- convert to float, add to last, and write
		// sign extend the nybble
		var iv int
		if (v & 0x8) != 0 {
			iv = int((^v + 1) & 0x7)
			iv = 0 - iv
		} else {
			iv = int(v & 0x7)
		}
		r = last + float32(iv)
	} else {
		r = w.getFloat32()
	}
	return r
}

func (w * nybbleInStream) readCompRowStart() (r float32) {
	// each row starts with esc, then startOfLineMarker
	em := w.getNybble()
	if em != escMarker {
		fmt.Printf("got bad nybble, expected esc got %02x\n", em)
	}
	// now the SOF marker
	sof := w.getInt16()
	if sof != startOfLineMarker {
		fmt.Printf("got bad startOfLineMarker, expected %04x got %04x\n", startOfFileMarker, sof)
	}
	
	// now get the elevation
	r = w.getFloat32()

	return
}

func ReadZCompressedMap(fname string) (* MapData, error) {
	ifd, ierr := os.Open(fname)
	if ierr != nil { panic(ierr) }
	defer ifd.Close()

	rd, gzerr := gzip.NewReader(ifd)
	if gzerr != nil { panic(gzerr) }
	return ReadCompressedMap(rd)
}

func ReadCompressedMap(instr io.Reader) (* MapData, error) {
	ns := nybbleInStream{rd: instr, odd: false}
	ns.initIn()
	
	m := new(MapData)

	ns.readCompHeader(&m.MD)

	fmt.Println(m.MD)

	m.Elevation = make([][]float32, m.MD.rows)

	for i := range m.Elevation {
		m.Elevation[i] = make([]float32, m.MD.cols)

		// now read each row
		m.Elevation[i][0] = ns.readCompRowStart()
		last_el := m.Elevation[i][0]
		
		for j := 1; j < m.MD.cols; j++ {
			// get the next elevation
			el := ns.readCompElevation(last_el)
			m.Elevation[i][j] = el
			last_el = el
		}
	}
	
	fmt.Printf("el[0][0] = %f\n", m.Elevation[0][0])
	
	return m, nil	
}
