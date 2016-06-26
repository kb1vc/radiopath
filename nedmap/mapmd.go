// Metadata reader for XML descriptions of maps from the USGS National map database. 
package nedmap


import (
	"io"
	"encoding/xml"
	"io/ioutil"
	"github.com/kb1vc/radiopath/location"
)

type MapInfo struct {
	name string // the name of this map segment
	ll   location.LatLon // location of the lower left corner of the map
	ur   location.LatLon // location of the upper right corner of the map
	rows int // number of rows in the elevation grid
	cols int // number of collumns in the elevation grid
}


type metadata_x struct {
	Name xml.Name `xml:"metadata"`
	Testint string `xml:"testint"`
	Idinfo idinfo_x  `xml:"idinfo"`
	Spdoinfo spdoinfo_x  `xml:"spdoinfo"`
	Inner string `xml:",innerxml"`
}

type idinfo_x struct {
	Name xml.Name 
	Spdom spdom_x `xml:"spdom"`
}

type spdoinfo_x struct { 
	Name xml.Name 
	Rastinfo rastinfo_x `xml:"rastinfo"`
}

type spdom_x struct {
	Name xml.Name 
	Bounding bounding_x `xml:"bounding"`
}

type bounding_x struct {
	Name xml.Name 
	North float64 `xml:"northbc"`
	South float64 `xml:"southbc"`
	East float64 `xml:"eastbc"`	
	West float64 `xml:"westbc"`
}

type rastinfo_x struct {
	Name xml.Name 
	Rowcount int `xml:"rowcount"`
	Colcount int `xml:"colcount"`
}

func GetInfo(instr io.Reader) (MapInfo, error) {
	xmlContent, err := ioutil.ReadAll(instr)
	if err != nil { panic(err) }

	md := metadata_x{}
	err2 := xml.Unmarshal(xmlContent, &md)
	if err2 != nil {
		panic(err)
	}

	var ret MapInfo
	ret.name = "Idunno"
	ret.ll.Lat, ret.ll.Lon = md.Idinfo.Spdom.Bounding.South, md.Idinfo.Spdom.Bounding.West
	ret.ur.Lat, ret.ur.Lon = md.Idinfo.Spdom.Bounding.North, md.Idinfo.Spdom.Bounding.East
	ret.rows, ret.cols = md.Spdoinfo.Rastinfo.Rowcount, md.Spdoinfo.Rastinfo.Colcount

	return ret, nil
}
