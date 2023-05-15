package main

import (
	"reflect"
	"testing"
)

func Test_getEvents(t *testing.T) {
	type args struct {
		content string
	}
	tests := []struct {
		name    string
		args    args
		want    []string
		wantErr bool
	}{
		{
			name: "single event",
			args: args{
				content: `{"version": "1.0", "timestamp": "2023-05-10T13:42:07.720Z", "message": {}}`,
			},
			want: []string{
				`{"version": "1.0", "timestamp": "2023-05-10T13:42:07.720Z", "message": {}}`,
			},
			wantErr: false,
		},
		{
			name: "multiple events",
			args: args{
				content: `{"version": "1.0", "timestamp": "2023-05-10T13:42:07.720Z", "message": {}}{"version": "1.0", "timestamp": "2023-05-11T13:42:07.720Z", "message": {}}`,
			},
			want: []string{
				`{"version": "1.0", "timestamp": "2023-05-10T13:42:07.720Z", "message": {}}`,
				`{"version": "1.0", "timestamp": "2023-05-11T13:42:07.720Z", "message": {}}`,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := getEvents(tt.args.content)
			if (err != nil) != tt.wantErr {
				t.Errorf("getEvents() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getEvents() got = %v, want %v", got, tt.want)
			}
		})
	}
}
