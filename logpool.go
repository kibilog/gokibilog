package gokibilog

import "sync"

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

func (l LogPool) getLogId() string {
	return l.logId
}

// Create new LogPool
func NewLogPool(logId string) *LogPool {
	l := LogPool{
		logId: logId,
	}
	l.messages = []*Message{}
	return &l
}
