package internal

type Mapping struct {
	Day   uint16
	Index uint8
}

type Bitset []byte

func ConvertDataset(mappings string) (Bitset, []Mapping) {
	mappingLen := len(mappings) >> 1
	results := make([]Mapping, mappingLen)
	day := uint16(0)
	j := 0
	for i := range mappingLen {
		day += uint16(mappings[j])
		j++
		index := mappings[j]
		j++

		results[i] = Mapping{
			Day:   day,
			Index: index,
		}
	}
	return newBitset(results), results
}

func newBitset(mapping []Mapping) Bitset {
	// js版と同じく、1bit=1日で保持する
	var maxDay uint16
	for _, m := range mapping {
		if m.Day > maxDay {
			maxDay = m.Day
		}
	}
	b := make(Bitset, (maxDay>>3)+1)
	for _, m := range mapping {
		b[m.Day>>3] |= 1 << (m.Day & 7)
	}
	return b
}

func (b Bitset) Has(day uint16) bool {
	idx := day >> 3
	if int(idx) >= len(b) {
		return false
	}
	return b[idx]&(1<<(day&7)) != 0
}
