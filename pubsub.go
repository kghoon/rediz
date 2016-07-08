package rediz

import (
	"log"
	"sync"

	"github.com/garyburd/redigo/redis"
)

type MessageHandler func(msg string)

type RedisPubSubAgent struct {
	wg              sync.WaitGroup
	subConn         *ResourceConn
	pubConn         *ResourceConn
	quit            bool
	messageHandlers map[string]MessageHandler
}

func (psa *RedisPubSubAgent) Activate(redisAddress string) {
	subConn, err := NewConn(redisAddress)
	if err != nil {
		panic(err)
	}
	pubConn, err := NewConn(redisAddress)
	if err != nil {
		panic(err)
	}

	psa.wg = sync.WaitGroup{}
	psa.subConn = &subConn
	psa.pubConn = &pubConn
	psa.quit = false
	psa.messageHandlers = make(map[string]MessageHandler)

	go receiveWorker(psa)
}

func (psa *RedisPubSubAgent) Deactivate(wait bool) {
	channels := make([]string, 0, len(psa.messageHandlers))
	for k := range psa.messageHandlers {
		channels = append(channels, k)
	}
	psa.subConn.Unsubscribe(channels)

	psa.messageHandlers = nil
	psa.quit = true

	if wait {
		psa.wg.Wait()
	}
}

func (psa *RedisPubSubAgent) Subscribe(channel string, msgHandle MessageHandler) error {
	err := psa.subConn.Subscribe(channel)
	if err == nil {
		psa.messageHandlers[channel] = msgHandle
	}
	return err
}

func (psa *RedisPubSubAgent) Publish(channel string, msg string) (reply interface{}, err error) {
	return psa.pubConn.Conn.Do("PUBLISH", channel, msg)
}

func receiveWorker(psa *RedisPubSubAgent) {
	psa.wg.Add(1)
	defer psa.wg.Done()
	defer psa.subConn.Close()
	defer psa.pubConn.Close()
	log.Println("Start Receive worker!!")

	for {
		if psa.quit {
			log.Println("Shutdown Receive worker!!")
			return
		} else {

			switch v := psa.subConn.Receive().(type) {
			case redis.Message:
				if h, exist := psa.messageHandlers[v.Channel]; exist {
					h(string(v.Data))
				}
			}
		}
	}
}
