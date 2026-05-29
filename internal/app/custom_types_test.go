package app

import (
	"bytes"
	"encoding/json"
	"reflect"
	"testing"
)

func TestNullBool_MarshalJSON(t *testing.T) {
	tests := []struct {
		name    string
		value   NullBool
		want    []byte
		wantErr bool
	}{
		{
			name:    "should be false",
			want:    []byte(`false`),
			wantErr: false,
		},
		{
			name:    "should be true",
			want:    []byte(`true`),
			value:   NullBool{HasValue: true, Value: true},
			wantErr: false,
		},
		{
			name:    "should be false when HasValue is false",
			want:    []byte(`false`),
			value:   NullBool{HasValue: false, Value: true},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.value.MarshalJSON()
			if (err != nil) != tt.wantErr {
				t.Errorf("MarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("MarshalJSON() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNullBool_UnmarshalJSON(t *testing.T) {
	type output struct {
		Value NullBool `json:"value"`
	}
	tests := []struct {
		name    string
		data    []byte
		want    output
		wantErr bool
	}{
		{
			name: "should have value set to false",
			data: []byte(`{"value": false}`),
			want: output{NullBool{HasValue: true, Value: false}},
		},
		{
			name: "should have value set to true",
			data: []byte(`{"value": true}`),
			want: output{NullBool{HasValue: true, Value: true}},
		},
		{
			name: "should have value unset",
			data: []byte("{}"),
			want: output{NullBool{HasValue: false, Value: false}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var got output
			if err := json.NewDecoder(bytes.NewReader(tt.data)).Decode(&got); (err != nil) != tt.wantErr {
				t.Errorf("UnmarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("UnmarshalJSON() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNullBool_UnmarshalText(t *testing.T) {
	tests := []struct {
		name    string
		text    []byte
		want    NullBool
		wantErr bool
	}{
		{
			name: "should have the value set to false",
			text: []byte("false"),
			want: NullBool{HasValue: true, Value: false},
		},
		{
			name: "should have the value set to true",
			text: []byte("false"),
			want: NullBool{HasValue: true, Value: false},
		},
		{
			name: "should have the value unset",
			text: []byte(""),
			want: NullBool{HasValue: false, Value: false},
		},
		{
			name:    "should return an error on wrong input",
			text:    []byte("wrong_input"),
			wantErr: true,
			want:    NullBool{HasValue: false, Value: false},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var got NullBool
			if err := got.UnmarshalText(tt.text); (err != nil) != tt.wantErr {
				t.Errorf("UnmarshalText() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("UnmarshalText() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNullBoolTransformer(t *testing.T) {
	type args struct {
		dst NullBool
		src NullBool
	}
	tests := []struct {
		name string
		args args
		want NullBool
	}{
		{
			name: "should overwrite true to false when the dst has the value",
			args: args{
				dst: NullBool{HasValue: true, Value: true},
				src: NullBool{HasValue: true, Value: false},
			},
			want: NullBool{HasValue: true, Value: false},
		},
		{
			name: "shouldn't overwrite when the value is unset",
			args: args{
				dst: NullBool{HasValue: true, Value: true},
				src: NullBool{HasValue: false, Value: false},
			},
			want: NullBool{HasValue: true, Value: true},
		},
		{
			name: "shouldn overwrite when the value is set and equal true",
			args: args{
				dst: NullBool{HasValue: true, Value: false},
				src: NullBool{HasValue: true, Value: true},
			},
			want: NullBool{HasValue: true, Value: true},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dst := tt.args.dst
			src := tt.args.src

			transformer := NullBoolTransformer(reflect.TypeOf(NullBool{}))

			transformer(reflect.ValueOf(&dst).Elem(), reflect.ValueOf(src))

			if !reflect.DeepEqual(dst, tt.want) {
				t.Errorf("NullBoolTransformer() = %v, want %v", dst, tt.want)
			}
		})
	}
}
