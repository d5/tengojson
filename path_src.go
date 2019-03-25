package tengojson

type processorType byte

const (
	transformer processorType = iota
	validator
)

type processor struct {
	t    processorType
	path string
	src  string
}
