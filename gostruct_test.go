package gostruct

import (
	"encoding/json"
	"strings"
	"testing"
)

func TestGenStruct(t *testing.T) {
	type args struct {
		name string
		i    string
	}
	tests := []struct {
		args args
		want string
	}{
		{args{"A", `{"hello":1}`}, ``},
	}
	for _, tt := range tests {
		var val interface{}
		json.Unmarshal([]byte(tt.args.i), &val)
		if got := string(GenStruct(tt.args.name, val)); strings.TrimSpace(got) != strings.TrimSpace(tt.want) {
			t.Errorf("GenStruct() = %v, want %v", got, tt.want)
		}
	}
}
