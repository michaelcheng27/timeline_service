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
			want:    Timeline{PagingToken: "someToken"},
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
