package gameresult

import (
	"context"
	"database/sql"
)

type GameResultRepository interface {
	FindForGameID(ctx context.Context, db *sql.DB, gameID int) (GameResult, error)
	Save(ctx context.Context, db *sql.DB, gr GameResult) error
}
