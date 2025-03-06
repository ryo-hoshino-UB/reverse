package infrastructure

import (
	"api/domain"
	othello "api/generated"
	"context"
	"log"
)

type SquareGateway struct {
	Queries *othello.Queries
}

func NewSquareGateway(q *othello.Queries) *SquareGateway {
	return &SquareGateway{
		Queries: q,
	}
}

func (s *SquareGateway) FindForTurnID(ctx context.Context, turnID int) ([]SquareRecord, error) {
	squares, err := s.Queries.GetSquaresByTurnID(ctx, int32(turnID))
	if err != nil {
		return nil, err
	}

	squareRecords := make([]SquareRecord, len(squares))
	for i, square := range squares {
		squareRecords[i] = SquareRecord{Square: square}
	}

	return squareRecords, nil
}

func (s *SquareGateway) Insert(ctx context.Context, turnID int, x int32, y int32, disc domain.Disc) (SquareRecord, error) {
	insertRes, err := s.Queries.CreateSquare(ctx, othello.CreateSquareParams{
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

func (s *SquareGateway) InsertAll(ctx context.Context, turnID int, board [][]domain.Disc) error {
	for y, line := range board {
		for x, disc := range line {
			_, err := s.Insert(ctx, turnID, int32(x), int32(y), disc)

			if err != nil {
				log.Println(err)
				return err
			}
		}
	}

	return nil
}
