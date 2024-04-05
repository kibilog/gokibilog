package gokibilog

import (
	"testing"
)

func Test_validatePool(t *testing.T) {
	type args struct {
		pool func() *LogPool
	}
	tests := []struct {
		name          string
		args          args
		wantErrsCount int
		messageCount  int
	}{
		{
			name: "valid",
			args: args{
				pool: func() *LogPool {
					l, _ := NewLogPool("01hggahp9skcph42wknxbckb46")
					m1, _ := NewMessage("test 1", LevelInfo)
					m2, _ := NewMessage("test 1", LevelInfo)
					l.AddMessage(m1)
					l.AddMessage(m2)
					return l
				},
			},
			wantErrsCount: 0,
			messageCount:  2,
		},
		{
			name: "valid",
			args: args{
				pool: func() *LogPool {
					l, _ := NewLogPool("01hggahp9skcph42wknxbckb46")
					m1, _ := NewMessage("test 1", LevelInfo)
					l.AddMessage(m1)
					l.AddMessage(nil)
					return l
				},
			},
			wantErrsCount: 0,
			messageCount:  1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			logPool := tt.args.pool()
			gotErrs := validatePool(logPool)

			if len(gotErrs) != tt.wantErrsCount || tt.messageCount != logPool.Len() {
				t.Errorf("validatePool(). Erros count = %v, want %v. Messages count %v, want %v.", len(gotErrs), tt.wantErrsCount, logPool.Len(), tt.messageCount)
			}

		})
	}
}
