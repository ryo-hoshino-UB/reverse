package dataaccess

import (
	othello "api/generated"
	"context"
)

type MoveGatewayImpl struct {
	Context context.Context
	Queries *othello.Queries
}

func (m MoveGatewayImpl) Insert(turnID int, x int32, y int32, disc int) (MoveRecord, error) {
	insertRes, err := m.Queries.CreateMove(m.Context, othello.CreateMoveParams{
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
