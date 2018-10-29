package gostruct

import (
	"encoding/json"
	"testing"
)

func TestGenStruct(t *testing.T) {
	type args struct {
		name string
		i    string
	}
	tests := []struct {
		args args
	}{
		{args{"A", `{"hello":1}`}},
	}
	for _, tt := range tests {
		var val interface{}
		json.Unmarshal([]byte(tt.args.i), &val)
		got := string(GenStruct(tt.args.name, val))
		t.Logf("// %v %v %v", tt.args.name, tt.args.i, got)
	}
}
