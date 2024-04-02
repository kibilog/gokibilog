package gokibilog

import (
	"reflect"
	"strings"
	"testing"
	"time"
)

func TestNewMessage(t *testing.T) {
	type args struct {
		message func() *Message
	}
	tests := []struct {
		name string
		args args
		want *Message
	}{
		{
			args: args{
				message: func() *Message {
					m, _ := NewMessage("test", LevelInfo)
					return m
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if reflect.TypeOf(tt.args.message()) != reflect.TypeOf(&Message{}) {
				t.Errorf("NewMessage() = %#v, want %#v", tt.args.message(), tt.want)
			}
		})
	}

	t.Run("empty message", func(t *testing.T) {
		_, err := NewMessage("", LevelInfo)
		if err == nil {
			t.Errorf("A message was created with an empty message and no error was caused")
		}
	})
}

func TestMessage_SetMessage(t *testing.T) {
	type args struct {
		message string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			args: args{
				message: "test 1",
			},
			want: "test 1",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m, _ := NewMessage(tt.args.message, LevelInfo)
			m.SetMessage(tt.args.message)
			if got := m.Message; got != tt.want {
				t.Errorf("SetMessage() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMessage_SetLevel(t *testing.T) {
	type args struct {
		level MessageLevel
	}
	tests := []struct {
		name string
		args args
		want MessageLevel
	}{
		{
			args: args{
				level: LevelDebug,
			},
			want: LevelDebug,
		},
		{
			args: args{
				level: LevelInfo,
			},
			want: LevelInfo,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m, _ := NewMessage("test", LevelDebug)
			m.SetLevel(tt.args.level)
			if got := m.Level; got != tt.want {
				t.Errorf("SetLevel() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMessage_SetCreatedAt(t *testing.T) {
	loc, _ := time.LoadLocation("Europe/Moscow")

	t.Run("nil default", func(t *testing.T) {
		m, _ := NewMessage("test", LevelDebug)

		if m.CreatedAt != nil {
			t.Errorf("SetCreatedAt() = %#v, want nil", m.CreatedAt)
		}
	})

	t.Run("UTC time", func(t *testing.T) {
		createdAt := time.Date(2024, 01, 30, 15, 16, 59, 0, loc)
		m, _ := NewMessage("test", LevelDebug)
		m.SetCreatedAt(createdAt)

		if m.CreatedAt.Format("2006-01-02T15:04:05") != createdAt.UTC().Format("2006-01-02T15:04:05") {
			t.Errorf("SetCreatedAt() = %#v, want %#v", m.CreatedAt.Format("2006-01-02T15:04:05"), createdAt.UTC().Format("2006-01-02T15:04:05"))
		}
	})

	t.Run("time not modified", func(t *testing.T) {
		createdAt := time.Date(2024, 01, 30, 15, 16, 59, 0, loc)
		m, _ := NewMessage("test", LevelDebug)
		m.SetCreatedAt(createdAt)

		if createdAt.Format("2006-01-02T15:04:05") == createdAt.UTC().Format("2006-01-02T15:04:05") {
			t.Errorf("Entered createdAt was modified!")
		}
	})
}

func TestMessage_SetPartition(t *testing.T) {

	t.Run("nil default", func(t *testing.T) {
		m, _ := NewMessage("test", LevelDebug)
		if m.Partition != nil {
			t.Errorf("SetPartition() = %#v, want nil", m.Partition)
		}
	})

	t.Run("uuid v4", func(t *testing.T) {
		uuid := "550e8400-e29b-11d4-a716-446655440000"
		m, _ := NewMessage("test", LevelDebug)
		err := m.SetPartition(uuid)
		if err != nil {
			t.Errorf("Error: %s", err.Error())
		}
		if *m.Partition != uuid {
			t.Errorf("SetPartition() = %v, want %v", *m.Partition, uuid)
		}
	})

	t.Run("uuid to lower", func(t *testing.T) {
		uuid := "1EC9414C-232A-6B00-B3C8-9E6BDECED846"
		m, _ := NewMessage("test", LevelDebug)
		err := m.SetPartition(uuid)
		if err != nil {
			t.Errorf("Error: %s", err.Error())
		}
		if *m.Partition != strings.ToLower(uuid) {
			t.Errorf("SetPartition() = %v, want %v", *m.Partition, strings.ToLower(uuid))
		}
	})

	t.Run("invalid length", func(t *testing.T) {
		uuid := "550e8400-e29b-11d4-a716-4466554400"
		m, _ := NewMessage("test", LevelDebug)
		err := m.SetPartition(uuid)
		if err == nil {
			t.Errorf("SetPartition: An invalid value of the wrong length was passed, but there is no error.")
		}
	})

	t.Run("not uuid pattern", func(t *testing.T) {
		for _, uuid := range []string{"550e8400-e29b-11d4-a7161446655440000", "550e8400-e29b-11d4-a716-44665544000J"} {
			m, _ := NewMessage("test", LevelDebug)
			err := m.SetPartition(uuid)
			if err == nil {
				t.Errorf("An incorrect UUID value with an incorrect pattern was passed, but an error was not received.")
				break
			}
		}
	})
}
