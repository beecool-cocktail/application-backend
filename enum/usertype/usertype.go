package usertype

type UserType int

const (
	Test UserType = iota
	Normal
)

func (h UserType) Int() int {
	switch h {
	case Test:
		return 0
	case Normal:
		return 1
	default:
		return 1
	}
}
