package presentation

import (
	"api/application/usecase"
	"api/domain/model/turn"
	othello "api/generated"
	"api/infrastructure/repository/game"
	"api/infrastructure/repository/gameresult"
	turnRepo "api/infrastructure/repository/turn"
	"api/xerrors"
	"context"
	"database/sql"
	"log"
	"net/http"
	"strconv"

	"github.com/go-errors/errors"
	"github.com/labstack/echo/v4"
)

type TurnPostRequest struct {
	TurnCount int                      `json:"turnCount"`
	Move      othello.CreateMoveParams `json:"move"`
}

type TurnGetResponse struct {
	TurnCount  int           `json:"turnCount"`
	Board      [][]turn.Disc `json:"board"`
	NextDisc   int           `json:"nextDisc"`
	WinnerDisc int           `json:"winnerDisc"`
}

func TurnRouter(ctx context.Context, db *sql.DB) func(e *echo.Echo) {
	gr := game.NewGameMySQLRepositoryImpl()
	tr := turnRepo.NewTurnMySQLRepositoryImpl()
	grr := gameresult.NewGameResultMySQLRepositoryImpl()
	registerTurn := usecase.NewRegisterTurn(tr, gr, grr)
	findLatestGameTurnByTurnCount := usecase.NewFindLatestGameTurnByTurnCount(tr, gr, grr)

	return func(e *echo.Echo) {
		turns := e.Group("/api/games/latest/turns")

		turns.POST("", func(c echo.Context) error {
			var turnReq TurnPostRequest
			if err := c.Bind(&turnReq); err != nil {
				return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request"})
			}

			disc, err := turn.ToDisc(int(turnReq.Move.Disc))
			if err != nil {
				return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid disc"})
			}

			err = registerTurn.Run(ctx, db, turnReq.TurnCount, disc, turnReq.Move.X, turnReq.Move.Y)
			if err != nil {
				log.Println(err)
				if errors.Is(err, xerrors.ErrBadRequest) {
					return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid turn"})
				}
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

			turnOutput, err := findLatestGameTurnByTurnCount.Run(ctx, db, turnCount)
			if errors.Is(err, xerrors.ErrNotFound) {
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
