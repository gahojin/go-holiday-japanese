package internal

import (
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

func ConvertDataset(names []string, mappings []StoreMapping) (*ParsedData, error) {
	namesLen := len(names)
	mappingLen := len(mappings)
	results := make([]Mapping, mappingLen)
	holidays := make(map[uint]int, mappingLen)
	day := uint(0)
	for i, mapping := range mappings {
		index := int(mapping.Index)
		// 名称インデックスのチェック
		if index >= namesLen {
			return nil, fmt.Errorf("invalid dataset. index overflow %d >= %d", index, namesLen)
		}

		day += uint(mapping.Diff)
		holidays[day] = index
		results[i] = Mapping{
			Day:   day,
			Index: index,
		}
	}
	return &ParsedData{
		Names:    names,
		Holidays: holidays,
		Mapping:  results,
	}, nil
}
