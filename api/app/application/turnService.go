package application

import (
	"api/dataaccess"
	"api/domain"
	othello "api/generated"
	"context"
	"database/sql"
	"log"
)

type TurnService struct{}

func NewTurnService() *TurnService {
	return &TurnService{}
}

func (t *TurnService) RegisterTurn(ctx context.Context, db *sql.DB, turnCount int, disc int, x int32, y int32) error {
	ggw := dataaccess.NewGameGateway(othello.New(db))
	tgw := dataaccess.NewTurnGateway(othello.New(db))
	sgw := dataaccess.NewSquareGateway(othello.New(db))
	mgw := dataaccess.NewMoveGateway(othello.New(db))

	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		log.Fatal(err)
		tx.Rollback()
		return err
	}

	// 1つ前のターンを取得
	gameRecord, err := ggw.FindLatest(ctx)
	if err != nil {
		log.Fatal(err)
		tx.Rollback()
		return err
	}

	prevTurnCount := turnCount - 1
	prevTurnRecord, err := tgw.FindForGameIDAndTurnCount(ctx, int(gameRecord.GetID()), prevTurnCount)
	if err != nil {
		log.Fatal(err)
		tx.Rollback()
		return err
	}

	squaresRecord, err := sgw.FindForTurnID(ctx, int(prevTurnRecord.GetID()))
	if err != nil {
		log.Fatal(err)
		tx.Rollback()
		return err
	}

	board := make([][]domain.Disc, 8)
	for i := range board {
		board[i] = make([]domain.Disc, 8)
	}
	for _, square := range squaresRecord {
		board[square.GetY()][square.GetX()] = domain.ToDisc(square.GetDisc())
	}

	// 石を置く
	board[y][x] = domain.ToDisc(disc)

	prevTurn := domain.NewTurn(prevTurnRecord.GetGameID(), prevTurnRecord.GetTurnCount(), domain.ToDisc(prevTurnRecord.GetNextDisc()), domain.Move{}, domain.NewBoard(board), prevTurnRecord.GetEndAt())

	newTurn := prevTurn.PlaceNext(domain.ToDisc(disc), domain.NewPoint(int(x), int(y)))

	// ターンを保存する
	turnRecord, err := tgw.Insert(ctx, newTurn.GetGameID(), newTurn.GetTurnCount(), newTurn.GetNextDisc())
	if err != nil {
		log.Fatal(err)
		tx.Rollback()
		return err
	}

	err = sgw.InsertAll(ctx, int(turnRecord.GetID()), newTurn.Board.GetDiscs())
	if err != nil {
		log.Fatal(err)
		tx.Rollback()
		return err
	}

	_, err = mgw.Insert(ctx, int(turnRecord.GetID()), x, y, disc)
	if err != nil {
		log.Fatal(err)
		tx.Rollback()
		return err
	}

	tx.Commit()
	return nil
}

type FindLatestGameTurnByTurnCountOutput struct {
	TurnCount  int             `json:"turnCount"`
	Board      [][]domain.Disc `json:"board"`
	NextDisc   int             `json:"nextDisc"`
	WinnerDisc int             `json:"winnerDisc"`
}

func (o FindLatestGameTurnByTurnCountOutput) GetTurnCount() int {
	return o.TurnCount
}

func (o FindLatestGameTurnByTurnCountOutput) GetBoard() [][]domain.Disc {
	return o.Board
}

func (o FindLatestGameTurnByTurnCountOutput) GetNextDisc() int {
	return o.NextDisc
}

func (o FindLatestGameTurnByTurnCountOutput) GetWinnerDisc() int {
	return o.WinnerDisc
}

func NewFindLatestGameTurnByTurnCountOutput(turnCount int, board [][]domain.Disc, nextDisc int, winnerDisc int) FindLatestGameTurnByTurnCountOutput {
	return FindLatestGameTurnByTurnCountOutput{
		TurnCount:  turnCount,
		Board:      board,
		NextDisc:   nextDisc,
		WinnerDisc: winnerDisc,
	}
}

func (t *TurnService) FindLatestGameTurnByTurnCount(ctx context.Context, db *sql.DB, turnCount int) (res FindLatestGameTurnByTurnCountOutput, found bool) {
	ggw := dataaccess.NewGameGateway(othello.New(db))
	tgw := dataaccess.NewTurnGateway(othello.New(db))
	sgw := dataaccess.NewSquareGateway(othello.New(db))

	gameRecord, err := ggw.FindLatest(ctx)
	if err != nil {
		log.Fatal(err)
		return FindLatestGameTurnByTurnCountOutput{}, false
	}

	turnRecord, err := tgw.FindForGameIDAndTurnCount(ctx, int(gameRecord.GetID()), turnCount)
	if err != nil {
		log.Fatal(err)
		return FindLatestGameTurnByTurnCountOutput{}, false
	}

	squaresRecord, err := sgw.FindForTurnID(ctx, int(turnRecord.GetID()))
	if err != nil {
		log.Fatal(err)
		return FindLatestGameTurnByTurnCountOutput{}, false
	}

	board := make([][]domain.Disc, 8)
	for i := range board {
		board[i] = make([]domain.Disc, 8)
	}

	for _, square := range squaresRecord {
		board[square.GetY()][square.GetX()] = domain.ToDisc(square.GetDisc())
	}

	res = NewFindLatestGameTurnByTurnCountOutput(turnRecord.GetTurnCount(), board, turnRecord.GetNextDisc(), 0)

	return res, true
}
