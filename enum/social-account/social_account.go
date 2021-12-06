package social_account

type SocialAccountType int

const (
	Google SocialAccountType = iota
)

func (s SocialAccountType) String() string {
	return [...]string{"Google"}[s]
}

func ParseSocialAccountType(s SocialAccountType) int {
	switch s {
	case Google:
		return 0
	default:
		return -1
	}
}