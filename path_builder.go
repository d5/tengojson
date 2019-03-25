package tengojson

type pathBuilder struct {
	parent     *pathBuilder
	path       string
	processors []processor
}

func (b *pathBuilder) Path(path string) *pathBuilder {
	return &pathBuilder{
		parent: b,
		path:   path,
	}
}
