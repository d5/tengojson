package tengojson

import "github.com/d5/tengo/script"

type Executor struct {
	compiled *script.Compiled
}

func (e *Executor) Run(input []byte) ([]byte, error) {
	_ = e.compiled.Set("__input__", input)

	if err := e.compiled.Run(); err != nil {
		return nil, err
	}

	output := e.compiled.Get("__output__")
	if err := output.Error(); err != nil {
		return nil, err
	}

	return output.Bytes(), nil
}
