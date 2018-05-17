package pool

import (
	"net"
	"sync"
)

type PooledConn struct {
	net.Conn
	sync.RWMutex
	pool     *pool
	unusable bool
}

func (pc *PooledConn) Close() error {
	pc.RLock()
	defer pc.RUnlock()

	if pc.unusable {
		if pc.Conn != nil {
			return pc.Conn.Close()
		}
		return nil
	}
	return pc.pool.put(pc.Conn)
}

func (pc *PooledConn) MarkUnusable() {
	pc.Lock()
	defer pc.Unlock()
	pc.unusable = true
}
