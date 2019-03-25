package tengojson

import "github.com/d5/tengo/objects"

type Validator interface {
	Validate(path string, object objects.Object) bool
}
