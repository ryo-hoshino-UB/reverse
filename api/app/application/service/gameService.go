package application

import (
	gameModel "api/domain/model/game"
	turnModel "api/domain/model/turn"
	"context"
	"database/sql"
	"log"
	"time"
)

type GameService struct{
	gameRepository gameModel.GameRepository
	turnRepository turnModel.TurnRepository
}

func NewGameService(gameRepository gameModel.GameRepository, turnRepository turnModel.TurnRepository) *GameService {
	return &GameService{
		gameRepository: gameRepository,
		turnRepository: turnRepository,
	}
}

func (g *GameService) StartNewGame(ctx context.Context, db *sql.DB) error {
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		log.Println(err)
		return err
	}

	game, err := g.gameRepository.Save(ctx, db)
	if err != nil {
		log.Println(err)
		tx.Rollback()
		return err
	}

	firstTurn := turnModel.NewFirstTurn(game.GetID(), time.Now())

	err = g.turnRepository.Save(ctx, db, firstTurn)
	if err != nil {
		log.Println(err)
		tx.Rollback()
		return err
	}

	err = tx.Commit()
	if err != nil {
		log.Println(err)
		tx.Rollback()
		return err
	}

	return nil
}
