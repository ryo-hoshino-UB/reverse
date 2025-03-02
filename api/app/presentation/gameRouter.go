package presentation

import (
	"api/dataaccess"
	othello "api/generated"
	"api/lib"
	"context"
	"database/sql"
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
)

func GameRouter(ctx context.Context, db *sql.DB) func(e *echo.Echo) {
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
		squareGatewayImpl := dataaccess.SquareGatewayImpl{
			Context: ctx,
			Queries: queries,
		}

		// ゲームに関連するAPI群
		games := e.Group("/api/games")

		// ゲーム作成API
		games.POST("", func(c echo.Context) error {
			return CreateGame(c, ctx, queries, gameGatewayImpl, turnGatewayImpl, squareGatewayImpl)
		})
	}
}

func CreateGame(c echo.Context, ctx context.Context, queries *othello.Queries,
	gameGatewayImpl dataaccess.GameGatewayImpl,
	turnGatewayImpl dataaccess.TurnGateway,
	squareGatewayImpl dataaccess.SquareGatewayImpl) error {

	db := dataaccess.ConnectDB()
	defer db.Close()
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		log.Fatal(err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to start transaction"})
	}

	gameRecord, err := gameGatewayImpl.Insert()
	if err != nil {
		log.Fatal(err)
		tx.Rollback()
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to insert game"})
	}

	insertedGameID := gameRecord.GetID()

	turnRecord, err := turnGatewayImpl.Insert(insertedGameID, 0, lib.BLACK)
	if err != nil {
		log.Fatal(err)
		tx.Rollback()
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to insert turn"})
	}

	err = squareGatewayImpl.InsertAll(int(turnRecord.GetID()), lib.INITIAL_BOARD)
	if err != nil {
		log.Fatal(err)
		tx.Rollback()
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to insert squares"})
	}

	err = tx.Commit()
	if err != nil {
		log.Fatal(err)
		tx.Rollback()
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to commit transaction"})
	}

	return c.JSON(http.StatusOK, gameRecord.Game.StartedAt)
}
