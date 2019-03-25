package tengojson

import (
	"fmt"
	"strings"

	"github.com/d5/tengo/script"
	"github.com/d5/tengo/stdlib"
)

type Builder struct {
	processors []processor
	errors     []string
}

func New() *Builder {
	return &Builder{}
}

//
func (b *Builder) Validate(path, src string) *Builder {
	b.processors = append(b.processors, processor{t: validator, path: path, src: src})
	return b
}

// same as Validate(path, "export " + src)
func (b *Builder) ValidateExpr(path, src string) *Builder {
	return b.Validate(path, "export "+src)
}

func (b *Builder) Transform(path, src string) *Builder {
	b.processors = append(b.processors, processor{t: transformer, path: path, src: src})
	return b
}

func (b *Builder) TransformExpr(path, src string) *Builder {
	return b.Transform(path, "export "+src)
}

func (b *Builder) Compile() (*Executor, error) {
	if len(b.errors) > 0 {
		return nil, fmt.Errorf(strings.Join(b.errors, "\n"))
	}

	var processorLoop []string
	for idx, proc := range b.processors {
		switch proc.t {
		case validator:
			code := fmt.Sprintf(`p = import("processor_%d")
if !p(immutable(parsed%s)) {
	return error("processor [%d]: validation failed")
}`,
				idx, pathToSelector(proc.path), idx)
			processorLoop = append(processorLoop, code)
		case transformer:
			code := fmt.Sprintf(`p = import("processor_%d")
parsed%s = p(parsed%s)`,
				idx, pathToSelector(proc.path), pathToSelector(proc.path))
			processorLoop = append(processorLoop, code)
		default:
			panic(fmt.Errorf("invalid process type: %d", proc.t))
		}
	}

	wrapper := `json := import("json")
p := undefined
__output__ := func() {
	parsed := json.parse(__input__)
	if !is_error(parsed) {
		` + strings.Join(processorLoop, "\n") + `
	}
	return json.stringify(parsed)
}()`

	s := script.New([]byte(wrapper))
	_ = s.Add("__input__", "")
	mods := stdlib.GetModuleMap(stdlib.AllModuleNames()...)
	mods.Remove("os")
	mods.Remove("fmt")
	for idx, proc := range b.processors {
		mods.AddSourceModule(fmt.Sprintf("processor_%d", idx), []byte(proc.src))
	}
	s.SetImports(mods)

	c, err := s.Compile()
	if err != nil {
		return nil, err
	}

	return &Executor{
		compiled: c,
	}, nil
}

func (b *Builder) addError(format string, args ...interface{}) {
	b.errors = append(b.errors, fmt.Sprintf(format, args...))
}

func pathToSelector(path string) string {
	return path
}
