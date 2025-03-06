package turn

import (
	"api/domain"
	gameresult "api/domain/model/gameResult"
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

func (t Turn) PlaceNext(disc domain.Disc, point domain.Point) (Turn, error) {
	if disc != t.NextDisc {
		panic("invalid disc")
	}

	move := domain.NewMove(disc, point)

	nextBoard, err := t.Board.Place(move)
	if err != nil {
		return Turn{}, err
	}

	nextDisc := t.decideNextDisc(nextBoard, disc)

	return NewTurn(t.GameID, t.TurnCount+1, nextDisc, move, nextBoard, time.Now()), nil
}

func (t Turn) decideNextDisc(board domain.Board, prevDisc domain.Disc) domain.Disc {
	existBlackValidMove := board.ExistValidMove(domain.Black)
	existWhiteValidMove := board.ExistValidMove(domain.White)

	if existBlackValidMove && existWhiteValidMove {
		if prevDisc == domain.Black {
			return domain.White
		} else {
			return domain.Black
		}
	}

	if existBlackValidMove && !existWhiteValidMove {
		return domain.Black
	}

	if !existBlackValidMove && existWhiteValidMove {
		return domain.White
	}

	if !existBlackValidMove && !existWhiteValidMove {
		return 0
	}

	panic("unreachable")
}

func (t Turn) WinnerDisc() gameresult.WinnerDisc {
	blackDiscCount := t.Board.CountDiscs(domain.Black)
	whiteDiscCount := t.Board.CountDiscs(domain.White)

	if blackDiscCount > whiteDiscCount {
		return gameresult.BlackWin
	}

	if blackDiscCount < whiteDiscCount {
		return gameresult.WhiteWin
	}

	return gameresult.Draw
}

func (t Turn) GameEnded() bool {
	return t.NextDisc == domain.Empty
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
