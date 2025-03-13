package presentation

import (
	"api/application/usecase"
	"api/infrastructure/query"
	"api/infrastructure/repository/game"
	"api/infrastructure/repository/turn"
	"context"
	"database/sql"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
)

type GameHistoriesResponse struct {
	GameID         int       `json:"id"`
	BlackMoveCount int       `json:"blackMoveCount"`
	WhiteMoveCount int       `json:"whiteMoveCount"`
	StartedAt      time.Time `json:"startedAt"`
	WinnerDisc     int       `json:"winnerDisc"`
	EndAt          time.Time `json:"endAt"`
}

func GameRouter(ctx context.Context, db *sql.DB) func(e *echo.Echo) {
	gr := game.NewGameMySQLRepositoryImpl()
	tr := turn.NewTurnMySQLRepositoryImpl()
	gqs := query.NewFindLastGamesMySQLQueryService()

	startNewGame := usecase.NewStartNewGame(gr, tr)
	findLastGames := usecase.NewFindLastGames(gqs)

	return func(e *echo.Echo) {
		// ゲームに関連するAPI群
		games := e.Group("/api/games")

		// ゲーム一覧取得API
		games.GET("", func(c echo.Context) error {
			limit := 10
			gameHistories, err := findLastGames.Run(ctx, db, limit)
			if err != nil {
				return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to get game histories"})
			}

			gameHistoriesResponse := make([]GameHistoriesResponse, len(gameHistories))
			for i, history := range gameHistories {
				gameHistoriesResponse[i] = GameHistoriesResponse{
					GameID:         history.GetGameID(),
					BlackMoveCount: history.GetBlackMoveCount(),
					WhiteMoveCount: history.GetWhiteMoveCount(),
					StartedAt:      history.GetStartedAt(),
					EndAt:          history.GetEndAt(),
				}
			}

			return c.JSON(http.StatusOK, gameHistoriesResponse)
		})

		// ゲーム作成API
		games.POST("", func(c echo.Context) error {
			if err := startNewGame.Run(ctx, db); err != nil {
				return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to create game"})
			}
			return c.JSON(http.StatusCreated, map[string]string{"message": "game created"})
		})
	}
}
