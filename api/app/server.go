package main

import (
	infrastructure "api/infrastructure"
	"api/presentation"
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	ctx := context.Background()

	db := infrastructure.ConnectDB()
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

	startWithGracefulShutdown(e, ":5002", 10*time.Second)
}

func startWithGracefulShutdown(e *echo.Echo, address string, timeout time.Duration) {
	// See: https://echo.labstack.com/docs/cookbook/graceful-shutdown

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	go func() {
		if err := e.Start(address); err != nil && err != http.ErrServerClosed {
			e.Logger.Fatal("shutting down the server")
		}
	}()

	// Wait for interrupt signal to gracefully shut down the server with a timeout.
	<-ctx.Done()

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	if err := e.Shutdown(ctx); err != nil {
		e.Logger.Fatal(err)
	}
}
