// @title           Autofort API
// @version         1.0
// @description     Autofort backend API
// @BasePath        /api
package main

import (
	"autofort/internal/adapter/postgres"
	"autofort/internal/adapter/wordorder"
	h "autofort/internal/controller/http"
	"autofort/internal/controller/http/router"
	"autofort/internal/migrate"
	"autofort/internal/usecase"
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "autofort/cmd/api/docs"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"

	"github.com/joho/godotenv"
)

func main() {

	if err := godotenv.Load("deploy/.env"); err != nil {
		log.Println("no .env file loaded:", err)
	}

	cfgConfig := migrate.LoadConfigFromEnv()

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	err := migrate.Up(ctx, cfgConfig)
	if err != nil {
		log.Fatal(err)
	}

	dsn, err := migrate.BuildPostgresDSN(cfgConfig)
	if err != nil {
		log.Fatal(err)
	}

	postgres, err := postgres.NewPostgres(ctx, dsn)
	if err != nil {
		log.Fatal(err)
	}
	defer postgres.Close()

	app := fiber.New()

	app.Use(logger.New(logger.Config{
		Format: "[${time}] ${ip} ${method} ${path} ${status} - ${latency}\n",
	}))
	app.Use(cors.New(cors.ConfigDefault))

	workorderClient, err := wordorder.NewWorkOrderClient(ctx, 10*time.Second)

	srv := usecase.NewService(postgres, postgres, postgres, workorderClient)

	handler := h.NewHandler(srv)

	router.Route(app, handler)

	port := ":" + os.Getenv("APP_PORT")

	go func() {
		if err := app.Listen(port); err != nil {
			log.Printf("fiber stopped:%v", err)
		}
	}()

	log.Printf("server started on %s", port)

	<-ctx.Done()

	log.Println("shutting down server...")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := app.ShutdownWithContext(shutdownCtx); err != nil {
		log.Printf("server shutdown failed: %v", err)
	}

	log.Println("server exited gracefully")
}
