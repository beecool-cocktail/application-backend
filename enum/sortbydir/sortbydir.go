package sortbydir

type SortByDir int

const (
	ASC SortByDir = iota
	DESC
)

func (s SortByDir) String() string {
	return [...]string{"ASC", "DESC"}[s]
}

func ParseSortByDirByInt(dir int) SortByDir {
	switch dir {
	case 0:
		return ASC
	case 1:
		return DESC
	default:
		return DESC
	}
}
func ParseStringBySortByDir(sortByDir SortByDir) string {
	switch sortByDir {
	case ASC:
		return "asc"
	case DESC:
		return "desc"
	default:
		return "desc"
	}
}

func MakeSortAndDir(sort, dir string) string {
	return sort + " " + dir
}