package application

import (
	"api/domain/model/game"
	"api/domain/model/gameresult"
	"api/domain/model/turn"
	"context"
	"database/sql"
	"log"
)

type TurnService struct{
	turnPort turn.TurnRepository
	gamePort game.GameRepository
	gameResPort gameresult.GameResultRepository
}

func NewTurnService(
	turnPort turn.TurnRepository,
	gamePort game.GameRepository,
	gameResPort gameresult.GameResultRepository,
) *TurnService {
	return &TurnService{
		turnPort: turnPort,
		gamePort: gamePort,
		gameResPort: gameResPort,
	}
}

func (t *TurnService) RegisterTurn(ctx context.Context, db *sql.DB, turnCount int, disc turn.Disc, x int32, y int32) error {
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		log.Println(err)
		tx.Rollback()
		return err
	}

	// 1つ前のターンを取得
	game, err := t.gamePort.FindLatest(ctx, db)
	if err != nil {
		log.Println(err)
		tx.Rollback()
		return err
	}

	prevTurnCount := turnCount - 1
	prevTurn, err := t.turnPort.FindForGameIDAndTurnCount(ctx, db, int(game.GetID()), prevTurnCount)
	// gameが取得できているのにturnが取得できない場合はアプリケーションのバグなので500を返す
	if err != nil {
		log.Println(err)
		tx.Rollback()
		return err
	}

	point, err := turn.NewPoint(int(x), int(y))
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
	err = t.turnPort.Save(ctx, db, newTurn)
	if err != nil {
		log.Println(err)
		tx.Rollback()
		return err
	}

	if newTurn.GameEnded() {
		winnerDisc := newTurn.WinnerDisc()
		gameResult := gameresult.NewGameResult(int(game.GetID()), winnerDisc)
		err = t.gameResPort.Save(ctx, db, *gameResult)
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
	Board      [][]turn.Disc `json:"board"`
	NextDisc   int             `json:"nextDisc"`
	WinnerDisc int             `json:"winnerDisc"`
}

func (o FindLatestGameTurnByTurnCountOutput) GetTurnCount() int {
	return o.TurnCount
}

func (o FindLatestGameTurnByTurnCountOutput) GetBoard() [][]turn.Disc {
	return o.Board
}

func (o FindLatestGameTurnByTurnCountOutput) GetNextDisc() int {
	return o.NextDisc
}

func (o FindLatestGameTurnByTurnCountOutput) GetWinnerDisc() int {
	return o.WinnerDisc
}

func NewFindLatestGameTurnByTurnCountOutput(turnCount int, board [][]turn.Disc, nextDisc int, winnerDisc int) FindLatestGameTurnByTurnCountOutput {
	return FindLatestGameTurnByTurnCountOutput{
		TurnCount:  turnCount,
		Board:      board,
		NextDisc:   nextDisc,
		WinnerDisc: winnerDisc,
	}
}

func (t *TurnService) FindLatestGameTurnByTurnCount(ctx context.Context, db *sql.DB, turnCount int) (res FindLatestGameTurnByTurnCountOutput, err error) {
	game, err := t.gamePort.FindLatest(ctx, db)

	turn, err := t.turnPort.FindForGameIDAndTurnCount(ctx, db, int(game.GetID()), turnCount)
	if err != nil {
		log.Println(err)
		return FindLatestGameTurnByTurnCountOutput{}, err
	}

	var gameResult gameresult.GameResult
	if turn.GameEnded() {
		gameResult, err = t.gameResPort.FindForGameID(ctx, db, int(game.GetID()))
		if err != nil {
			log.Println(err)
			return FindLatestGameTurnByTurnCountOutput{}, err
		}
	}

	res = NewFindLatestGameTurnByTurnCountOutput(turn.GetTurnCount(), turn.Board.Discs, int(turn.GetNextDisc()), int(gameResult.GetWinnerDisc()))

	return res, nil
}
