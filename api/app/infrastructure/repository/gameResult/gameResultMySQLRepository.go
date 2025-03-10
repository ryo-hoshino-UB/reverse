package gameresult

import (
	"api/domain/model/gameresult"
	othello "api/generated"
	"context"
	"database/sql"
)

type GameResultMySQLRepositoryImpl struct{}

func NewGameResultMySQLRepositoryImpl() *GameResultMySQLRepositoryImpl {
	return &GameResultMySQLRepositoryImpl{}
}

func (r *GameResultMySQLRepositoryImpl) FindForGameID(ctx context.Context, db *sql.DB, gameID int) (gameresult.GameResult, error) {
	grgw := NewGameResultGateway(othello.New(db))

	gr, err := grgw.SelectGameResult(ctx, gameID)
	if err != nil {
		return gameresult.GameResult{}, err
	}

	return *gameresult.NewGameResult(
		gr.GetGameID(),
		gameresult.ToWinnerDisc(gr.GetWinnerDisc()),
	), nil
}

func (r *GameResultMySQLRepositoryImpl) Save(ctx context.Context, db *sql.DB, gr gameresult.GameResult) error {
	grgw := NewGameResultGateway(othello.New(db))

	_, err := grgw.Insert(ctx, gr.GetGameID(), int(gr.GetWinnerDisc()), gr.GetEndAt())
	if err != nil {
		return err
	}

	return nil
}
