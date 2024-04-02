package gokibilog

import (
	"reflect"
	"testing"
)

func TestGetInstance(t *testing.T) {
	tests := []struct {
		name string
		want *Kibilog
	}{
		{
			want: GetInstance(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetInstance(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetInstance() = %#v, want %#v", got, tt.want)
			}
		})
	}
}

func TestKibilog_GetSetAuthToken(t *testing.T) {
	type args struct {
		authToken string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			args: args{
				authToken: "token",
			},
			want: "token",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			GetInstance().SetAuthToken(tt.args.authToken)
			if got := GetInstance().getAuthToken(); got != tt.want {
				t.Errorf("getAuthToken() = %v, want %v", got, tt.want)
			}
		})
	}
}
