package codec

import (
	"encoding/binary"
	"io"
)

func Read(r io.Reader) (*DubboResponse, error) {
	m := &DubboResponse{}
	err := m.Decode(r)
	if err != nil {
		return nil, err
	}
	return m, nil
}

func (m *DubboResponse) Decode(r io.Reader) error {
	_, err := io.ReadFull(r, m.DubboHeader[:])
	if err != nil {
		return err
	}

	len := int(binary.BigEndian.Uint32(m.DubboHeader[12:16]))
	data := make([]byte, len)
	_, err = io.ReadFull(r, data[:])
	if err != nil {
		return err
	}
	// do something with the json data

	return nil
}
