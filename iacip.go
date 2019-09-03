package iacip

import (
	"bytes"
	"encoding/csv"
	"io"
	"net"
	"sort"
)

// Finder _
type Finder struct {
	elements []element
}

// New _
func New(areaReader, ipReader io.Reader) Finder {
	areas := mustReadCSV(areaReader)
	m := make(map[string]location, len(areas)-1)
	for i := 1; i < len(areas); i++ {
		m[areas[i][0]] = locationOf(areas[i])
	}

	var tmp []element
	for _, line := range mustReadCSV(ipReader) {
		tmp = append(tmp, element{
			low: net.ParseIP(line[0]).To16(),
			l:   m[line[2]],
		})
	}

	return Finder{
		elements: tmp,
	}
}

// Lookup _
func (f Finder) Lookup(ipString string) (country, region, city string) {
	ip := net.ParseIP(ipString)
	if ip == nil {
		return
	}

	trial := ip.To16()
	if trial == nil {
		return
	}

	i := sort.Search(len(f.elements), func(i int) bool { return bytes.Compare(f.elements[i].low, trial) > 0 }) - 1

	if i >= 0 {
		l := f.elements[i].l
		return l.Country, l.Region, l.City
	}

	return
}

type location struct {
	Country string
	Region  string
	City    string
}

type element struct {
	low []byte
	l   location
}

const (
	global  = "全球"
	china   = "中国大陆"
	unknown = "NULL"
)

func locationOf(data []string) (l location) {
	switch f1, f2 := data[1], data[2]; {
	case f2 == unknown:

	case f2 == global:
		l.Country = f1

	case f2 == china:
		l.Country, l.Region = f2, f1

	default:
		l.Country, l.Region, l.City = china, f2, f1
	}

	return
}

func mustReadCSV(reader io.Reader) [][]string {
	lines, err := csv.NewReader(reader).ReadAll()
	if err != nil {
		panic(err)
	}

	return lines
}
