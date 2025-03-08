package gameresult

import (
	othello "api/generated"
	"api/infrastructure"
	"context"
	"database/sql"
)

type GameResultRepository struct {
}

func NewGameResultRepository() *GameResultRepository {
	return &GameResultRepository{}
}

func (r *GameResultRepository) FindForGameID(ctx context.Context, db *sql.DB, gameID int) (GameResult, error) {
	grgw := infrastructure.NewGameResultGateway(othello.New(db))

	gr, err := grgw.SelectGameResult(ctx, gameID)
	if err != nil {
		return GameResult{}, err
	}

	return *NewGameResult(
		gr.GetGameID(),
		ToWinnerDisc(gr.GetWinnerDisc()),
	), nil
}

func (r *GameResultRepository) Save(ctx context.Context, db *sql.DB, gr GameResult) error {
	grgw := infrastructure.NewGameResultGateway(othello.New(db))

	_, err := grgw.Insert(ctx, gr.GetGameID(), int(gr.GetWinnerDisc()), gr.GetEndAt())
	if err != nil {
		return err
	}

	return nil
}
