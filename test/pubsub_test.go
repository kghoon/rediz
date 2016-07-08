package test

import (
	"fmt"
	"log"
	"rediz"
	"testing"
	"time"
)

const ChannelName = "test_channel"

func TestPubSubHandlers(t *testing.T) {
	psAgent := rediz.RedisPubSubAgent{}
	psAgent.Activate(RedisAddress)

	err := psAgent.Subscribe(ChannelName, func(msg string) {
		log.Printf("handle msg [%s]\n", msg)
	})

	if err != nil {
		t.Fatal(err)
	}

	for i := 0; i < 5; i++ {
		msg := fmt.Sprintf("testMsg-%d", i)
		_, err := psAgent.Publish(ChannelName, msg)
		if err != nil {
			t.Fatal(err)
		}
		log.Printf("publish : %s", msg)
		time.Sleep(1 * time.Second)
	}
	psAgent.Deactivate(true)

	log.Println("Done")
}
