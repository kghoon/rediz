package test

import (
	"testing"
	"rediz"
	"log"
	"time"
	"github.com/garyburd/redigo/redis"
)

func TestPubSubHandlers(t *testing.T) {
	c, err := rediz.NewConn(RedisAddress)
	if err != nil {
		t.Fatal(err)
	}

	waiting := make(chan bool)

	c.SyncDo("FLUSHDB")

	psc := c.PubSubConn()
	if psc != c.PubSubConn() {
		t.Fail()
	}

	psc.OnMessage("test-channel", func(channel string, data []byte) {
		if message := string(data); message != "Hello, World!?" {
			t.Fail()
		}
		waiting<-true
	})

	psc.Subscribe("test-channel")

	pub, err := rediz.NewConn(RedisAddress)
	if err != nil {
		t.Fatal(err)
	}

	times := 100

	go func() {
		for i := 0; i < times; i++ {
			_, err := redis.Int64(pub.Do("PUBLISH", "test-channel", "Hello, World!?"))
			if err != nil {
				t.Fatal(err)
			}
		}
	}()

	go func() {
		time.Sleep(10 * time.Second)
		waiting<-false
	}()

	for i := 0; i < times; i++ {
		if ok := <-waiting; !ok {
			t.Fail()
		}
	}

	log.Println("Done")
}
