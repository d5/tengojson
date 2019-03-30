package tengojson

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/d5/tengo/script"
	"github.com/d5/tengo/stdlib"
)

type Builder struct {
	path       string
	processors []processor
	functions  []string
	children   []*Builder
}

func New() *Builder {
	return &Builder{}
}

func (b *Builder) At(path string, fn func(*Builder)) *Builder {
	child := &Builder{path: path}
	fn(child)
	b.children = append(b.children, child)
	return b
}

func (b *Builder) Do(src string) *Builder {
	b.functions = append(b.functions, src)
	return b
}

func (b *Builder) On(path, src string) *Builder {
	b.processors = append(b.processors, processor{
		t:    transformer,
		path: path,
		src:  src,
	})
	return b
}

type KeyPath []string

func Key(elems ...string) KeyPath {
	return elems
}

func examplecode() {
	pipeline, err := New().
		Do(`text := import("text")`).
		On("name", `text.to_lower`).
		On("products[0].name", `func(v) { if !v { return error("product name missing") } }`).
		Compile()
	if err != nil {
		panic(err)
	}

	out, err := pipeline.Run([]byte(`{}`))
	if err != nil {
		panic(err)
	}
	fmt.Println(string(out))
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
	var code bytes.Buffer

	for _, src := range b.functions {
		code.WriteString(src)
		code.WriteByte('\n')
	}

	var processorLoop []string
	for _, proc := range b.processors {
		switch proc.t {
		case validator:
			// ?
		case transformer:
			code := `
__val__ = __root__` + pathToSelector(proc.path) + `
if !is_undefined(__val__) {
	__ret__ := func() {
		__eval__ := ` + proc.src + `
		return is_callable(__eval__) ? __eval__(__val__) : __eval__
	}()
	if is_error(__ret__) { 
		return __ret__ 
	} else if !is_undefined(__ret__) { 
		__root__` + pathToSelector(proc.path) + ` = __ret__
	}
}`
			processorLoop = append(processorLoop, code)
		default:
			panic(fmt.Errorf("invalid process type: %d", proc.t))
		}
	}

	wrapper := `
json := import("json")
` + scriptOutputVarName + ` := func() {
	__root__ := json.parse(` + scriptInputVarName + `)
	if is_error(__root__) { return __root__ }

	__val__ := undefined

` + strings.Join(processorLoop, "\n") + `

	return json.stringify(__root__)
}()`

	s := script.New([]byte(wrapper))
	_ = s.Add(scriptInputVarName, "")
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

func pathToSelector(path string) string {
	return path
}

func pathClean(p string) string {
	return p
}

func pathJoin(elems ...string) string {
	return strings.Join(elems, ".")
}
