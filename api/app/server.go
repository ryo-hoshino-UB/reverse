package main

import (
	othello "api/generated"
	"context"
	"database/sql"
	"log"
	"net/http"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

const (
	EMPTY = 0
	BLACK = 1
	WHITE = 2
)

var (
	INITIAL_BOARD = [8][8]int{
		{EMPTY, EMPTY, EMPTY, EMPTY, EMPTY, EMPTY, EMPTY, EMPTY},
		{EMPTY, EMPTY, EMPTY, EMPTY, EMPTY, EMPTY, EMPTY, EMPTY},
		{EMPTY, EMPTY, EMPTY, EMPTY, EMPTY, EMPTY, EMPTY, EMPTY},
		{EMPTY, EMPTY, EMPTY, BLACK, WHITE, EMPTY, EMPTY, EMPTY},
		{EMPTY, EMPTY, EMPTY, WHITE, BLACK, EMPTY, EMPTY, EMPTY},
		{EMPTY, EMPTY, EMPTY, EMPTY, EMPTY, EMPTY, EMPTY, EMPTY},
		{EMPTY, EMPTY, EMPTY, EMPTY, EMPTY, EMPTY, EMPTY, EMPTY},
		{EMPTY, EMPTY, EMPTY, EMPTY, EMPTY, EMPTY, EMPTY, EMPTY},
	}
	SQUARE_COUNT = len(INITIAL_BOARD) * len(INITIAL_BOARD[0])
)

func main() {
	ctx := context.Background()

	db := connectDB()
	defer db.Close()

	queries := othello.New(db)

	e := echo.New()
	e.Use(middleware.BodyDump(func(c echo.Context, reqBody, resBody []byte) {
		log.Printf("%s %s %d %s %s", c.Request().Method, c.Request().RequestURI, c.Response().Status, reqBody, resBody)
	}))

	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"http://localhost:5001"}}))

	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})

	// game開始時に叩くAPI
	// games tableに新規レコードを追加する
	// turns tableに新規レコード(初期盤面)を追加する
	// squares tableに初期盤面の各マス目の情報を追加する
	e.POST("/api/games", func(c echo.Context) error {
		var startedAt = time.Now()

		tx, err := db.BeginTx(ctx, nil)
		if err != nil {
			log.Fatal(err)
		}

		insertGameRes, err := queries.CreateGame(ctx)
		if err != nil {
			log.Fatal(err)
			tx.Rollback()
		}

		insertedGameID, err := insertGameRes.LastInsertId()
		if err != nil {
			log.Fatal(err)
			tx.Rollback()
		}

		insertTurnRes, err := queries.CreateTurn(ctx, othello.CreateTurnParams{
			GameID:    int32(insertedGameID),
			TurnCount: 0,
			NextDisc:  sql.NullInt32{Int32: BLACK, Valid: true},
		})

		turnID, err := insertTurnRes.LastInsertId()
		if err != nil {
			log.Fatal(err)
			tx.Rollback()
		}

		for y, line := range INITIAL_BOARD {
			for x, disc := range line {
				_, err := queries.CreateSquare(ctx, othello.CreateSquareParams{
					TurnID: int32(turnID),
					X:      int32(x),
					Y:      int32(y),
					Disc:   int32(disc),
				})

				if err != nil {
					log.Fatal(err)
					tx.Rollback()
				}
			}
			
		}

		if err != nil {
			log.Fatal(err)
			tx.Rollback()
		}

		err = tx.Commit()
		if err != nil {
			log.Fatal(err)
			tx.Rollback()
		}

		return c.JSON(http.StatusOK, startedAt)
	})

	e.Logger.Fatal(e.Start(":5002"))
}

func connectDB() *sql.DB {
	db, err := sql.Open("mysql", "othello_user:pass@tcp(localhost:3306)/othello?parseTime=true")
	if err != nil {
		log.Fatal(err)
	}

	return db
}
