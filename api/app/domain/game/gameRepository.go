package game

import (
	othello "api/generated"
	infrastucture "api/infrastructure"
	"context"
	"database/sql"
	"log"
)

type GameRepository struct{}

func NewGameRepository() *GameRepository {
	return &GameRepository{}
}

func (r GameRepository) FindLatest(ctx context.Context, db *sql.DB) (Game, error) {
	ggw := infrastucture.NewGameGateway(othello.New(db))

	gameRecord, err := ggw.FindLatest(ctx)
	if err != nil {
		log.Fatal(err)
		return Game{}, err
	}

	return NewGame(gameRecord.GetID(), gameRecord.GetStartedAt()), nil
}

func (r GameRepository) Save(ctx context.Context, db *sql.DB) (Game, error) {
	ggw := infrastucture.NewGameGateway(othello.New(db))

	gameRecord, err := ggw.Insert(ctx)
	if err != nil {
		log.Fatal(err)
		return Game{}, err
	}

	return NewGame(gameRecord.GetID(), gameRecord.GetStartedAt()), nil
}
