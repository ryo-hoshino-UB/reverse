package domain

type Point struct {
	X int
	Y int
}

func NewPoint(x, y int) Point {
	return Point{
		X: x,
		Y: y,
	}
}

func (p Point) GetX() int {
	return p.X
}

func (p Point) GetY() int {
	return p.Y
}
