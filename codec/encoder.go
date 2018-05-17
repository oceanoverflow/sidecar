package codec

import (
	"io"
	"sync/atomic"
)

func init() {
	uniqueID = uint64(0)
}

var uniqueID uint64

func GetUniqueID() uint64 {
	return atomic.AddUint64(&uniqueID, uint64(1))
}

func NewRequest() *DubboRequest {
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
		DubboVersion:   "2.6.0",
		ServiceName:    "com.alibaba.dubbo.performance.demo.provider.IHelloService",
		ServiceVersion: "2.0.0",
		MethodName:     "hash",
		ParameterTypes: "Ljava/lang/String;",
	}
}

func (r DubboRequest) Encode() []byte {
	// first and foremost, encode the json
	// get the total length of the data
	// make one byte slice using make([]byte, length)
	// fill the slice with respective data
	// finally return it
	// things to take care of
	return nil
}

func (r DubboRequest) WriteTo(w io.Writer) error {
	_, err := w.Write(r.DubboHeader[:])
	if err != nil {
		return err
	}
	// do something with the extra parts
	return nil
}

func encodeAttachments() {

}
