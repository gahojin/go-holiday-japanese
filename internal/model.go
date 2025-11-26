package internal

import (
	"bytes"
	"encoding/binary"
	"encoding/gob"
	"fmt"
)

type StoreMapping struct {
	Diff  uint8
	Index uint8
}

type StoreData struct {
	Names   []string
	Mapping []StoreMapping
}

type Mapping struct {
	Day   uint
	Index int
}

type ParsedData struct {
	Names    []string
	Holidays map[uint]int
	Mapping  []Mapping
}

func (m *StoreMapping) GobEncode() ([]byte, error) {
	buf := bytes.NewBuffer(nil)
	err := binary.Write(buf, binary.LittleEndian, m.Diff)
	if err != nil {
		return nil, err
	}
	err = binary.Write(buf, binary.LittleEndian, m.Index)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func (m *StoreMapping) GobDecode(data []byte) error {
	buf := bytes.NewReader(data)
	err := binary.Read(buf, binary.LittleEndian, &m.Diff)
	if err != nil {
		return err
	}
	err = binary.Read(buf, binary.LittleEndian, &m.Index)
	if err != nil {
		return err
	}
	return nil
}

func DecodeDataset(rawDataset []byte) (*ParsedData, error) {
	reader := bytes.NewReader(rawDataset)
	decoder := gob.NewDecoder(reader)

	var data StoreData
	if err := decoder.Decode(&data); err != nil {
		return nil, fmt.Errorf("failed to decode dataset: %w", err)
	}

	namesLen := len(data.Names)
	mappingLen := len(data.Mapping)
	mappings := make([]Mapping, mappingLen)
	holidays := make(map[uint]int, mappingLen)
	day := uint(0)
	for i, mapping := range data.Mapping {
		index := int(mapping.Index)
		// 名称インデックスのチェック
		if index >= namesLen {
			return nil, fmt.Errorf("invalid dataset. index overflow %d >= %d", index, namesLen)
		}

		day += uint(mapping.Diff)
		holidays[day] = index
		mappings[i] = Mapping{
			Day:   day,
			Index: index,
		}
	}
	return &ParsedData{
		Names:    data.Names,
		Holidays: holidays,
		Mapping:  mappings,
	}, nil
}
