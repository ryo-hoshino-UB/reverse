package infrastructure

import othello "api/generated"

type SquareRecord struct {
	Square othello.Square
}

func (s SquareRecord) GetID() int {
	return int(s.Square.ID)
}

func (s SquareRecord) GetTurnD() int {
	return int(s.Square.TurnID)
}

func (s SquareRecord) GetX() int {
	return int(s.Square.X)
}

func (s SquareRecord) GetY() int {
	return int(s.Square.Y)
}

func (s SquareRecord) GetDisc() int {
	return int(s.Square.Disc)
}
