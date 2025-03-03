package presentation

import (
	"api/application"
	othello "api/generated"
	"context"
	"database/sql"
	"log"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

type RegisterTurnRequest struct {
	TurnCount int                      `json:"turnCount"`
	Move      othello.CreateMoveParams `json:"move"`
}

func TurnRouter(ctx context.Context, db *sql.DB) func(e *echo.Echo) {
	return func(e *echo.Echo) {
		turns := e.Group("/api/games/latest/turns")
		ts := application.NewTurnService()

		turns.POST("", func(c echo.Context) error {
			var turnReq RegisterTurnRequest
			if err := c.Bind(&turnReq); err != nil {
				return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request"})
			}

			err := ts.RegisterTurn(ctx, turnReq.TurnCount, int(turnReq.Move.Disc), turnReq.Move.X, turnReq.Move.Y)
			if err != nil {
				log.Fatal(err)
				return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to register turn"})
			}

			return c.JSON(http.StatusCreated, nil)
		})

		turns.GET("/:turnCount", func(c echo.Context) error {
			turnCountStr := c.Param("turnCount")
			turnCount, err := strconv.Atoi(turnCountStr)
			if err != nil {
				return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid turn count"})
			}

			res, found := ts.FindLatestGameTurnByTurnCount(ctx, turnCount)
			if !found {
				return c.JSON(http.StatusNotFound, map[string]string{"error": "turn not found"})
			}

			return c.JSON(http.StatusOK, res)
		})
	}
}
