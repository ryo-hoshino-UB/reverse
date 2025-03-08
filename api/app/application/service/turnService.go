package application

import (
	"api/domain"
	"api/domain/model/game"
	gameresult "api/domain/model/gameResult"
	"api/domain/model/turn"
	"context"
	"database/sql"
	"log"
)

type TurnService struct{}

func NewTurnService() *TurnService {
	return &TurnService{}
}

func (t *TurnService) RegisterTurn(ctx context.Context, db *sql.DB, turnCount int, disc domain.Disc, x int32, y int32) error {
	tr := turn.NewTurnRepository()
	gr := game.NewGameRepository()
	grr := gameresult.NewGameResultRepository()

	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		log.Println(err)
		tx.Rollback()
		return err
	}

	// 1つ前のターンを取得
	game, err := gr.FindLatest(ctx, db)
	if err != nil {
		log.Println(err)
		tx.Rollback()
		return err
	}

	prevTurnCount := turnCount - 1
	prevTurn, err := tr.FindForGameIDAndTurnCount(ctx, db, int(game.GetID()), prevTurnCount)
	// gameが取得できているのにturnが取得できない場合はアプリケーションのバグなので500を返す
	if err != nil {
		log.Println(err)
		tx.Rollback()
		return err
	}

	point, err := domain.NewPoint(int(x), int(y))
	if err != nil {
		log.Println(err)
		tx.Rollback()
		return err
	}

	// 石を置く
	newTurn, err := prevTurn.PlaceNext(disc, point)
	if err != nil {
		log.Println(err)
		tx.Rollback()
		return err
	}

	// ターンを保存する
	err = tr.Save(ctx, db, newTurn)
	if err != nil {
		log.Println(err)
		tx.Rollback()
		return err
	}

	if newTurn.GameEnded() {
		winnerDisc := newTurn.WinnerDisc()
		gameResult := gameresult.NewGameResult(int(game.GetID()), winnerDisc)
		err = grr.Save(ctx, db, *gameResult)
		if err != nil {
			log.Println(err)
			tx.Rollback()
			return err
		}
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

func (t *TurnService) FindLatestGameTurnByTurnCount(ctx context.Context, db *sql.DB, turnCount int) (res FindLatestGameTurnByTurnCountOutput, err error) {
	tr := turn.NewTurnRepository()
	gr := game.NewGameRepository()
	grr := gameresult.NewGameResultRepository()

	game, err := gr.FindLatest(ctx, db)

	turn, err := tr.FindForGameIDAndTurnCount(ctx, db, int(game.GetID()), turnCount)
	if err != nil {
		log.Println(err)
		return FindLatestGameTurnByTurnCountOutput{}, err
	}

	var gameResult gameresult.GameResult
	if turn.GameEnded() {
		gameResult, err = grr.FindForGameID(ctx, db, int(game.GetID()))
		if err != nil {
			log.Println(err)
			return FindLatestGameTurnByTurnCountOutput{}, err
		}
	}

	res = NewFindLatestGameTurnByTurnCountOutput(turn.GetTurnCount(), turn.Board.Discs, int(turn.GetNextDisc()), int(gameResult.GetWinnerDisc()))

	return res, nil
}
