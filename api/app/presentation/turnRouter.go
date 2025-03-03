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

type TurnPostRequest struct {
	TurnCount int                      `json:"turnCount"`
	Move      othello.CreateMoveParams `json:"move"`
}

type TurnGetResponse struct {
	TurnCount  int       `json:"turnCount"`
	Board      [8][8]int `json:"board"`
	NextDisc   int       `json:"nextDisc"`
	WinnerDisc int       `json:"winnerDisc"`
}

func TurnRouter(ctx context.Context, db *sql.DB) func(e *echo.Echo) {
	return func(e *echo.Echo) {
		turns := e.Group("/api/games/latest/turns")
		ts := application.NewTurnService()

		turns.POST("", func(c echo.Context) error {
			var turnReq TurnPostRequest
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

			turnOutput, found := ts.FindLatestGameTurnByTurnCount(ctx, turnCount)
			if !found {
				return c.JSON(http.StatusNotFound, map[string]string{"error": "turn not found"})
			}
			
			res := TurnGetResponse{
				TurnCount:  turnOutput.TurnCount,
				Board:      turnOutput.Board,
				NextDisc:   turnOutput.NextDisc,
				WinnerDisc: turnOutput.WinnerDisc,
			}

			return c.JSON(http.StatusOK, res)
		})
	}
}
