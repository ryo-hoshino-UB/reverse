package domain

import "api/xerrors"

type Disc int

const (
	Empty Disc = iota
	Black
	White
	Wall
)

func ToDisc(v int) (Disc, error) {
	switch v {
	case 0:
		return Empty, nil
	case 1:
		return Black, nil
	case 2:
		return White, nil
	case 3:
		return Wall, nil
	default:
		return 0, xerrors.ErrBadRequest
	}
}

func IsOppositeDisc(d1, d2 Disc) bool {
	return ((d1 == Black && d2 == White) ||
		(d1 == White && d2 == Black))
}
