package gokibilog

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
	"time"
	"unicode/utf8"
)

type MessageLevel int

// Constants of message types according to the RFC 5424 standard
const (
	LevelDebug     MessageLevel = 10
	LevelInfo      MessageLevel = 20
	LevelNotice    MessageLevel = 30
	LevelWarning   MessageLevel = 40
	LevelError     MessageLevel = 50
	LevelCritical  MessageLevel = 60
	LevelAlert     MessageLevel = 70
	LevelEmergency MessageLevel = 80
)

// [Message] stores information that will be transmitted to the log in Kibilog.com
type Message struct {
	Message   string       `json:"message"`
	CreatedAt *int64       `json:"createdAt"`
	Level     MessageLevel `json:"level"`
	Params    any          `json:"params"`
	Partition *string      `json:"partition"`
}

// The text of the message to be saved.
func (m *Message) SetMessage(message string) {
	m.Message = message
}

// The time that the message will display. It is assumed that it indicates the time when the message occurred.
// If it is not passed, we will substitute a value equal to the time we received the request.
func (m *Message) SetCreatedAt(createdAt time.Time) {
	createdAtInt := createdAt.UTC().Unix()
	m.CreatedAt = &createdAtInt
}

// [Message] level according to RFC 5424 standard.
//
// Available value:
//
// - [LevelDebug]
//
// - [LevelInfo]
//
// - [LevelNotice]
//
// - [LevelWarning]
//
// - [LevelError]
//
// - [LevelCritical]
//
// - [LevelAlert]
//
// - [LevelEmergency]
func (m *Message) SetLevel(level MessageLevel) {
	m.Level = level
}

// If necessary, additional parameters can be registered to form an array with scalar values.
// The transmitted value must be able to be processed via "encoding/json".
// Available values: ~array, ~map, ~struct
func (m *Message) SetParams(params any) {
	m.Params = params
}

// If we need to group messages, we need to form a message partition value.
// The partition value must be strings of UUID or nil.
func (m *Message) SetPartition(partition any) error {
	switch pType := partition.(type) {
	case string:
		pStr := strings.ToLower(strings.Trim(partition.(string), " "))
		if utf8.RuneCountInString(pStr) != 36 {
			return fmt.Errorf("The length of the UUID is 36, the length of the transmitted value is %d.", utf8.RuneCountInString(pStr))
		}
		reg := regexp.MustCompile("[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}")
		if !reg.MatchString(pStr) {
			return fmt.Errorf("The passed value does not look like a UUID.")
		}
		m.Partition = &pStr
	case nil:
		m.Partition = nil
		return nil
	default:
		return fmt.Errorf("Partition must be string or nil, typed %v.", pType)
	}
	return nil
}

// NewMessage create new [Message]
func NewMessage(message string, level MessageLevel) (*Message, error) {
	message = strings.Trim(message, "\r\n\t ")
	if utf8.RuneCountInString(message) < 1 {
		return nil, errors.New("the message cannot be empty")
	}
	m := Message{
		Message:   message,
		CreatedAt: nil,
		Level:     level,
		Params:    nil,
		Partition: nil,
	}
	return &m, nil
}
