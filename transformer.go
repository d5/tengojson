package tengojson

import "github.com/d5/tengo/objects"

type Transformer interface {
	Transform(path string, object objects.Object) (objects.Object, error)
}
