package cockarticletype

type CocktailArticleType int

const (
	Draft CocktailArticleType = iota
	Formal
)

func (s CocktailArticleType) String() string {
	return [...]string{"draft", "normal"}[s]
}

func (s CocktailArticleType) Int() int {
	switch s {
	case Draft:
		return 0
	case Formal:
		return 1
	default:
		return -1
	}
}
