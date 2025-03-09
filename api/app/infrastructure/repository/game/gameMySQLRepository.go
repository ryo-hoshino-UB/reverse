package game

import (
	"api/domain/model/game"
	othello "api/generated"
	"api/xerrors"
	"context"
	"database/sql"
	"log"
)

type GameMySQLRepositoryImpl struct{}

func NewGameMySQLRepositoryImpl() *GameMySQLRepositoryImpl {
	return &GameMySQLRepositoryImpl{}
}

func (g GameMySQLRepositoryImpl) FindLatest(ctx context.Context, db *sql.DB) (game.Game, error) {
	ggw := NewGameGateway(othello.New(db))
	gameRecord, err := ggw.FindLatest(ctx)
	if err != nil {
		if err == sql.ErrNoRows {
			return game.Game{}, xerrors.ErrNotFound
		}
		log.Println(err)
		return game.Game{}, err
	}

	return game.NewGame(gameRecord.GetID(), gameRecord.GetStartedAt()), nil
}

func (g GameMySQLRepositoryImpl) Save(ctx context.Context, db *sql.DB) (game.Game, error) {
	ggw := NewGameGateway(othello.New(db))
	gameRecord, err := ggw.Save(ctx)
	if err != nil {
		log.Println(err)
		return game.Game{}, err
	}

	return game.NewGame(gameRecord.GetID(), gameRecord.GetStartedAt()), nil
}
