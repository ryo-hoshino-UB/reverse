package turn

import (
	"api/domain"
	"time"
)

type Turn struct {
	GameID    int
	TurnCount int
	NextDisc  domain.Disc
	Move      domain.Move
	Board     domain.Board
	EndAt     time.Time
}

func NewTurn(gameID, turnCount int, nextDisc domain.Disc, move domain.Move, board domain.Board, endAt time.Time) Turn {
	return Turn{
		GameID:    gameID,
		TurnCount: turnCount,
		NextDisc:  nextDisc,
		Move:      move,
		Board:     board,
		EndAt:     endAt,
	}
}

func (t Turn) PlaceNext(disc domain.Disc, point domain.Point) Turn {
	if disc != t.NextDisc {
		panic("invalid disc")
	}

	move := domain.NewMove(disc, point)

	nextBoard := t.Board.Place(move)

	var nextDisc domain.Disc
	if disc == domain.Black {
		nextDisc = domain.White
	} else {
		nextDisc = domain.Black
	}

	return NewTurn(t.GameID, t.TurnCount+1, nextDisc, move, nextBoard, time.Now())
}

func (t Turn) GetGameID() int {
	return t.GameID
}

func (t Turn) GetTurnCount() int {
	return t.TurnCount
}

func (t Turn) GetNextDisc() domain.Disc {
	return t.NextDisc
}

func (t Turn) GetMove() domain.Move {
	return t.Move
}

func (t Turn) GetBoard() domain.Board {
	return t.Board
}

func (t Turn) GetEndAt() time.Time {
	return t.EndAt
}

func NewFirstTurn(gameID int, endAt time.Time) Turn {
	return NewTurn(gameID, 0, domain.Black, domain.Move{}, domain.NewInitialBoard(), endAt)
}
