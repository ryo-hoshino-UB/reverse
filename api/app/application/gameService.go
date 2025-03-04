package application

import (
	"api/dataaccess"
	othello "api/generated"
	"api/lib"
	"context"
	"database/sql"
	"log"
)

type GameService struct{}

func NewGameService() *GameService {
	return &GameService{}
}

func (g *GameService) StartNewGame(ctx context.Context, db *sql.DB) error {
	ggw := dataaccess.NewGameGateway(othello.New(db))
	tgw := dataaccess.NewTurnGateway(othello.New(db))
	sgw := dataaccess.NewSquareGateway(othello.New(db))

	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		log.Fatal(err)
		return err
	}

	gameRecord, err := ggw.Insert(ctx)
	if err != nil {
		log.Fatal(err)
		tx.Rollback()
		return err
	}

	insertedGameID := gameRecord.GetID()

	turnRecord, err := tgw.Insert(ctx, insertedGameID, 0, lib.BLACK)
	if err != nil {
		log.Fatal(err)
		tx.Rollback()
		return err
	}

	err = sgw.InsertAll(ctx, int(turnRecord.GetID()), lib.INITIAL_BOARD)
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
