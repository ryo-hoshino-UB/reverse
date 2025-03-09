package turn

import (
	"api/domain/model/turn"
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

func (t *TurnGateway) FindForGameIDAndTurnCount(ctx context.Context, gameID int, turnCount int) (turn.Turn, error) {
	turnRecord, err := t.Queries.GetTurnByGameIDAndTurnCount(ctx, othello.GetTurnByGameIDAndTurnCountParams{
		GameID:    int32(gameID),
		TurnCount: int32(turnCount),
	})
	if err != nil {
		return turn.Turn{}, err
	}

	return turn.Turn{
		ID : int(turnRecord.ID),
		GameID:    int(turnRecord.GameID),
		TurnCount: int(turnRecord.TurnCount),
		NextDisc:  turn.Disc(turnRecord.NextDisc.Int32),
	}, nil
}

func (t *TurnGateway) Save(ctx context.Context, inputTurn turn.Turn) (turn.Turn, error) {
	insertRes, err := t.Queries.CreateTurn(ctx, othello.CreateTurnParams{
		GameID:    int32(inputTurn.GetGameID()),
		TurnCount: int32(inputTurn.GetTurnCount()),
		NextDisc:  sql.NullInt32{Int32: int32(inputTurn.GetNextDisc()), Valid: true},
	})
	if err != nil {
		return turn.Turn{}, err
	}

	turnID, err := insertRes.LastInsertId()
	if err != nil {
		return turn.Turn{}, err
	}

	return turn.Turn{
			ID: int(turnID),
			GameID:    int(inputTurn.GetGameID()),
			TurnCount: int(inputTurn.GetTurnCount()),
			NextDisc:  inputTurn.GetNextDisc(),
	}, nil
}
