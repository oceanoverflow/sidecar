package codec

import (
	"bytes"
	"sync/atomic"
)

func init() {
	uniqueID = uint64(0)
}

var uniqueID uint64

func GetUniqueID() uint64 {
	return atomic.AddUint64(&uniqueID, uint64(1))
}

func NewRequest(arguments []byte) *DubboRequest {
	dubboHeader := DubboHeader([16]byte{})
	dubboHeader[0] = MagicHigh
	dubboHeader[1] = MagicLow
	dubboHeader.SetMessageType(Request)
	dubboHeader.SetTwoWay(true)
	dubboHeader.SetEvent(false)
	dubboHeader.SetSerializeType(JSON)
	dubboHeader.SetRequestID(GetUniqueID())

	return &DubboRequest{
		DubboHeader:    &dubboHeader,
		DubboVersion:   "2.0.1",
		ServiceName:    "com.alibaba.dubbo.performance.demo.provider.IHelloService",
		MethodName:     "hash",
		ParameterTypes: "Ljava/lang/String;",
		Arguments:      arguments,
		Attachment:     map[string]string{"path": "com.alibaba.dubbo.performance.demo.provider.IHelloService"},
	}
}

func (r *DubboRequest) Encode() []byte {
	body := r.EncodeBody()
	length := len(body)
	r.DubboHeader.SetDataLength(uint32(length))

	ret := make([]byte, length+16)
	copy(ret[:16], r.DubboHeader[:])
	copy(ret[16:], body[:])

	return ret
}

func (r *DubboRequest) EncodeBody() []byte {
	var b bytes.Buffer
	b.Write(WriteObject(r.DubboVersion))
	b.Write(WriteObject(r.ServiceName))
	b.Write(WriteObject(nil))
	b.Write(WriteObject(r.MethodName))
	b.Write(WriteObject(r.ParameterTypes))
	b.Write(WriteObject(r.Arguments))
	b.Write(WriteObject(r.Attachment))
	return b.Bytes()
}
