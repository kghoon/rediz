package rediz

import (
	"github.com/youtube/vitess/go/pools"
	"golang.org/x/net/context"
	"time"
)

type ConnPool struct {
	*pools.ResourcePool
}

func NewConnPool(address string, capacity, maxCapacity int, idleTimeout time.Duration) *ConnPool {
	p := pools.NewResourcePool(func() (pools.Resource, error) {
		return NewConn(address)
	}, capacity, maxCapacity, idleTimeout)

	return &ConnPool{p}
}

func (pool *ConnPool) GetConn() (conn *ResourceConn, err error) {
	r, err := pool.Get(context.TODO())
	c := r.(ResourceConn)

	return &c, err
}

func (pool *ConnPool) PutConn(conn *ResourceConn) {
	if pool != nil {
		pool.Put(conn)
	}
}