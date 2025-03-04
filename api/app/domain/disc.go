package domain

type Disc int

const (
	Empty Disc = iota
	Black
	White
)

func ToDisc(v int) Disc {
	switch v {
	case 0:
		return Empty
	case 1:
		return Black
	case 2:
		return White
	default:
		return Empty
	}
}
