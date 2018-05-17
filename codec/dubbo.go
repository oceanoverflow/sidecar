package codec

import (
	"encoding/binary"
)

const (
	MagicHigh byte = 0xda
	MagicLow  byte = 0xbb
)

type MessageType byte

const (
	Response MessageType = iota
	Request
)

type SerializeType byte

const (
	JSON SerializeType = iota
)

type StatusType uint8

const (
	OK                                StatusType = 20
	CLIENT_TIMEOUT                               = 30
	SERVER_TIMEOUT                               = 31
	BAD_REQUEST                                  = 40
	BAD_RESPONSE                                 = 50
	SERVICE_NOT_FOUND                            = 60
	SERVICE_ERROR                                = 70
	SERVER_ERROR                                 = 80
	CLIENT_ERROR                                 = 90
	SERVER_THREADPOOL_EXHAUSTED_ERROR            = 100
)

type ReturnValueType int

const (
	RESPONSE_WITH_EXCEPTION ReturnValueType = iota
	RESPONSE_VALUE
	RESPONSE_NULL_VALUE
)

// MagicHigh byte
// MagicLow  byte
// Misc      byte
// status    byte
// RequestID uint64
// Length    uint32
type DubboHeader [16]byte

// * Dubbo version
// * Service name
// * Service version
// * Method name
// * Method parameter types
// * Method arguments
// * Attachments
type DubboRequest struct {
	*DubboHeader
	DubboVersion   string
	ServiceName    string
	ServiceVersion string
	MethodName     string
	ParameterTypes string
	Arguments      []byte
	Attachment     map[string]string
}

// * Return value type, identifies what kind of value returns from server side: RESPONSE_NULL_VALUE - 2, RESPONSE_VALUE - 1, RESPONSE_WITH_EXCEPTION - 0.
// * Return value, the real value returns from server.
type DubboResponse struct {
	*DubboHeader
	Type  int
	Value []byte
}

func (h DubboHeader) CheckMagicNumber() bool {
	return (h[0] == MagicHigh) && (h[1] == MagicLow)
}

func (h DubboHeader) MessageType() MessageType {
	return MessageType((h[2] & 0x80) >> 7)
}

func (h *DubboHeader) SetMessageType(mt MessageType) {
	h[2] = h[2] | (byte(mt) << 7)
}

func (h DubboHeader) IsTwoWay() bool {
	return h[2]&0x40 == 0x40
}

func (h *DubboHeader) SetTwoWay(twoway bool) {
	if twoway {
		h[2] = h[2] | 0x40
	} else {
		h[2] = h[2] &^ 0x40
	}
}

func (h DubboHeader) IsEvent() bool {
	return h[2]&0x20 == 0x20
}

func (h *DubboHeader) SetEvent(evt bool) {
	if evt {
		h[2] = h[2] | 0x20
	} else {
		h[2] = h[2] &^ 0x20
	}
}

func (h DubboHeader) SerializeType() SerializeType {
	return SerializeType(h[2] & 0x1f)
}

// ???
func (h *DubboHeader) SetSerializeType(st SerializeType) {
	h[2] = h[2] & byte(st)
}

func (h DubboHeader) StatusType() StatusType {
	return StatusType(h[3])
}

func (h *DubboHeader) SetStatusType(st StatusType) {
	h[3] = byte(st)
}

func (h DubboHeader) RequestID() uint64 {
	return binary.BigEndian.Uint64(h[4:12])
}

func (h *DubboHeader) SetRequestID(id uint64) {
	binary.BigEndian.PutUint64(h[4:12], id)
}

func (h DubboHeader) DataLength() uint32 {
	return binary.BigEndian.Uint32(h[12:16])
}

func (h *DubboHeader) SetDataLength(length uint32) {
	binary.BigEndian.PutUint32(h[13:16], length)
}
