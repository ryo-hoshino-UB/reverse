package turn

import (
	"api/domain/model/gameresult"
	"time"
)

type Turn struct {
	ID 	  int
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

func (t Turn) PlaceNext(disc Disc, point Point) (Turn, error) {
	if disc != t.NextDisc {
		panic("invalid disc")
	}

	move := NewMove(disc, point)

	nextBoard, err := t.Board.Place(move)
	if err != nil {
		return Turn{}, err
	}

	nextDisc := t.decideNextDisc(nextBoard, disc)

	return NewTurn(t.GameID, t.TurnCount+1, nextDisc, move, nextBoard, time.Now()), nil
}

func (t Turn) decideNextDisc(board Board, prevDisc Disc) Disc {
	existBlackValidMove := board.ExistValidMove(Black)
	existWhiteValidMove := board.ExistValidMove(White)

	if existBlackValidMove && existWhiteValidMove {
		if prevDisc == Black {
			return White
		} else {
			return Black
		}
	}

	if existBlackValidMove && !existWhiteValidMove {
		return Black
	}

	if !existBlackValidMove && existWhiteValidMove {
		return White
	}

	if !existBlackValidMove && !existWhiteValidMove {
		return 0
	}

	panic("unreachable")
}

func (t Turn) WinnerDisc() gameresult.WinnerDisc {
	blackDiscCount := t.Board.CountDiscs(Black)
	whiteDiscCount := t.Board.CountDiscs(White)

	if blackDiscCount > whiteDiscCount {
		return gameresult.BlackWin
	}

	if blackDiscCount < whiteDiscCount {
		return gameresult.WhiteWin
	}

	return gameresult.Draw
}

func (t Turn) GameEnded() bool {
	return t.NextDisc == Empty
}

func (t Turn) GetID() int {
	return t.ID
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

func NewFirstTurn(gameID int, endAt time.Time) Turn {
	return NewTurn(gameID, 0, Black, Move{}, NewInitialBoard(), endAt)
}
