package main

import (
	"api/dataaccess"
	othello "api/generated"
	"context"
	"database/sql"
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
	INITIAL_BOARD = [][]int{
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

	var gameGatewayImpl = dataaccess.GameGatewayImpl{
		Context: ctx,
		Queries: queries,
	}
	var turnGatewayImpl = dataaccess.TurnGateway{
		Context: ctx,
		Queries: queries,
	}
	var moveGatewayImpl = dataaccess.MoveGatewayImpl{
		Context: ctx,
		Queries: queries,
	}
	var squareGatewayImpl = dataaccess.SquareGatewayImpl{
		Context: ctx,
		Queries: queries,
	}

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

		gameRecord, err := gameGatewayImpl.Insert()
		if err != nil {
			log.Fatal(err)
			tx.Rollback()
		}

		insertedGameID := gameRecord.GetID()

		turnRecord, err := turnGatewayImpl.Insert(insertedGameID, 0, BLACK)

		err = squareGatewayImpl.InsertAll(int(turnRecord.GetID()), INITIAL_BOARD)
		if err != nil {
			log.Fatal(err)
			tx.Rollback()
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
		if disc == BLACK {
			nextDisc = WHITE
		} else {
			nextDisc = BLACK
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

	e.GET("/api/games/latest/turns/:turnCount", func(c echo.Context) error {
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

	e.Logger.Fatal(e.Start(":5002"))
}

func connectDB() *sql.DB {
	db, err := sql.Open("mysql", "othello_user:pass@tcp(localhost:3306)/othello?parseTime=true")
	if err != nil {
		log.Fatal(err)
	}

	return db
}
