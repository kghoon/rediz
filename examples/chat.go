package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"path"
	"rediz"
	"regexp"
	"strings"
)

const (
	RedisAddress   = "local-redis:6379"
	KeyChatChannel = "test-chat"
	FmtMessage     = "msg,%s:%s" //fromUserName:message
)

var reMessage = regexp.MustCompile(fmt.Sprintf(FmtMessage, "(.+)", "(.+)"))

func main() {
	if len(os.Args) != 2 {
		log.Printf("usage: %s [USERNAME]", path.Base(os.Args[0]))
		return
	}
	chatName := os.Args[1]
	log.Printf("Welcome '%s'! Let's start chatting\n", chatName)

	psAgent := rediz.RedisPubSubAgent{}
	psAgent.Activate(RedisAddress)

	psAgent.Subscribe(KeyChatChannel, func(msg string) {
		if tokens := reMessage.FindStringSubmatch(msg); len(tokens) > 0 {
			var userName string
			if tokens[1] == chatName {
				userName = "ME"
			} else {
				userName = tokens[1]
			}
			log.Printf("[%s]: %s\n", userName, tokens[2])
		}
	})

	reader := bufio.NewReader(os.Stdin)

	for {
		text, _ := reader.ReadString('\n')

		if strings.HasPrefix(text, "exit") {
			break
		}

		message := fmt.Sprintf(FmtMessage, chatName, text)
		psAgent.Publish(KeyChatChannel, message)
	}

	psAgent.Deactivate(true)
}
