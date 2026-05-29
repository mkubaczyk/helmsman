package app

import (
	"encoding/json"
	"reflect"
	"strconv"

	"github.com/invopop/jsonschema"
)

// truthy and falsy NullBool values
var (
	True  = NullBool{HasValue: true, Value: true}
	False = NullBool{HasValue: true, Value: false}
)

// NullBool represents a bool that may be null.
type NullBool struct {
	Value    bool
	HasValue bool // true if bool is not null
}

func (b NullBool) MarshalJSON() ([]byte, error) {
	value := b.HasValue && b.Value
	return json.Marshal(value)
}

func (b *NullBool) UnmarshalJSON(data []byte) error {
	var unmarshalledJson bool

	err := json.Unmarshal(data, &unmarshalledJson)
	if err != nil {
		return err
	}

	b.Value = unmarshalledJson
	b.HasValue = true

	return nil
}

func (b *NullBool) UnmarshalText(text []byte) error {
	str := string(text)
	if len(str) < 1 {
		return nil
	}

	value, err := strconv.ParseBool(str)
	if err != nil {
		return err
	}

	b.HasValue = true
	b.Value = value

	return nil
}

// JSONSchema instructs the jsonschema generator to represent NullBool type as boolean
func (NullBool) JSONSchema() *jsonschema.Schema {
	return &jsonschema.Schema{
		Type: "boolean",
	}
}

type MergoTransformer func(typ reflect.Type) func(dst, src reflect.Value) error

func (m MergoTransformer) Transformer(typ reflect.Type) func(dst, src reflect.Value) error {
	return m(typ)
}

// NullBoolTransformer is a custom imdario/mergo transformer for the NullBool type
func NullBoolTransformer(typ reflect.Type) func(dst, src reflect.Value) error {
	if typ != reflect.TypeOf(NullBool{}) {
		return nil
	}

	return func(dst, src reflect.Value) error {
		if src.FieldByName("HasValue").Bool() {
			dst.Set(src)
		}
		return nil
	}
}
