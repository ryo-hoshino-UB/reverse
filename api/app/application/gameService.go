package application

import (
	"api/domain/game"
	"api/domain/turn"
	"context"
	"database/sql"
	"log"
	"time"
)

type GameService struct{}

func NewGameService() *GameService {
	return &GameService{}
}

func (g *GameService) StartNewGame(ctx context.Context, db *sql.DB) error {
	gr := game.NewGameRepository()
	tr := turn.NewTurnRepository()

	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		log.Fatal(err)
		return err
	}

	game, err := gr.Save(ctx, db)
	if err != nil {
		log.Fatal(err)
		tx.Rollback()
		return err
	}

	firstTurn := turn.NewFirstTurn(game.GetID(), time.Now())

	err = tr.Save(ctx, db, firstTurn)
	if err != nil {
		log.Fatal(err)
		tx.Rollback()
		return err
	}

	err = tx.Commit()
	if err != nil {
		log.Fatal(err)
		tx.Rollback()
		return err
	}

	return nil
}
