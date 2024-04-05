package gokibilog

import (
	"fmt"
	"regexp"
	"strings"
	"sync"
)

type LogPool struct {
	mu       sync.Mutex
	logId    string
	messages []*Message
}

// AddMessage is a method for filling [LogPool] with messages
func (l *LogPool) AddMessage(message *Message) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.messages = append(l.messages, message)
}

func (l *LogPool) Len() int {
	return len(l.messages)
}

func (l *LogPool) getLogId() string {
	return l.logId
}

func (l *LogPool) removeNilMessages() {
	messagesNew := []*Message{}
	for _, message := range l.messages {
		if message != nil {
			messagesNew = append(messagesNew, message)
		}
	}
	l.messages = messagesNew
}

// Create new LogPool
func NewLogPool(logId string) (*LogPool, error) {
	logId = strings.Trim(logId, " ")
	reg := regexp.MustCompile("[0-7][0-9a-hjkmnp-tv-z]{25}")
	if !reg.MatchString(logId) {
		return nil, fmt.Errorf("The \"%s\" is not similar to the log id Kibilog.com", logId)
	}
	l := LogPool{
		logId: logId,
	}
	l.messages = []*Message{}
	return &l, nil
}
