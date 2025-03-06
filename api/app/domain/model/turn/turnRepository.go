package turn

import (
	"api/domain"
	othello "api/generated"
	infrastucture "api/infrastructure"
	"api/xerrors"
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
	tgw := infrastucture.NewTurnGateway(othello.New(db))
	sgw := infrastucture.NewSquareGateway(othello.New(db))
	mgw := infrastucture.NewMoveGateway(othello.New(db))

	turnRecord, err := tgw.Insert(ctx, turn.GetGameID(), turn.GetTurnCount(), turn.GetNextDisc())
	if err != nil {
		log.Println(err)
		return err
	}

	err = sgw.InsertAll(ctx, int(turnRecord.GetID()), turn.Board.Discs)
	if err != nil {
		log.Println(err)
		return err
	}

	_, err = mgw.Insert(ctx, int(turnRecord.GetID()), int32(turn.Move.Point.X), int32(turn.Move.Point.Y), int(turn.Move.Disc))
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}

func (t *TurnRepository) FindForGameIDAndTurnCount(ctx context.Context, db *sql.DB, gameID int, turnCount int) (Turn Turn, err error) {
	defer xerrors.Wrap(&err, "FindForGameIDAndTurnCount")
	tgw := infrastucture.NewTurnGateway(othello.New(db))
	sgw := infrastucture.NewSquareGateway(othello.New(db))
	mgw := infrastucture.NewMoveGateway(othello.New(db))

	turnRecord, err := tgw.FindForGameIDAndTurnCount(ctx, int(gameID), turnCount)
	if err != nil {
		log.Println(err)
		return Turn, err
	}

	squaresRecord, err := sgw.FindForTurnID(ctx, int(turnRecord.GetID()))
	if err != nil {
		log.Println(err)
		return Turn, err
	}

	board := make([][]domain.Disc, 8)
	for i := range board {
		board[i] = make([]domain.Disc, 8)
	}

	for _, square := range squaresRecord {
		disc, err := domain.ToDisc((square.GetDisc()))
		if err != nil {
			return Turn, err
		}

		board[square.GetY()][square.GetX()] = disc
	}

	moveRecord, err := mgw.FindForTurnID(ctx, int(turnRecord.GetID()))
	if err != nil {
		return Turn, err
	}

	point, err := domain.NewPoint(int(moveRecord.GetX()), int(moveRecord.GetY()))
	if err != nil {
		return Turn, err
	}

	move := domain.NewMove(domain.Disc(moveRecord.GetDisc()), point)

	return NewTurn(
		int(turnRecord.GetGameID()),
		int(turnRecord.GetTurnCount()),
		domain.Disc(turnRecord.GetNextDisc()),
		move,
		domain.NewBoard(board),
		turnRecord.GetEndAt(),
	), nil
}
