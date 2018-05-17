package pool

type LeakyBuffer struct {
	bufSize  int
	freeList chan []byte
}

func NewLeakyBuffer(n, bufSize int) *LeakyBuffer {
	return &LeakyBuffer{
		bufSize:  bufSize,
		freeList: make(chan []byte, n),
	}
}

func (lb *LeakyBuffer) Get() (b []byte) {
	select {
	case b = <-lb.freeList:
	default:
		b = make([]byte, lb.bufSize)
	}
	return
}

func (lb *LeakyBuffer) Put(b []byte) {
	if len(b) != lb.bufSize {
		panic("invalid buffer size that's put into leaky buffer")
	}
	select {
	case lb.freeList <- b:
	default:
	}
	return
}
