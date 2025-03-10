package game

import (
	"context"
	"database/sql"
)

type GameRepository interface {
	FindLatest(ctx context.Context, db *sql.DB) (Game, error)
	Save(ctx context.Context, db *sql.DB) (Game, error)
}
