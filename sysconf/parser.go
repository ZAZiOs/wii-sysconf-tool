package sysconf

import (
	"encoding/binary"
	"fmt"
)

func Parse(data []byte) (*Sysconf, error) {
	
	var magic [4]byte
	copy(magic[:], data[0:4])
	if magic != [4]byte{'S', 'C', 'v', '0'} {
		return nil, fmt.Errorf("invalid magic: got %s, need SCv0", magic)
	}


	itemCount := binary.BigEndian.Uint16(data[4:6])
	itemOffsets := make([]uint16, itemCount)
	for i := 0; i < int(itemCount); i++ {
		itemOffsets[i] = binary.BigEndian.Uint16(data[6+i*2 : 6+i*2+2])
	}

	offsetPastLastItem := binary.BigEndian.Uint16(data[6+int(itemCount)*2 : 6+int(itemCount)*2+2])

	header := Header{
		Magic: magic,
		ItemCount:   itemCount,
		ItemOffsets: itemOffsets,
		OffsetPastLastItem: offsetPastLastItem,
	}

	// reading items

	items := make([]Item, itemCount)

	for i, offset := range itemOffsets {
		// Проверка на корректность смещения
		if offset >= uint16(len(data)) {
			return nil, fmt.Errorf("invalid item offset: %d", offset)
		}

		itemData := data[offset:]
		if len(itemData) < 2 {
			return nil, fmt.Errorf("insufficient data for item %d", i)
		}

		headerByte := data[offset]

		itemType := ItemType(headerByte >> 5)
		nameLen := uint8(headerByte & 0x1F) + 1
		
		nameStart := offset + 1
		nameEnd := nameStart + uint16(nameLen)
		name := string(data[nameStart:nameEnd])
		
		dataStart := nameEnd
		var itemBytes []byte

		switch itemType {
			case BIGARRAY:
				if len(data[dataStart:]) < 2 {
					return nil, fmt.Errorf("insufficient data for BIGARRAY length in item %d", i)
				}
				length := binary.BigEndian.Uint16(data[dataStart : dataStart+2]) + 1
				dataStart += 2
				if len(data[dataStart:]) < int(length) {
					return nil, fmt.Errorf("insufficient data for BIGARRAY in item %d", i)
				}
				itemBytes = data[dataStart : dataStart+length]
			case SMALLARRAY:
				if len(data[dataStart:]) < 1 {
					return nil, fmt.Errorf("insufficient data for SMALLARRAY length in item %d", i)
				}
				length := uint16(data[dataStart]) + 1
				dataStart += 1
				if uint16(len(data[dataStart:])) < length {
					return nil, fmt.Errorf("insufficient data for SMALLARRAY in item %d", i)
				}
				itemBytes = data[dataStart : dataStart+length]
			case BYTE:
				if len(data[dataStart:]) < 1 {
					return nil, fmt.Errorf("insufficient data for BYTE in item %d", i)
				}
				itemBytes = data[dataStart : dataStart+1]
			case SHORT:
				if len(data[dataStart:]) < 2 {
					return nil, fmt.Errorf("insufficient data for SHORT in item %d", i)
				}
				itemBytes = data[dataStart : dataStart+2]
			case LONG:
				if len(data[dataStart:]) < 4 {
					return nil, fmt.Errorf("insufficient data for LONG in item %d", i)
				}
				itemBytes = data[dataStart : dataStart+4]
			case LONGLONG:
				if len(data[dataStart:]) < 8 {
					return nil, fmt.Errorf("insufficient data for LONGLONG in item %d", i)
				}
				itemBytes = data[dataStart : dataStart+8]
			case BOOL:
				if len(data[dataStart:]) < 1 {
					return nil, fmt.Errorf("insufficient data for BOOL in item %d", i)
				}
				itemBytes = data[dataStart : dataStart+1]
			default:
				return nil, fmt.Errorf("unknown item type %d in item %d", itemType, i)
		}
		items[i] = Item{
			Type: itemType,
			Name: name,
			Data: itemBytes,
		}
	}

	sys := &Sysconf{
		Header: header,
		Items: items,
	}

	return sys, nil
}
