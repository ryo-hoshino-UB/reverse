package application

import (
	"api/dataaccess"
	othello "api/generated"
	"api/lib"
	"context"
	"log"
)

type TurnService struct{}

func NewTurnService() *TurnService {
	return &TurnService{}
}

func (t *TurnService) RegisterTurn(ctx context.Context, turnCount int, disc int, x int32, y int32) error {
	db := dataaccess.ConnectDB()

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

	board := make([][]int, 8)
	for i := range board {
		board[i] = make([]int, 8)
	}
	for _, square := range squaresRecord {
		board[square.GetY()][square.GetX()] = int(square.GetDisc())
	}

	// 石を置く
	board[y][x] = disc

	// ターンを保存する
	var nextDisc int
	if disc == lib.BLACK {
		nextDisc = lib.WHITE
	} else {
		nextDisc = lib.BLACK
	}

	turnRecord, err := tgw.Insert(ctx, gameRecord.GetID(), turnCount, nextDisc)
	if err != nil {
		log.Fatal(err)
		tx.Rollback()
		return err
	}

	err = sgw.InsertAll(ctx, int(turnRecord.GetID()), board)
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
	TurnCount  int       `json:"turnCount"`
	Board      [8][8]int `json:"board"`
	NextDisc   int       `json:"nextDisc"`
	WinnerDisc int       `json:"winnerDisc"`
}

func (o FindLatestGameTurnByTurnCountOutput) GetTurnCount() int {
	return o.TurnCount
}

func (o FindLatestGameTurnByTurnCountOutput) GetBoard() [8][8]int {
	return o.Board
}

func (o FindLatestGameTurnByTurnCountOutput) GetNextDisc() int {
	return o.NextDisc
}

func (o FindLatestGameTurnByTurnCountOutput) GetWinnerDisc() int {
	return o.WinnerDisc
}

func NewFindLatestGameTurnByTurnCountOutput(turnCount int, board [8][8]int, nextDisc int, winnerDisc int) FindLatestGameTurnByTurnCountOutput {
	return FindLatestGameTurnByTurnCountOutput{
		TurnCount:  turnCount,
		Board:      board,
		NextDisc:   nextDisc,
		WinnerDisc: winnerDisc,
	}
}

func (t *TurnService) FindLatestGameTurnByTurnCount(ctx context.Context, turnCount int) (res FindLatestGameTurnByTurnCountOutput, found bool) {
	db := dataaccess.ConnectDB()

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

	board := [8][8]int{}
	for _, square := range squaresRecord {
		board[square.GetY()][square.GetX()] = int(square.GetDisc())
	}

	res = NewFindLatestGameTurnByTurnCountOutput(turnRecord.GetTurnCount(), board, turnRecord.GetNextDisc(), 0)
	
	return res, true
}
