package tengojson_test

import (
	"testing"

	"github.com/d5/tengo/assert"
	"github.com/d5/tengojson"
)

func TestExecutor(t *testing.T) {
	c, err := tengojson.New().
		Validate(".age", `export func(v) { return is_undefined(v) || is_int(v) || is_float(v) }`).
		Validate(".name", `export is_string`).
		Transform(".is_male", `export func(v) { if v == "male" { return true } }`).
		Compile()

	assert.NoError(t, err)
	out, err := c.Run([]byte(`{"name": "tengo"}`))
	assert.NoError(t, err)
	assert.Equal(t, `{"name":"tengo"}`, string(out))

	out, err = c.Run([]byte(`{"name2": "tengo"}`))
	assert.Error(t, err)
}
