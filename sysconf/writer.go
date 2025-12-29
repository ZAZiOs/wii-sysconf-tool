package sysconf

import (
	"encoding/binary"
	"fmt"
)

func Write(sys *Sysconf) ([]byte, error) {
	if sys == nil {
		return nil, fmt.Errorf("sysconf is nil")
	}

	data := make([]byte, 0x4000)

	Magic := [4]byte{'S', 'C', 'v', '0'}
	itemCount := len(sys.Items)
	itemOffsets := make([]uint16, itemCount)

	copy(data[0:4], Magic[:])
	binary.BigEndian.PutUint16(data[4:6], uint16(itemCount))

	currentOffset := uint16(6 + itemCount*2 + 2) // header size

	for i, item := range sys.Items {
		itemOffsets[i] = currentOffset

		nameLen := len(item.Name)
		if nameLen > 31 {
			return nil, fmt.Errorf("item name too long: %s", item.Name)
		}

		headerByte := byte(item.Type<<5) | byte(nameLen-1)
		data[currentOffset] = headerByte
		currentOffset++

		copy(data[currentOffset:currentOffset+uint16(nameLen)], []byte(item.Name))
		currentOffset += uint16(nameLen)
		switch item.Type {
			case BIGARRAY:
				if len(item.Data) > 0xFFFF {
					return nil, fmt.Errorf("BIGARRAY too large")
				}
				binary.BigEndian.PutUint16(data[currentOffset:currentOffset+2], uint16(len(item.Data)-1))
				currentOffset += 2
				copy(data[currentOffset:currentOffset+uint16(len(item.Data))], item.Data)
				currentOffset += uint16(len(item.Data))
			case SMALLARRAY:
				if len(item.Data) > 0xFF {
					return nil, fmt.Errorf("SMALLARRAY too large")
				}
				data[currentOffset] = byte(len(item.Data) - 1)
				currentOffset++
				copy(data[currentOffset:currentOffset+uint16(len(item.Data))], item.Data)
				currentOffset += uint16(len(item.Data))
			case BYTE:
				if len(item.Data) != 1 {
					return nil, fmt.Errorf("BYTE item must have exactly 1 byte of data")
				}
				data[currentOffset] = item.Data[0]
				currentOffset++
			case SHORT:
				if len(item.Data) != 2 {
					return nil, fmt.Errorf("SHORT item must have exactly 2 bytes of data")
				}
				copy(data[currentOffset:currentOffset+2], item.Data)
				currentOffset += 2
			case LONG:
				if len(item.Data) != 4 {
					return nil, fmt.Errorf("LONG item must have exactly 4 bytes of data")
				}
				copy(data[currentOffset:currentOffset+4], item.Data)
				currentOffset += 4
			case LONGLONG:
				if len(item.Data) != 8 {
					return nil, fmt.Errorf("LONGLONG item must have exactly 8 bytes of data")
				}
				copy(data[currentOffset:currentOffset+8], item.Data)
				currentOffset += 8
			case BOOL:
				if len(item.Data) != 1 {
					return nil, fmt.Errorf("BOOL item must have exactly 1 byte of data")
				}
				data[currentOffset] = item.Data[0]
				currentOffset++
			default:
				return nil, fmt.Errorf("unknown item type %d for item %s", item.Type, item.Name)
		}

	}

	for i, off := range itemOffsets {
		pos := 6 + i*2
		binary.BigEndian.PutUint16(data[pos:pos+2], off)
	}

	offsetPast := currentOffset
	binary.BigEndian.PutUint16(data[6+itemCount*2:6+itemCount*2+2], offsetPast)
	
	copy(data[0x3FFC:0x4000], sys.EOF[:])

	return data, nil
}