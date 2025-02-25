package dataaccess

import (
	othello "api/generated"
	"context"
	"log"
)

type GameGateway interface {
	FindLatest()
	Insert()
}

type GameGatewayImpl struct {
	Context context.Context
	Queries *othello.Queries
}

func (g GameGatewayImpl) FindLatest() (GameRecord, error) {
	latestGame, err := g.Queries.GetLatestGame(g.Context)
	if err != nil {
		log.Println(err)
		return GameRecord{}, err
	}

	return GameRecord{
		Game: latestGame,
	}, nil
}

func (g GameGatewayImpl) Insert() (GameRecord, error) {
	insertRes, err := g.Queries.CreateGame(g.Context)
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
