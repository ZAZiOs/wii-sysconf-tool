package sysconf

import (
	"fmt"
)

type ItemType uint8

const (
	BIGARRAY   ItemType = 1
	SMALLARRAY ItemType = 2
	BYTE       ItemType = 3
	SHORT      ItemType = 4
	LONG       ItemType = 5
	LONGLONG   ItemType = 6
	BOOL       ItemType = 7
)

func (t ItemType) String() string {
	switch t {
	case BIGARRAY:
		return "BIGARRAY"
	case SMALLARRAY:
		return "SMALLARRAY"
	case BYTE:
		return "BYTE"
	case SHORT:
		return "SHORT"
	case LONG:
		return "LONG"
	case LONGLONG:
		return "LONGLONG"
	case BOOL:
		return "BOOL"
	default:
		return fmt.Sprintf("UNKNOWN(%d)", uint8(t))
	}
}


type Header struct {
	Magic       [4]byte   // SCv0
	ItemCount   uint16
	ItemOffsets []uint16
	OffsetPastLastItem uint16
}

type Item struct {
	Type ItemType
	Name string
	Data []byte
}

// full sysconf file structure
type Sysconf struct {
	Header      Header
	Items       []Item
	LookupTable [39]uint16
	EOF         [4]byte // SCed
}
