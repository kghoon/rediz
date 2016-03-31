package test

import (
	"testing"
	"runtime"
	"sync"
	"fmt"
	"rediz"
	"github.com/garyburd/redigo/redis"
)

func TestSyncDo(t *testing.T) {

	runtime.GOMAXPROCS(runtime.NumCPU())

	wg := new(sync.WaitGroup)

	c, err := rediz.NewConn(RedisAddress)
	if err != nil {
		t.Fatal(err)
	}
	defer c.Close()

	if _, err := redis.String(c.SyncDo("FLUSHDB")); err != nil {
		t.Fatal(err)
	}

	for i := 0; i < 100; i++ {
		wg.Add(1)

		go func(i int) {
			_, err := c.SyncDo("INCR", "test_sync")
			if err != nil {
				t.Fatal(err)
			}
			wg.Done()
		}(i)
	}

	wg.Wait()
	result, err := redis.Int64(c.SyncDo("INCR", "test_sync"))
	if err != nil {
		t.Fatal(err)
	}

	if result != 101 {
		t.Errorf("Expect 101 but %d\n", result)
	}

	fmt.Println("END")
}