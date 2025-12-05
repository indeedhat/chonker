package types

type Report interface {
	Error() error
	Message() string
	Title() string
}
