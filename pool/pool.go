package pool

import (
	"errors"
	"fmt"
	"net"
	"sync"
)

var (
	ErrClosed = errors.New("pool is closed")
)

type Pool interface {
	Get() (net.Conn, error)
	Len() int
	Close()
}

type pool struct {
	sync.RWMutex
	conns   chan net.Conn
	factory Factory
}

type Factory func() (net.Conn, error)

func New(initialCap, maxCap int, factory Factory) (Pool, error) {
	if initialCap < 0 || maxCap <= 0 || initialCap > maxCap {
		return nil, errors.New("invalid capacity settings")
	}

	p := &pool{
		conns:   make(chan net.Conn, maxCap),
		factory: factory,
	}

	for i := 0; i < initialCap; i++ {
		conn, err := factory()
		if err != nil {
			p.Close()
			return nil, fmt.Errorf("factory is not able to fill the pool %s", err)
		}
		p.conns <- conn
	}
	return p, nil
}

func (p *pool) Get() (net.Conn, error) {
	conns, factory := p.getConnsAndFactory()
	if conns == nil {
		return nil, ErrClosed
	}
	select {
	case conn := <-conns:
		if conn == nil {
			return nil, ErrClosed
		}
		return p.wrapConn(conn), nil
	default:
		conn, err := factory()
		if err != nil {
			return nil, err
		}
		return p.wrapConn(conn), nil
	}
}

func (p *pool) Len() int {
	conns, _ := p.getConnsAndFactory()
	return len(conns)
}

func (p *pool) Close() {
	p.Lock()
	conns := p.conns
	p.conns = nil
	p.factory = nil
	p.Unlock()

	if conns == nil {
		return
	}

	close(conns)
	for conn := range conns {
		conn.Close()
	}
}

func (p *pool) put(conn net.Conn) error {
	if conn == nil {
		return errors.New("connection is nil, rejectiong")
	}

	p.RLock()
	defer p.RUnlock()

	if p.conns == nil {
		return conn.Close()
	}

	select {
	case p.conns <- conn:
		return nil
	default:
		return conn.Close()
	}
}

func (p *pool) getConnsAndFactory() (chan net.Conn, Factory) {
	p.RLock()
	conns := p.conns
	factory := p.factory
	p.RUnlock()
	return conns, factory
}

func (p *pool) wrapConn(conn net.Conn) net.Conn {
	c := &PooledConn{
		pool: p,
	}
	c.Conn = conn
	return c
}
