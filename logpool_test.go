package gokibilog

import (
	"reflect"
	"testing"
)

func TestNewLogPool(t *testing.T) {
	type args struct {
		logId string
	}
	tests := []struct {
		name string
		args args
		want func() *LogPool
	}{
		{
			args: args{
				logId: "01hggahp9skcph42wknxbckb46",
			},
			want: func() *LogPool {
				l, _ := NewLogPool("01hggahp9skcph42wknxbckb46")
				return l
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l, _ := NewLogPool(tt.args.logId)
			if got := l; !reflect.DeepEqual(got, tt.want()) {
				t.Errorf("NewLogPool() = %v, want %v", got, tt.want())
			}
		})
	}

	t.Run("wrong logId", func(t *testing.T) {
		_, err := NewLogPool("01hggahp9skcph42wknxbckb")
		if err == nil {
			t.Errorf("An incorrect LogID was passed to the LogPool and this did not cause an error")
		}
	})
}

func TestLogPool_AddMessage(t *testing.T) {
	type args struct {
		logId    string
		messages func() []*Message
	}
	tests := []struct {
		name string
		args args
	}{
		{
			args: args{
				logId: "01hggahp9skcph42wknxbckb46",
				messages: func() []*Message {
					m1, _ := NewMessage("test", LevelInfo)
					m2, _ := NewMessage("test", LevelInfo)
					return []*Message{
						m1,
						m2,
					}
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l, _ := NewLogPool(tt.args.logId)
			for _, m := range tt.args.messages() {
				l.AddMessage(m)
			}
			if len(tt.args.messages()) != l.Len() {
				t.Errorf("AddMessage(): The number of messages differs. Added: %d, total: %d", l.Len(), len(tt.args.messages()))
			}
		})
	}
}

func TestLogPool_getLogId(t *testing.T) {
	type args struct {
		logId string
	}
	tests := []struct {
		name string
		args args
	}{
		{
			args: args{
				logId: "01hggahp9skcph42wknxbckb46",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l, _ := NewLogPool(tt.args.logId)
			if l.getLogId() != tt.args.logId {
				t.Errorf("getLogId() = %v, want %v", l.getLogId(), tt.args.logId)
			}
		})
	}
}

func TestLogPool_removeNilMessages(t *testing.T) {
	type args struct {
		logId    string
		messages func() []*Message
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			args: args{
				logId: "01hggahp9skcph42wknxbckb46",
				messages: func() []*Message {
					a := []*Message{}
					m1, _ := NewMessage("test", LevelInfo)
					a = append(a, m1)
					m2, _ := NewMessage("test", LevelInfo)
					a = append(a, m2)
					return a
				},
			},
			want: 2,
		},
		{
			args: args{
				logId: "01hggahp9skcph42wknxbckb46",
				messages: func() []*Message {
					a := []*Message{}
					m1, _ := NewMessage("test", LevelInfo)
					a = append(a, m1)
					a = append(a, nil)
					return a
				},
			},
			want: 1,
		},
		{
			args: args{
				logId: "01hggahp9skcph42wknxbckb46",
				messages: func() []*Message {
					a := []*Message{}
					a = append(a, nil)
					return a
				},
			},
			want: 0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l, _ := NewLogPool(tt.args.logId)
			l.messages = tt.args.messages()
			l.removeNilMessages()
			if l.Len() != tt.want {
				t.Errorf("removeNilMessages() and Len() after = %v, want %v", l.Len(), tt.args.logId)
			}
		})
	}
}
