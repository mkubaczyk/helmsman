// go:build exclude

package main

import (
	"encoding/json"
	"os"

	"github.com/invopop/jsonschema"
	"github.com/mkubaczyk/helmsman/internal/app"
)

func main() {
	r := new(jsonschema.Reflector)
	r.AllowAdditionalProperties = true
	if err := r.AddGoComments("github.com/mkubaczyk/helmsman", "./internal/app"); err != nil {
		panic(err)
	}
	s := r.Reflect(&app.State{})
	data, _ := json.MarshalIndent(s, "", "  ")
	os.WriteFile("schema.json", data, 0o644)
}
