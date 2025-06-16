package holiday

import (
	"bytes"
	"encoding/binary"
	"time"
)

type Mapping struct {
	Day   uint16
	Index uint8
}

type Name struct {
	Ja string `json:"ja"`
	En string `json:"en"`
}

type StoreData struct {
	Names   []Name
	Mapping []Mapping

	holidaySet map[uint16]int
}

type Holiday struct {
	Name
	Date time.Time `json:"date"`
}

func (m *Mapping) GobEncode() ([]byte, error) {
	buf := bytes.NewBuffer(nil)
	err := binary.Write(buf, binary.LittleEndian, m.Day)
	if err != nil {
		return nil, err
	}
	err = binary.Write(buf, binary.LittleEndian, m.Index)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func (m *Mapping) GobDecode(data []byte) error {
	buf := bytes.NewReader(data)
	err := binary.Read(buf, binary.LittleEndian, &m.Day)
	if err != nil {
		return err
	}

	err = binary.Read(buf, binary.LittleEndian, &m.Index)
	if err != nil {
		return err
	}
	return nil
}
