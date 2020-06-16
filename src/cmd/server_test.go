package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"reflect"
	"testing"
)

func TestServe(t *testing.T) {
	type args struct {
		request TimelineRequest
	}
	tests := []struct {
		name    string
		args    args
		want    Timeline
		wantErr bool
	}{
		{
			name:    "dummy",
			args:    args{TimelineRequest{}},
			want:    Timeline{PagingToken: nil},
			wantErr: false,
		},
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Serve(tt.args.request)
			var buf bytes.Buffer

			body, err := json.Marshal(got)
			json.HTMLEscape(&buf, body)
			t.Log(fmt.Sprint("this is body s", buf.String()))
			if (err != nil) != tt.wantErr {
				t.Errorf("Serve() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Serve() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_getMoments(t *testing.T) {
	type args struct {
		PagingToken *string
	}
	tests := []struct {
		name      string
		args      args
		want      *[]Moment
		wantToken *string
		wantErr   bool
	}{
		// TODO: Add test cases.
		{
			name: "test",
			args: args{
				PagingToken: nil,
			},
			want:      &[]Moment{},
			wantToken: nil,
			wantErr:   false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, nextToken, err := getMoments(tt.args.PagingToken)
			t.Errorf("nextToken = %v", *nextToken)
			if (err != nil) != tt.wantErr {
				t.Errorf("getMoments() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, *tt.want) {
				t.Errorf("getMoments() = %v, want %v", got, *tt.want)
			}
			got, nextToken, err = getMoments(nextToken)
			t.Errorf(" got = %v, nextToken = %v, err= %v", *got, *nextToken, err)
			got, nextToken, err = getMoments(nextToken)
			t.Errorf(" got = %v, nextToken = %v, err= %v", got, nextToken, nextToken)
		})
	}
}
