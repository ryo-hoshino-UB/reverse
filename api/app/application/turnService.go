package application

import (
	"api/domain"
	"api/domain/game"
	"api/domain/turn"
	"context"
	"database/sql"
	"log"
)

type TurnService struct{}

func NewTurnService() *TurnService {
	return &TurnService{}
}

func (t *TurnService) RegisterTurn(ctx context.Context, db *sql.DB, turnCount int, disc int, x int32, y int32) error {
	tr := turn.NewTurnRepository()
	gr := game.NewGameRepository()

	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		log.Fatal(err)
		tx.Rollback()
		return err
	}

	// 1つ前のターンを取得
	game, err := gr.FindLatest(ctx, db)
	if err != nil {
		log.Fatal(err)
		tx.Rollback()
		return err
	}

	prevTurnCount := turnCount - 1
	prevTurn, err := tr.FindForGameIDAndTurnCount(ctx, db, int(game.GetID()), prevTurnCount)
	if err != nil {
		log.Fatal(err)
		tx.Rollback()
		return err
	}

	// 石を置く
	newTurn := prevTurn.PlaceNext(domain.ToDisc(disc), domain.NewPoint(int(x), int(y)))

	// ターンを保存する
	err = tr.Save(ctx, db, newTurn)
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
	tr := turn.NewTurnRepository()
	gr := game.NewGameRepository()

	gameRecord, err := gr.FindLatest(ctx, db)

	turnRecord, err := tr.FindForGameIDAndTurnCount(ctx, db, int(gameRecord.GetID()), turnCount)
	if err != nil {
		log.Fatal(err)
		return FindLatestGameTurnByTurnCountOutput{}, false
	}

	res = NewFindLatestGameTurnByTurnCountOutput(turnRecord.GetTurnCount(), turnRecord.Board.Discs, int(turnRecord.GetNextDisc()), 0)

	return res, true
}
