package main

import (
	"api/dataaccess"
	"api/presentation"
	"context"
	"log"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	ctx := context.Background()

	db := dataaccess.ConnectDB()
	defer db.Close()

	e := echo.New()
	e.Use(middleware.BodyDump(func(c echo.Context, reqBody, resBody []byte) {
		log.Printf("%s %s %d %s %s", c.Request().Method, c.Request().RequestURI, c.Response().Status, reqBody, resBody)
	}))

	e.Use(middleware.Recover())

	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"http://localhost:5001"}}))

	presentation.GameRouter(ctx, db)(e)
	presentation.TurnRouter(ctx, db)(e)

	e.Logger.Fatal(e.Start(":5002"))
}
