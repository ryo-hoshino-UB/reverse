package game

import (
	"api/domain/model/game"
	othello "api/generated"
	"context"
	"log"
)

type GameGateway struct {
	Queries *othello.Queries
}

func NewGameGateway(q *othello.Queries) *GameGateway {
	return &GameGateway{
		Queries: q,
	}
}

func (g *GameGateway) FindLatest(ctx context.Context) (game.Game, error) {
	latestGame, err := g.Queries.GetLatestGame(ctx)
	if err != nil {
		log.Println(err)
		return game.Game{}, err
	}

	return game.Game{
		ID: int(latestGame.ID),
		StartedAt: latestGame.StartedAt,
	}, nil
}

func (g *GameGateway) Save(ctx context.Context) (game.Game, error) {
	insertRes, err := g.Queries.CreateGame(ctx)
	if err != nil {
		log.Println(err)
		return game.Game{}, err
	}

	gameID, err := insertRes.LastInsertId()
	if err != nil {
		log.Println(err)
		return game.Game{}, err
	}

	return game.Game{
		ID: int(gameID),
	}, nil
}
