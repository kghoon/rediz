package test

import (
	"testing"
	"rediz"
	"github.com/garyburd/redigo/redis"
)

func TestConnPool(t *testing.T) {
	pool := rediz.NewConnPool(RedisAddress, Capacity, MaxCap, Timeout)

	conn, err := pool.GetConn()
	if err != nil {
		t.Fail()
	}
	defer pool.PutConn(conn)

	if _, err := redis.Int64(conn.Do("RPUSH", "TestKey", 1)); err != nil {
		t.Fail()
	}

	if n, err := redis.Int64(conn.Do("LPOP", "TestKey")); err != nil || n != 1{
		t.Fail()
	}
}
