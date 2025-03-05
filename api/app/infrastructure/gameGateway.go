package infrastucture

import (
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

func (g *GameGateway) FindLatest(ctx context.Context) (GameRecord, error) {
	latestGame, err := g.Queries.GetLatestGame(ctx)
	if err != nil {
		log.Println(err)
		return GameRecord{}, err
	}

	return GameRecord{
		Game: latestGame,
	}, nil
}

func (g *GameGateway) Insert(ctx context.Context) (GameRecord, error) {
	insertRes, err := g.Queries.CreateGame(ctx)
	if err != nil {
		log.Println(err)
		return GameRecord{}, err
	}

	gameID, err := insertRes.LastInsertId()
	if err != nil {
		log.Println(err)
		return GameRecord{}, err
	}

	return GameRecord{
		Game: othello.Game{
			ID: int32(gameID),
		},
	}, nil
}
