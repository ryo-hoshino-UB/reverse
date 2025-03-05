package infrastucture

import othello "api/generated"

type MoveRecord struct {
	Move othello.Move
}

func NewMoveRecord(move othello.Move) MoveRecord {
	return MoveRecord{
		Move: move,
	}
}

func (m MoveRecord) GetID() int {
	return int(m.Move.ID)
}

func (m MoveRecord) GetTurnID() int {
	return int(m.Move.TurnID)
}

func (m MoveRecord) GetDisc() int {
	return int(m.Move.Disc)
}

func (m MoveRecord) GetX() int {
	return int(m.Move.X)
}

func (m MoveRecord) GetY() int {
	return int(m.Move.Y)
}
