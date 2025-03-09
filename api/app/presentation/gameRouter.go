package presentation

import (
	"api/application/usecase"
	"api/infrastructure/repository/game"
	"api/infrastructure/repository/turn"
	"context"
	"database/sql"
	"net/http"

	"github.com/labstack/echo/v4"
)

func GameRouter(ctx context.Context, db *sql.DB) func(e *echo.Echo) {
	gr := game.NewGameMySQLRepositoryImpl()
	tr := turn.NewTurnMySQLRepositoryImpl()
	startNewGame := usecase.NewStartNewGame(gr, tr)
	
	return func(e *echo.Echo) {
		// ゲームに関連するAPI群
		games := e.Group("/api/games")

		// ゲーム作成API
		games.POST("", func(c echo.Context) error {
			if err := startNewGame.Run(ctx, db); err != nil {
				return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to create game"})
			}
			return c.JSON(http.StatusCreated, map[string]string{"message": "game created"})
		})
	}
}
