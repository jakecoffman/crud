package crud

import (
	"encoding/json"
	"testing"
)

func TestField_ToJsonSchema(t *testing.T) {
	field := Object(map[string]Field{
		"name":       String().Required().Example("Bob"),
		"arrayMatey": Array().Items(Number()),
	})
	schema := field.ToJsonSchema()

	data, err := json.Marshal(schema)
	if err != nil {
		t.Fatal(err)
	}
	if string(data) != `{"type":"object","properties":{"arrayMatey":{"type":"array","items":{"type":"number"}},"name":{"type":"string","example":"Bob"}},"required":["name"]}` {
		t.Errorf(string(data))
	}
}
