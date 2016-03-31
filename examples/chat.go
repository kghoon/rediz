package main

import (
	"rediz"
	"regexp"
	"fmt"
	"log"
	"github.com/garyburd/redigo/redis"
	"bufio"
	"os"
	"strconv"
	"strings"
)

const (
	RedisAddress = "192.168.99.100:6379"
	KeyChatChannel = "test-chat"
	KeyUserId = "user-id"
	FmtMessage = "msg,%s:%s"
)

var reMessage = regexp.MustCompile(fmt.Sprintf(FmtMessage, "([\\d]+)", "(.+)"))

func main() {
	pub, sub := NewConnPair()
	defer pub.Close()
	defer sub.Close()

	userId := GetUserId(pub)

	psc := sub.PubSubConn()

	psc.OnMessage(KeyChatChannel, func(channel string, data []byte) {
		if tokens := reMessage.FindStringSubmatch(string(data)); len(tokens) > 0 {
			if id, err := strconv.ParseInt(tokens[1], 10, 64); err != nil {
				log.Println("Error:", err)
			} else if id == userId {
				tokens[1] = "ME"
			}
			log.Printf("[%s]: %s\n", tokens[1], tokens[2])
		}
	})

	psc.Subscribe(KeyChatChannel)
	defer psc.Unsubscribe(KeyChatChannel)

	log.Printf("Welcome! Your ID is %d \n", userId)
	reader := bufio.NewReader(os.Stdin)

	for {
		text, _ := reader.ReadString('\n')

		if strings.HasPrefix(text, "exit") {
			break
		}

		message := fmt.Sprintf(FmtMessage, strconv.FormatInt(userId, 10), text)
		pub.SyncDo("PUBLISH", KeyChatChannel, message)
	}
}

func NewConnPair() (*rediz.ResourceConn, *rediz.ResourceConn) {
	pub, err := rediz.NewConn(RedisAddress)
	if err != nil {
		panic(err)
	}

	sub, err := rediz.NewConn(RedisAddress)
	if err != nil {
		panic(err)
	}

	return &pub, &sub
}

func GetUserId(c *rediz.ResourceConn) int64 {
	id, err := redis.Int64(c.SyncDo("INCR", KeyUserId))
	if err != nil {
		panic(err)
	}
	return id
}
