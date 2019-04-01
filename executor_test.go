package tengojson_test

import (
	"encoding/json"
	"testing"

	"github.com/d5/tengo/script"
	"github.com/d5/tengo/stdlib"
	"github.com/d5/tengojson"
	"github.com/stretchr/testify/assert"
)

func TestExecutor(t *testing.T) {
	input := []byte(`
{
	"age": 36,
	"name": "Tengo",
	"address": {
		"city": "Los Angeles",
		"country": "USA",
		"zip": "90005"
	},
	"male": true, 
	"tags": ["tag1", "tag2", "tag3"]
}
`)

	c, err := tengojson.New().
		Do(`x := import("enum")`).
		On(".age", `string`).
		On(".address.zip", `func(v) { if len(v) != 5 { return error("wrong zip code") } }`).
		On(".tags", `func(v) { if !x.all(v, x.value) { return error("invalid tag") } } `).
		Compile()
	if !assert.NoError(t, err) {
		return
	}
	output, err := c.Run(input)
	assertEqualJSON(t, output, ".age", "36")
}

func assertEqualJSON(t *testing.T, b []byte, key string, expected interface{}) bool {
	v := make(map[string]interface{})
	err := json.Unmarshal(b, &v)
	if !assert.NoError(t, err) {
		return false
	}
	s := script.New([]byte(`
json := import("json")
parsed := json.parse(input)
output := is_error(parsed) ? parsed : parsed` + key))
	s.SetImports(stdlib.GetModuleMap("json"))
	_ = s.Add("input", b)
	c, err := s.Run()
	if !assert.NoError(t, err) {
		return false
	}
	return assert.Equal(t, expected, c.Get("output").Value())
}
