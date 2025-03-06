package domain

import "api/xerrors"

const (
	MIN_POINT = 0
	MAX_POINT = 7
)

type Point struct {
	X int
	Y int
}

func NewPoint(x, y int) (Point, error) {
	if x < MIN_POINT || x > MAX_POINT || y < MIN_POINT || y > MAX_POINT {
		return Point{}, xerrors.ErrBadRequest
	}
	return Point{
		X: x,
		Y: y,
	}, nil
}

func (p Point) GetX() int {
	return p.X
}

func (p Point) GetY() int {
	return p.Y
}
