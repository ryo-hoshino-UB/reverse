package game

import (
	othello "api/generated"
	infrastucture "api/infrastructure"
	"api/xerrors"
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
		if err == sql.ErrNoRows {
			return Game{}, xerrors.ErrNotFound
		}
		log.Println(err)
		return Game{}, err
	}

	return NewGame(gameRecord.GetID(), gameRecord.GetStartedAt()), nil
}

func (r GameRepository) Save(ctx context.Context, db *sql.DB) (Game, error) {
	ggw := infrastucture.NewGameGateway(othello.New(db))

	gameRecord, err := ggw.Insert(ctx)
	if err != nil {
		log.Println(err)
		return Game{}, err
	}

	return NewGame(gameRecord.GetID(), gameRecord.GetStartedAt()), nil
}
