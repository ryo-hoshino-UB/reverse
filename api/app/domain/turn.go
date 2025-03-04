package domain

import "time"

type Turn struct {
	GameID    int
	TurnCount int
	NextDisc  Disc
	Move      Move
	Board     Board
	EndAt     time.Time
}

func NewTurn(gameID, turnCount int, nextDisc Disc, move Move, board Board, endAt time.Time) Turn {
	return Turn{
		GameID:    gameID,
		TurnCount: turnCount,
		NextDisc:  nextDisc,
		Move:      move,
		Board:     board,
		EndAt:     endAt,
	}
}

func (t Turn) PlaceNext(disc Disc, point Point) Turn {
	if disc != t.NextDisc {
		panic("invalid disc")
	}

	move := NewMove(disc, point)

	nextBoard := t.Board.Place(move)

	var nextDisc Disc
	if disc == Black {
		nextDisc = White
	} else {
		nextDisc = Black
	}

	return NewTurn(t.GameID, t.TurnCount+1, nextDisc, move, nextBoard, time.Now())
}

func (t Turn) GetGameID() int {
	return t.GameID
}

func (t Turn) GetTurnCount() int {
	return t.TurnCount
}

func (t Turn) GetNextDisc() Disc {
	return t.NextDisc
}

func (t Turn) GetMove() Move {
	return t.Move
}

func (t Turn) GetBoard() Board {
	return t.Board
}

func (t Turn) GetEndAt() time.Time {
	return t.EndAt
}
