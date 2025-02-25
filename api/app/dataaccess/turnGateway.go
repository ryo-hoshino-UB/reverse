package dataaccess

import (
	othello "api/generated"
	"context"
	"database/sql"
)

type TurnGateway struct {
	Context context.Context
	Queries *othello.Queries
}

func (t TurnGateway) FindForGameIDAndTurnCount(gameID int, turnCount int) (TurnRecord, error) {
	turnRecord, err := t.Queries.GetTurnByGameIDAndTurnCount(t.Context, othello.GetTurnByGameIDAndTurnCountParams{
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

func (t TurnGateway) Insert(gameID int, turnCount int, nextDisc int) (TurnRecord, error) {
	insertRes, err := t.Queries.CreateTurn(t.Context, othello.CreateTurnParams{
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
