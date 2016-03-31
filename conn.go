package rediz

import (
	"github.com/garyburd/redigo/redis"
	"sync"
)

type ResourceConn struct {
	redis.Conn
	mutex *sync.Mutex
	psc   *PubSubConn
}

// Close will close the redis connection.
func (r ResourceConn) Close() {
	r.Conn.Close()
}

// SyncDo is doing a job (send a command and return response) with lock
func (r *ResourceConn) SyncDo(commandName string, args... interface{}) (reply interface{}, err error) {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	if len(args) == 0 {
		// To avoid error by weird parameter check of redigo
		return r.Do(commandName)
	}

	return r.Do(commandName, args...)
}

// NewConn will return a connection with server
func NewConn(address string) (ResourceConn, error) {
	c, err := redis.Dial("tcp", address)

	return ResourceConn{
		Conn: c,
		mutex: new(sync.Mutex),
	}, err
}