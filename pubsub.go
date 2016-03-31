package rediz

import (
	"github.com/garyburd/redigo/redis"
)

type MessageHandler func(channel string, data []byte)

type SubscriptionHandler func(channel, kind string, count int)

type PubSubConn struct {
	*redis.PubSubConn
	parent               *ResourceConn
	messageHandlers      map[string]MessageHandler
	subscriptionHandlers map[string]SubscriptionHandler
}

func (c *ResourceConn) PubSubConn() *PubSubConn {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if c.psc == nil {

		c.psc = &PubSubConn{
			PubSubConn: &redis.PubSubConn{c.Conn},
			parent: c,
			messageHandlers: make(map[string]MessageHandler),
			subscriptionHandlers: make(map[string]SubscriptionHandler),
		}

		go c.psc.receiver()
	}

	return c.psc
}

func (c *PubSubConn) OnMessage(channel string, cb MessageHandler) {
	c.messageHandlers[channel] = cb
}

func (c *PubSubConn) OnSubscription(channel string, cb SubscriptionHandler) {
	c.subscriptionHandlers[channel] = cb
}

func (psc *PubSubConn) receiver() {
	for {
		switch v := psc.Receive().(type) {
		case redis.Message:
			psc.onMessage(v.Channel, v.Data)
		case redis.Subscription:
			psc.onSubscription(v.Channel, v.Kind, v.Count)
		}
	}
}

func (psc *PubSubConn) onMessage(channel string, data []byte) {
	if h, exist := psc.messageHandlers[channel]; exist {
		h(channel, data)
	}
}

func (psc *PubSubConn) onSubscription(channel, kind string, count int) {
	if h, exist := psc.subscriptionHandlers[channel]; exist {
		h(channel, kind, count)
	}
}
