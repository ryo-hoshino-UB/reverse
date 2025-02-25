package dataaccess

import (
	othello "api/generated"
	"context"
	"log"
)

type SquareGatewayImpl struct {
	Context context.Context
	Queries *othello.Queries
}

func (s SquareGatewayImpl) FindForTurnID(turnID int) ([]SquareRecord, error) {
	squares, err := s.Queries.GetSquaresByTurnID(s.Context, int32(turnID))
	if err != nil {
		return nil, err
	}

	squareRecords := make([]SquareRecord, len(squares))
	for i, square := range squares {
		squareRecords[i] = SquareRecord{Square: square}
	}

	return squareRecords, nil
}

func (s SquareGatewayImpl) Insert(turnID int, x int32, y int32, disc int) (SquareRecord, error) {
	insertRes, err := s.Queries.CreateSquare(s.Context, othello.CreateSquareParams{
		TurnID: int32(turnID),
		X:      x,
		Y:      y,
		Disc:   int32(disc),
	})
	if err != nil {
		return SquareRecord{}, err
	}

	squareID, err := insertRes.LastInsertId()
	if err != nil {
		return SquareRecord{}, err
	}

	return SquareRecord{
		Square: othello.Square{
			ID:     int32(squareID),
			TurnID: int32(turnID),
			X:      x,
			Y:      y,
			Disc:   int32(disc),
		},
	}, nil

}

func (s SquareGatewayImpl) InsertAll(turnID int, board [][]int) error {
	for y, line := range board {
		for x, disc := range line {
			_, err := s.Insert(turnID, int32(x), int32(y), disc)

			if err != nil {
				log.Fatal(err)
				return err
			}
		}
	}

	return nil
}
