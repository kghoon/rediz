package rediz

import (
	"sync"

	"github.com/garyburd/redigo/redis"
)

type ResourceConn struct {
	*redis.PubSubConn
	mutex *sync.Mutex
}

// Close will close the redis connection.
func (r *ResourceConn) Close() {
	r.Conn.Close()
}

func (r *ResourceConn) SyncDo(commandName string, args ...interface{}) (reply interface{}, err error) {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	if len(args) == 0 {
		// To avoid error by weird parameter check of redigo
		return r.Conn.Do(commandName)
	}

	return r.Conn.Do(commandName, args...)
}

// NewConn will return a connection with server
func NewConn(address string) (ResourceConn, error) {
	c, err := redis.Dial("tcp", address)

	return ResourceConn{
		PubSubConn: &redis.PubSubConn{c},
		mutex:      new(sync.Mutex),
	}, err
}
