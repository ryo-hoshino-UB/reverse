package usecase

import (
	q "api/application/query"
	"context"
	"database/sql"
)

type FindLastGames struct {
	queryService q.FindLastGamesQueryService
}

func NewFindLastGames(queryService q.FindLastGamesQueryService) *FindLastGames {
	return &FindLastGames{queryService: queryService}
}

func (f *FindLastGames) Run(ctx context.Context, db *sql.DB, limit int) ([]q.FindLastGamesQueryModel, error) {
	gameHistories, err := f.queryService.Query(ctx, db, limit)
	if err != nil {
		return nil, err
	}

	return gameHistories, nil
}
