package yaccjson

import (
	"reflect"
	"testing"
)

func TestParseJson(t *testing.T) {
	type args struct {
		input string
		debug bool
	}
	tests := []struct {
		name    string
		args    args
		want    interface{}
		wantErr bool
	}{
		{``, args{`{"a":1}`, false}, map[string]interface{}{
			"a": 1,
		}, false},
		{``, args{`{"b":"str"}`, false}, map[string]interface{}{
			"b": "str",
		}, false},
		{``, args{`{"c":[1,"str"]}`, false}, map[string]interface{}{
			"c": []interface{}{1, "str"},
		}, false},
		{``, args{`{"d":{"d1":1, "d2":"str"}}`, false}, map[string]interface{}{
			"d": map[string]interface{}{"d1": 1, "d2": "str"},
		}, false},
		{``, args{`{"e":[]}`, false}, map[string]interface{}{"e": []interface{}{}}, false},
		{``, args{`{"f":{}}`, false}, map[string]interface{}{"a": 1}, false},
		{``, args{`{"g":"h w"}`, false}, map[string]interface{}{"a": 1}, false},
		{``, args{`{"h":"true"}`, false}, map[string]interface{}{"a": 1}, false},
		{``, args{`{"i":"false"}`, false}, map[string]interface{}{"a": 1}, false},
		{``, args{`{"i":"null"}`, false}, map[string]interface{}{"a": 1}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseJson(tt.args.input, tt.args.debug)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseJson() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ParseJson() = %v, want %v", got, tt.want)
			}
		})
	}
}
