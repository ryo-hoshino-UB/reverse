package turn

import (
	"api/domain/model/turn"
	othello "api/generated"
	"api/xerrors"
	"context"
	"database/sql"
	"log"
)


type TurnMySQLRepositoryImpl struct {
	turnRepository turn.TurnRepository
}

func NewTurnMySQLRepositoryImpl() *TurnMySQLRepositoryImpl {
	return &TurnMySQLRepositoryImpl{}
}

func (t *TurnMySQLRepositoryImpl) Save(ctx context.Context, db *sql.DB, turn turn.Turn) error {
	tgw := NewTurnGateway(othello.New(db))
	sgw := NewSquareGateway(othello.New(db))
	mgw := NewMoveGateway(othello.New(db))

	turnRecord, err := tgw.Save(ctx, turn)
	if err != nil {
		log.Println(err)
		return err
	}

	err = sgw.InsertAll(ctx, turnRecord.GetID(), turn.Board.Discs)
	if err != nil {
		log.Println(err)
		return err
	}

	_, err = mgw.Insert(ctx, turnRecord.GetID(), int32(turn.Move.Point.X), int32(turn.Move.Point.Y), int(turn.Move.Disc))
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}

func (t *TurnMySQLRepositoryImpl) FindForGameIDAndTurnCount(ctx context.Context, db *sql.DB, gameID int, turnCount int) (Turn turn.Turn, err error) {
	defer xerrors.Wrap(&err, "FindForGameIDAndTurnCount")
	tgw := NewTurnGateway(othello.New(db))
	sgw := NewSquareGateway(othello.New(db))
	mgw := NewMoveGateway(othello.New(db))

	turnRecord, err := tgw.FindForGameIDAndTurnCount(ctx, int(gameID), turnCount)
	if err != nil {
		log.Println(err)
		return Turn, err
	}

	squaresRecord, err := sgw.FindForTurnID(ctx, int(turnRecord.ID))
	if err != nil {
		log.Println(err)
		return Turn, err
	}

	board := make([][]turn.Disc, 8)
	for i := range board {
		board[i] = make([]turn.Disc, 8)
	}

	for _, square := range squaresRecord {
		disc, err := turn.ToDisc((square.GetDisc()))
		if err != nil {
			return Turn, err
		}

		board[square.GetY()][square.GetX()] = disc
	}

	moveRecord, err := mgw.FindForTurnID(ctx, int(turnRecord.ID))
	if err != nil {
		return Turn, err
	}

	point, err := turn.NewPoint(int(moveRecord.GetX()), int(moveRecord.GetY()))
	if err != nil {
		return Turn, err
	}

	move := turn.NewMove(turn.Disc(moveRecord.GetDisc()), point)

	return turn.NewTurn(
		int(turnRecord.GetGameID()),
		int(turnRecord.GetTurnCount()),
		turn.Disc(turnRecord.GetNextDisc()),
		move,
		turn.NewBoard(board),
		turnRecord.GetEndAt(),
	), nil
}
