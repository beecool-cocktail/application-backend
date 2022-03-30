package httpaction

type HttpAction int

const (
	Delete HttpAction = iota
	Keep
	Edit
	Add
)

func (h HttpAction) String() string {
	return [...]string{"delete", "keep", "edit", "add"}[h]
}

func (h HttpAction) Int() int {
	switch h {
	case Delete:
		return 0
	case Keep:
		return 1
	case Edit:
		return 2
	case Add:
		return 3
	default:
		return -1
	}
}
