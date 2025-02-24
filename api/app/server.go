package main

import (
	othello "api/generated"
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"strconv"
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

	e.Use(middleware.Recover())

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

	e.POST("/api/games/latest/turns", func(c echo.Context) error {
		var turnReq struct {
			TurnCount int                      `json:"turnCount"`
			Move      othello.CreateMoveParams `json:"move"`
		}

		if err := c.Bind(&turnReq); err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request"})
		}
		log.Printf("turnReq: %+v", turnReq)

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
		fmt.Printf("latestGame: %+v\n", 1)
		latestGame, err := queries.GetLatestGame(ctx)
		if err != nil {
			log.Fatal(err)
			tx.Rollback()
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to get latest game"})
		}

		prevTurnCount := turnCount - 1
		fmt.Printf("get turn: %+v\n", 1)
		turn, err := queries.GetTurnByGameIDAndTurnCount(ctx, othello.GetTurnByGameIDAndTurnCountParams{
			GameID:    latestGame.ID,
			TurnCount: int32(prevTurnCount),
		})
		if err != nil {
			log.Fatal(err)
			tx.Rollback()
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to get turn"})
		}

		fmt.Printf("get square: %+v\n", 1)
		squares, err := queries.GetSquaresByTurnID(ctx, turn.ID)
		if err != nil {
			log.Fatal(err)
			tx.Rollback()
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to get squares"})
		}

		board := [8][8]int{}
		for _, square := range squares {
			board[square.Y][square.X] = int(square.Disc)
		}

		// 石を置く
		board[y][x] = disc

		// ターンを保存する
		var nextDisc int
		if disc == BLACK {
			nextDisc = WHITE
		} else {
			nextDisc = BLACK
		}

		fmt.Printf("insert turn: %+v\n", 1)
		fmt.Printf("game_id: %+v\n", latestGame.ID)
		fmt.Printf("turnCount: %+v\n", turnCount)
		fmt.Printf("nextDisc: %+v\n", nextDisc)
		insertTurnRes, err := queries.CreateTurn(ctx, othello.CreateTurnParams{
			GameID:    int32(latestGame.ID),
			TurnCount: int32(turnCount),
			NextDisc:  sql.NullInt32{Int32: int32(nextDisc), Valid: true},
		})
		// fmt.Println(err)
		fmt.Printf("insertTurnRes: %+v\n", insertTurnRes)
		// _, err = tx.Exec("insert into turns (game_id, turn_count, next_disc) values (?, ?, ?)", latestGame.ID, turnCount, nextDisc)
		// if err != nil {
		// 	log.Fatal(err)
		// 	tx.Rollback()
		// 	return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to insert turn"})
		// }


		turnID, err := insertTurnRes.LastInsertId()
		if err != nil {
			log.Fatal(err)
			tx.Rollback()
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to insert turn"})
		}

		fmt.Printf("insert square: %+v\n", 1)
		for y, line := range board {
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
					return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to insert square"})
				}
			}
		}

		fmt.Printf("insert move: %+v\n", 1)
		_, err = queries.CreateMove(ctx, othello.CreateMoveParams{
			TurnID: int32(turnID),
			X:      int32(x),
			Y:      int32(y),
			Disc:   int32(disc),
		})
		if err != nil {
			log.Fatal(err)
			tx.Rollback()
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to insert move"})
		}
		
		tx.Commit()
		return c.JSON(http.StatusOK, turnReq)
	})

	e.GET("/api/games/latest/turns/:turnCount", func(c echo.Context) error {
		turnCountStr := c.Param("turnCount")
		turnCount, err := strconv.Atoi(turnCountStr)
		if err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid turn count"})
		}

		latestGame, err := queries.GetLatestGame(ctx)
		if err != nil {
			log.Fatal(err)
		}

		turn, err := queries.GetTurnByGameIDAndTurnCount(ctx, othello.GetTurnByGameIDAndTurnCountParams{
			GameID:    latestGame.ID,
			TurnCount: int32(turnCount),
		})

		if err != nil {
			log.Fatal(err)
		}

		squares, err := queries.GetSquaresByTurnID(ctx, turn.ID)
		if err != nil {
			log.Fatal(err)
		}

		board := [8][8]int{}
		for _, square := range squares {
			board[square.Y][square.X] = int(square.Disc)
		}

		res := struct {
			TurnCount  int       `json:"turnCount"`
			Board      [8][8]int `json:"board"`
			NextDisc   int       `json:"nextDisc"`
			WinnerDisc int       `json:"winnerDisc"`
		}{
			TurnCount:  int(turn.TurnCount),
			Board:      board,
			NextDisc:   int(turn.NextDisc.Int32),
			WinnerDisc: 0,
		}

		return c.JSON(http.StatusOK, res)
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
