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

	code.WriteString(`
__json__ := import("json")
` + scriptOutputVarName + ` := func() {
	__root__ := __json__.parse(` + scriptInputVarName + `)
	if is_error(__root__) { return __root__ }

	__val__ := undefined

` + strings.Join(processorLoop, "\n") + `

	return __json__.stringify(__root__)
}()`)

	s := script.New(code.Bytes())
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

func (b *Builder) Run(input []byte) (output []byte, err error) {
	compiled, err := b.Compile()
	if err != nil {
		return nil, err
	}

	return compiled.Run(input)
}

func pathToSelector(path string) string {
	if path == "." { return ""}

	return path
}