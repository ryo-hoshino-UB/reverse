package query

import (
	queryModel "api/application/query"
	othello "api/generated"
	"context"
	"database/sql"
)

type FindLastGamesMySQLQueryService struct {
}

func NewFindLastGamesMySQLQueryService() *FindLastGamesMySQLQueryService {
	return &FindLastGamesMySQLQueryService{}
}

func (s *FindLastGamesMySQLQueryService) Query(ctx context.Context, db *sql.DB, limit int) ([]queryModel.FindLastGamesQueryModel, error) {
	queries := othello.New(db)
	gameHistoryRecords, err := queries.GetGameHistories(ctx, int32(limit))
	if err != nil {
		return nil, err
	}

	gameHistories := make([]queryModel.FindLastGamesQueryModel, len(gameHistoryRecords))
	for i, record := range gameHistoryRecords {
		gameHistories[i] = queryModel.NewFindLastGamesQueryModel(
			int(record.GameID),
			int(record.BlackMoveCount),
			int(record.WhiteMoveCount),
			int(record.WinnerDisc.Int32),
			record.StartedAt,
			record.EndAt.Time,
		)
	}

	return gameHistories, nil
}
