package internal

type Mapping struct {
	Day   uint
	Index int
}

func ConvertDataset(mappings string) (map[uint]int, []Mapping) {
	mappingLen := len(mappings) >> 1
	results := make([]Mapping, mappingLen)
	holidays := make(map[uint]int, mappingLen)
	day := uint(0)
	j := 0
	for i := 0; i < mappingLen; i++ {
		day += uint(mappings[j])
		j++
		index := int(mappings[j])
		j++

		holidays[day] = index
		results[i] = Mapping{
			Day:   day,
			Index: index,
		}
	}
	return holidays, results
}
