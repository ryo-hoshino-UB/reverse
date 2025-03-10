package usecase

import (
	"api/domain/model/game"
	"api/domain/model/gameresult"
	"api/domain/model/turn"
	"context"
	"database/sql"
	"log"
)

type RegisterTurn struct {
	turnPort    turn.TurnRepository
	gamePort    game.GameRepository
	gameResPort gameresult.GameResultRepository
}

func NewRegisterTurn(
	turnPort turn.TurnRepository,
	gamePort game.GameRepository,
	gameResPort gameresult.GameResultRepository,
) *RegisterTurn {
	return &RegisterTurn{
		turnPort:    turnPort,
		gamePort:    gamePort,
		gameResPort: gameResPort,
	}
}

func (t *RegisterTurn) Run(ctx context.Context, db *sql.DB, turnCount int, disc turn.Disc, x int32, y int32) error {
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
