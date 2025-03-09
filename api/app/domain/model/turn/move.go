package turn

type Move struct {
	Disc  Disc
	Point Point
}

func NewMove(disc Disc, point Point) Move {
	return Move{
		Disc:  disc,
		Point: point,
	}
}

func (m Move) GetDisc(move Move) Disc {
	return move.Disc
}

func (m Move) GetPoint(move Move) Point {
	return move.Point
}
