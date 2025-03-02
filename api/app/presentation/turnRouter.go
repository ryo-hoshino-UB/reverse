package presentation

import (
	"api/dataaccess"
	othello "api/generated"
	"api/lib"
	"context"
	"database/sql"
	"log"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

func TurnRouter(ctx context.Context, db *sql.DB) func(e *echo.Echo) {
	return func(e *echo.Echo) {
		queries := othello.New(db)
		gameGatewayImpl := dataaccess.GameGatewayImpl{
			Context: ctx,
			Queries: queries,
		}
		turnGatewayImpl := dataaccess.TurnGateway{
			Context: ctx,
			Queries: queries,
		}
		moveGatewayImpl := dataaccess.MoveGatewayImpl{
			Context: ctx,
			Queries: queries,
		}
		squareGatewayImpl := dataaccess.SquareGatewayImpl{
			Context: ctx,
			Queries: queries,
		}

		turns := e.Group("/api/games/latest/turns")

		turns.POST("", func(c echo.Context) error {
			var turnReq struct {
				TurnCount int                      `json:"turnCount"`
				Move      othello.CreateMoveParams `json:"move"`
			}

			if err := c.Bind(&turnReq); err != nil {
				return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request"})
			}

			turnCount := turnReq.TurnCount
			disc := int(turnReq.Move.Disc)
			x := turnReq.Move.X
			y := turnReq.Move.Y

			tx, err := db.BeginTx(ctx, nil)
			if err != nil {
				log.Fatal(err)
				tx.Rollback()
				return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to start transaction"})
			}

			// 1つ前のターンを取得
			gameRecord, err := gameGatewayImpl.FindLatest()
			if err != nil {
				log.Fatal(err)
				tx.Rollback()
				return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to get latest game"})
			}

			prevTurnCount := turnCount - 1
			prevTurnRecord, err := turnGatewayImpl.FindForGameIDAndTurnCount(int(gameRecord.GetID()), prevTurnCount)
			if err != nil {
				log.Fatal(err)
				tx.Rollback()
				return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to get turn"})
			}

			squaresRecord, err := squareGatewayImpl.FindForTurnID(int(prevTurnRecord.GetID()))
			if err != nil {
				log.Fatal(err)
				tx.Rollback()
				return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to get squares"})
			}

			board := make([][]int, 8)
			for i := range board {
				board[i] = make([]int, 8)
			}
			for _, square := range squaresRecord {
				board[square.GetY()][square.GetX()] = int(square.GetDisc())
			}

			// 石を置く
			board[y][x] = disc

			// ターンを保存する
			var nextDisc int
			if disc == lib.BLACK {
				nextDisc = lib.WHITE
			} else {
				nextDisc = lib.BLACK
			}

			turnRecord, err := turnGatewayImpl.Insert(gameRecord.GetID(), turnCount, nextDisc)
			if err != nil {
				log.Fatal(err)
				tx.Rollback()
				return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to insert turn"})
			}

			err = squareGatewayImpl.InsertAll(int(turnRecord.GetID()), board)
			if err != nil {
				log.Fatal(err)
				tx.Rollback()
				return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to insert squares"})
			}

			_, err = moveGatewayImpl.Insert(int(turnRecord.GetID()), x, y, disc)
			if err != nil {
				log.Fatal(err)
				tx.Rollback()
				return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to insert move"})
			}

			tx.Commit()
			return c.JSON(http.StatusOK, turnReq)
		})

		turns.GET("/:turnCount", func(c echo.Context) error {
			turnCountStr := c.Param("turnCount")
			turnCount, err := strconv.Atoi(turnCountStr)
			if err != nil {
				return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid turn count"})
			}

			gameRecord, err := gameGatewayImpl.FindLatest()
			if err != nil {
				log.Fatal(err)
			}

			turnRecord, err := turnGatewayImpl.FindForGameIDAndTurnCount(int(gameRecord.GetID()), turnCount)
			if err != nil {
				log.Fatal(err)
			}

			if err != nil {
				log.Fatal(err)
			}

			squaresRecord, err := squareGatewayImpl.FindForTurnID(int(turnRecord.GetID()))
			if err != nil {
				log.Fatal(err)
			}

			board := [8][8]int{}
			for _, square := range squaresRecord {
				board[square.GetY()][square.GetX()] = int(square.GetDisc())
			}

			res := struct {
				TurnCount  int       `json:"turnCount"`
				Board      [8][8]int `json:"board"`
				NextDisc   int       `json:"nextDisc"`
				WinnerDisc int       `json:"winnerDisc"`
			}{
				TurnCount:  int(turnRecord.GetTurnCount()),
				Board:      board,
				NextDisc:   int(turnRecord.GetNextDisc()),
				WinnerDisc: 0,
			}

			return c.JSON(http.StatusOK, res)
		})
	}
}
