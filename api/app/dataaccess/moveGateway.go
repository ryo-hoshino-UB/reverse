package dataaccess

import (
	othello "api/generated"
	"context"
)

type MoveGateway struct {
	q *othello.Queries
}

func NewMoveGateway(q *othello.Queries) *MoveGateway {
	return &MoveGateway{
		q: q,
	}
}

func (m MoveGateway) FindForTurnID(ctx context.Context, turnID int) (MoveRecord, error) {
	moveGetRes, err := m.q.GetMoveByTurnID(ctx, int32(turnID))
	if err != nil {
		return MoveRecord{}, err
	}

	return NewMoveRecord(moveGetRes), nil
}

func (m MoveGateway) Insert(ctx context.Context, turnID int, x int32, y int32, disc int) (MoveRecord, error) {
	insertRes, err := m.q.CreateMove(ctx, othello.CreateMoveParams{
		TurnID: int32(turnID),
		X:      int32(x),
		Y:      int32(y),
		Disc:   int32(disc),
	})
	if err != nil {
		return MoveRecord{}, err
	}

	moveID, err := insertRes.LastInsertId()
	if err != nil {
		return MoveRecord{}, err
	}

	return MoveRecord{
		Move: othello.Move{
			ID:     int32(moveID),
			TurnID: int32(turnID),
			X:      int32(x),
			Y:      int32(y),
			Disc:   int32(disc),
		},
	}, nil
}
