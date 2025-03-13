package query

import (
	"context"
	"database/sql"
	"time"
)

type FindLastGamesQueryModel struct {
	GameID         int
	BlackMoveCount int
	WhiteMoveCount int
	WinnerDisc     int
	StartedAt      time.Time
	EndAt          time.Time
}

func NewFindLastGamesQueryModel(gameID int, blackMoveCount int, whiteMoveCount int, winnerDisc int, startedAt time.Time, endAt time.Time) FindLastGamesQueryModel {
	return FindLastGamesQueryModel{
		GameID:         gameID,
		BlackMoveCount: blackMoveCount,
		WhiteMoveCount: whiteMoveCount,
		WinnerDisc:     winnerDisc,
		StartedAt:      startedAt,
		EndAt:          endAt,
	}
}

func (m *FindLastGamesQueryModel) GetGameID() int {
	return m.GameID
}

func (m *FindLastGamesQueryModel) GetBlackMoveCount() int {
	return m.BlackMoveCount
}

func (m *FindLastGamesQueryModel) GetWhiteMoveCount() int {
	return m.WhiteMoveCount
}

func (m *FindLastGamesQueryModel) GetWinnerDisc() int {
	return m.WinnerDisc
}

func (m *FindLastGamesQueryModel) GetStartedAt() time.Time {
	return m.StartedAt
}

func (m *FindLastGamesQueryModel) GetEndAt() time.Time {
	return m.EndAt
}

type FindLastGamesQueryService interface {
	Query(ctx context.Context, db *sql.DB, limit int) ([]FindLastGamesQueryModel, error)
}
