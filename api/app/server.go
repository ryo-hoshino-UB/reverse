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

	e.POST("/api/games", func(c echo.Context) error {
		var startedAt = time.Now()
		log.Println("Started at: ", startedAt)

		tx, err := db.BeginTx(ctx, nil)
		if err != nil {
			log.Fatal(err)
		}

		err = queries.CreateGame(ctx)
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
