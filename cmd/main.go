package main

import (
	"context"
	"log"
	"log/slog"
	"os"

	"github.com/cabrn3t/todo-list-server/config"
	"github.com/cabrn3t/todo-list-server/internal/database"
	"github.com/cabrn3t/todo-list-server/internal/handlers"
	"github.com/cabrn3t/todo-list-server/internal/repositories"
	"github.com/cabrn3t/todo-list-server/internal/routes"
	"github.com/gofiber/fiber/v3"
)

func main() {
	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, nil)))

	cfg := config.GetConfig()
	slog.Info("config read successfully")

	conn, err := database.Connect(context.Background(), *cfg, 5)

	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	app := fiber.New()

	repo := repositories.New(conn)
	taskHandler := handlers.NewTaskHandler(*repo)
	taskRoutes := routes.NewTaskRoutes(*taskHandler)
	taskRoutes.RegisterRoutes(app)

	slog.Info("server started successfully")
	log.Fatal(app.Listen(":3000"))
}
