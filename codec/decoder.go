package codec

import (
	"encoding/binary"
	"io"
)

func NewResponse() *DubboResponse {
	return &DubboResponse{}
}

func Read(r io.Reader) (*DubboResponse, error) {
	m := NewResponse()
	b, err := m.Decode(r)
	if err != nil {
		return nil, err
	}
	m.ReturnValue = b
	return m, nil
}

func (m *DubboResponse) Decode(r io.Reader) ([]byte, error) {
	_, err := io.ReadFull(r, m.DubboHeader[:])
	if err != nil {
		return nil, err
	}

	length := int(binary.BigEndian.Uint32(m.DubboHeader[12:16]))
	body := make([]byte, length)
	_, err = io.ReadFull(r, body[:])
	if err != nil {
		return nil, err
	}

	return body[2 : length-1], nil
}
