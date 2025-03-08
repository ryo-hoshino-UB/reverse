package infrastructure

import (
	othello "api/generated"
	"api/xerrors"
	"context"
	"time"
)

type GameResultGateway struct {
	q *othello.Queries
}

func NewGameResultGateway(q *othello.Queries) *GameResultGateway {
	return &GameResultGateway{
		q: q,
	}
}

func (g *GameResultGateway) SelectGameResult(ctx context.Context, gameID int) (GameResultRecord, error) {
	gameResult, err := g.q.GetGameResultByGameID(ctx, int32(gameID))
	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			return GameResultRecord{}, xerrors.ErrNotFound
		}
		return GameResultRecord{}, err
	}

	return GameResultRecord{
		ID:         int(gameResult.ID),
		GameID:     int(gameResult.GameID),
		WinnerDisc: int(gameResult.WinnerDisc),
		EndAt:      gameResult.EndAt,
	}, nil
}

func (g *GameResultGateway) Insert(ctx context.Context, gameID int, winnerDisc int, endAt time.Time) (GameResultRecord, error) {
	insertRes, err := g.q.CreateGameResult(ctx, othello.CreateGameResultParams{
		GameID:     int32(gameID),
		WinnerDisc: int32(winnerDisc),
	})
	if err != nil {
		return GameResultRecord{}, err
	}

	gameResultID, err := insertRes.LastInsertId()
	if err != nil {
		return GameResultRecord{}, err
	}

	return GameResultRecord{
		ID:         int(gameResultID),
		GameID:     gameID,
		WinnerDisc: winnerDisc,
		EndAt:      endAt,
	}, nil
}
