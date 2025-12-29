package sysconf

import (
	"encoding/json"
	"fmt"
)

type itemJSONMap map[string]struct {
	Type  string `json:"Type"`
	UTF8  string `json:"utf8,omitempty"`
	UTF16 string `json:"utf16,omitempty"`
	HEX   string `json:"hex"`
}

func ToJSON(sys *Sysconf) ([]byte, error) {
	if sys == nil {
		return nil, fmt.Errorf("sysconf is nil")
	}

	itemsMap := make(itemJSONMap)

	for _, item := range sys.Items {
		dataHex := ""
		for _, b := range item.Data {
			dataHex += fmt.Sprintf("%02X", b)
		}

		dataUTF8 := ""
		dataUTF16 := ""

		switch item.Type {
		case BIGARRAY, SMALLARRAY:
			dataUTF8 = string(item.Data)
			dataUTF16 = encodeUTF16(item.Data)
		}

		itemsMap[item.Name] = struct {
			Type  string `json:"Type"`
			UTF8  string `json:"utf8,omitempty"`
			UTF16 string `json:"utf16,omitempty"`
			HEX   string `json:"hex"`
		}{
			Type:  item.Type.String(),
			UTF8:  dataUTF8,
			UTF16: dataUTF16,
			HEX:   dataHex,
		}
	}


	data, err := json.MarshalIndent(itemsMap, "", "  ")

	if err != nil {
		return nil, fmt.Errorf("failed to marshal sysconf to JSON: %w", err)
	}

	return data, nil
}


func encodeUTF16(b []byte) string {
	if len(b)%2 != 0 {
		return "NOT_UTF16_ALIGNED"
	}
	buf := make([]rune, len(b)/2)
	for i := 0; i < len(b); i += 2 {
		buf[i/2] = rune(b[i])<<8 | rune(b[i+1])
	}
	return string(buf)
}

func FromJSON(jsonBytes []byte) (*Sysconf, error) {
	var itemsMap map[string]struct {
		Type  string `json:"Type"`
		UTF8  string `json:"utf8,omitempty"`
		UTF16 string `json:"utf16,omitempty"`
		HEX   string `json:"hex"`
	}

	if err := json.Unmarshal(jsonBytes, &itemsMap); err != nil {
		return nil, err
	}

	var items []Item
	for name, v := range itemsMap {
		itemType, err := parseItemType(v.Type)
		if err != nil {
			return nil, err
		}

		data, err := hexStringToBytes(v.HEX)
		if err != nil {
			return nil, fmt.Errorf("invalid HEX data for item %s: %w", name, err)
		}

		items = append(items, Item{
			Type: itemType,
			Name: name,
			Data: data,
		})
	}

	sys := &Sysconf{
		Items: items,
		EOF:  [4]byte{'S', 'C', 'e', 'd'},
	}
	return sys, nil
}

func parseItemType(s string) (ItemType, error) {
	switch s {
	case "BIGARRAY":
		return BIGARRAY, nil
	case "SMALLARRAY":
		return SMALLARRAY, nil
	case "BYTE":
		return BYTE, nil
	case "SHORT":
		return SHORT, nil
	case "LONG":
		return LONG, nil
	case "LONGLONG":
		return LONGLONG, nil
	case "BOOL":
		return BOOL, nil
	default:
		return 0, fmt.Errorf("unknown ItemType: %s", s)
	}
}

func hexStringToBytes(s string) ([]byte, error) {
	if len(s)%2 != 0 {
		return nil, fmt.Errorf("hex string has odd length")
	}
	b := make([]byte, len(s)/2)
	for i := 0; i < len(s); i += 2 {
		var byteVal byte
		_, err := fmt.Sscanf(s[i:i+2], "%02X", &byteVal)
		if err != nil {
			return nil, fmt.Errorf("invalid hex byte: %s", s[i:i+2])
		}
		b[i/2] = byteVal
	}
	return b, nil
}