package hardware

import "github.com/indeedhat/chonker/internal/types"

func Disponse(weight float32) error {
	panic("not implemented")
}

type Report struct {
	err     error
	title   string
	message string
}

// Error implements types.Report.
func (r Report) Error() error {
	return r.err
}

// Message implements types.Report.
func (r Report) Message() string {
	return r.message
}

// Title implements types.Report.
func (r Report) Title() string {
	return r.title
}

var _ types.Report = (*Report)(nil)
