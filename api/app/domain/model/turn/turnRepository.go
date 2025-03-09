package turn

import (
	"context"
	"database/sql"
)

type TurnRepository interface {
	Save(ctx context.Context, db *sql.DB, turn Turn) error
	FindForGameIDAndTurnCount(ctx context.Context, db *sql.DB, gameID int, turnCount int) (Turn, error)
}
