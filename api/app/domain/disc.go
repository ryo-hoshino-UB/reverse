package domain

type Disc int

const (
	Empty Disc = iota
	Black
	White
	Wall
)

func ToDisc(v int) Disc {
	switch v {
	case 0:
		return Empty
	case 1:
		return Black
	case 2:
		return White
	case 3:
		return Wall
	default:
		return Empty
	}
}

func IsOppositeDisc(d1, d2 Disc) bool {
	return ((d1 == Black && d2 == White) ||
		(d1 == White && d2 == Black))
}
