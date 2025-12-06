package hardware

type scale struct {
}

func newScale() scale {
	return scale{}
}

func (s scale) Await(weight float32) <-chan struct{} {
	panic("not implemented")
}

func (s scale) Close() {
	panic("not implemented")
}
