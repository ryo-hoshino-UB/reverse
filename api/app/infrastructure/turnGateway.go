package infrastucture

import (
	"api/domain"
	othello "api/generated"
	"context"
	"database/sql"
)

type TurnGateway struct {
	Queries *othello.Queries
}

func NewTurnGateway(q *othello.Queries) *TurnGateway {
	return &TurnGateway{
		Queries: q,
	}
}

func (t *TurnGateway) FindForGameIDAndTurnCount(ctx context.Context, gameID int, turnCount int) (TurnRecord, error) {
	turnRecord, err := t.Queries.GetTurnByGameIDAndTurnCount(ctx, othello.GetTurnByGameIDAndTurnCountParams{
		GameID:    int32(gameID),
		TurnCount: int32(turnCount),
	})
	if err != nil {
		return TurnRecord{}, err
	}
	return TurnRecord{
		Turn: turnRecord,
	}, nil
}

func (t *TurnGateway) Insert(ctx context.Context, gameID int, turnCount int, nextDisc domain.Disc) (TurnRecord, error) {
	insertRes, err := t.Queries.CreateTurn(ctx, othello.CreateTurnParams{
		GameID:    int32(gameID),
		TurnCount: int32(turnCount),
		NextDisc:  sql.NullInt32{Int32: int32(nextDisc), Valid: true},
	})
	if err != nil {
		return TurnRecord{}, err
	}

	turnID, err := insertRes.LastInsertId()
	if err != nil {
		return TurnRecord{}, err
	}

	return TurnRecord{
		Turn: othello.Turn{
			ID:        int32(turnID),
			GameID:    int32(gameID),
			TurnCount: int32(turnCount),
			NextDisc:  sql.NullInt32{Int32: int32(nextDisc), Valid: true},
		},
	}, nil
}
