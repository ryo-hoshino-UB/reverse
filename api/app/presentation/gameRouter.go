package presentation

import (
	"api/application"
	"context"
	"database/sql"
	"net/http"

	"github.com/labstack/echo/v4"
)

func GameRouter(ctx context.Context, db *sql.DB) func(e *echo.Echo) {
	return func(e *echo.Echo) {
		// ゲームに関連するAPI群
		games := e.Group("/api/games")
		gs := application.NewGameService()

		// ゲーム作成API
		games.POST("", func(c echo.Context) error {
			if err := gs.StartNewGame(ctx, db); err != nil {
				return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to create game"})
			}
			return c.JSON(http.StatusCreated, map[string]string{"message": "game created"})
		})
	}
}
