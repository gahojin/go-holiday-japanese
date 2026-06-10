package internal

type Mapping struct {
	Day   uint32
	Index uint8
}

func ConvertDataset(mappings string) (map[uint32]uint8, []Mapping) {
	mappingLen := len(mappings) >> 1
	results := make([]Mapping, mappingLen)
	holidays := make(map[uint32]uint8, mappingLen)
	day := uint32(0)
	j := 0
	for i := range mappingLen {
		day += uint32(mappings[j])
		j++
		index := mappings[j]
		j++

		holidays[day] = index
		results[i] = Mapping{
			Day:   day,
			Index: index,
		}
	}
	return holidays, results
}
