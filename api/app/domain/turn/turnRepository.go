package turn

import (
	"api/dataaccess"
	"api/domain"
	othello "api/generated"
	"context"
	"database/sql"
	"log"
)

type TurnRepository struct {
}

func NewTurnRepository() *TurnRepository {
	return &TurnRepository{}
}

func (t *TurnRepository) Save(ctx context.Context, db *sql.DB, turn Turn) error {
	tgw := dataaccess.NewTurnGateway(othello.New(db))
	sgw := dataaccess.NewSquareGateway(othello.New(db))
	mgw := dataaccess.NewMoveGateway(othello.New(db))

	turnRecord, err := tgw.Insert(ctx, turn.GetGameID(), turn.GetTurnCount(), turn.GetNextDisc())
	if err != nil {
		log.Fatal(err)
		return err
	}

	err = sgw.InsertAll(ctx, int(turnRecord.GetID()), turn.Board.GetDiscs())
	if err != nil {
		log.Fatal(err)
		return err
	}

	_, err = mgw.Insert(ctx, int(turnRecord.GetID()), int32(turn.Move.Point.X), int32(turn.Move.Point.Y), int(turn.Move.Disc))
	if err != nil {
		log.Fatal(err)
		return err
	}
	return nil
}

func (t *TurnRepository) FindForGameIDAndTurnCount(ctx context.Context, db *sql.DB, gameID int, turnCount int) (Turn, error) {
	tgw := dataaccess.NewTurnGateway(othello.New(db))
	sgw := dataaccess.NewSquareGateway(othello.New(db))
	mgw := dataaccess.NewMoveGateway(othello.New(db))

	turnRecord, err := tgw.FindForGameIDAndTurnCount(ctx, int(gameID), turnCount)
	if err != nil {
		log.Fatal(err)
		return Turn{}, err
	}

	squaresRecord, err := sgw.FindForTurnID(ctx, int(turnRecord.GetID()))
	if err != nil {
		log.Fatal(err)
		return Turn{}, err
	}

	board := make([][]domain.Disc, 8)
	for i := range board {
		board[i] = make([]domain.Disc, 8)
	}

	for _, square := range squaresRecord {
		board[square.GetY()][square.GetX()] = domain.ToDisc(square.GetDisc())
	}

	moveRecord, err := mgw.FindForTurnID(ctx, int(turnRecord.GetID()))
	if err != nil {
		return Turn{}, err
	}

	move := domain.NewMove(domain.Disc(moveRecord.GetDisc()), domain.NewPoint(int(moveRecord.GetX()), int(moveRecord.GetY())))

	return NewTurn(
		int(turnRecord.GetGameID()),
		int(turnRecord.GetTurnCount()),
		domain.Disc(turnRecord.GetNextDisc()),
		move,
		domain.NewBoard(board),
		turnRecord.GetEndAt(),
	), nil
}
