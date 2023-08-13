package main

import (
	"context"
	"encoding/json"
	"log"

	"github.com/gofiber/fiber/v2"
	"sahabatgolkar.com/database"
	"sahabatgolkar.com/routes"
)

type Survey struct {
	Name string
	Umur int
}

func main() {
	ctx := context.TODO()
	app := fiber.New(fiber.Config{
		JSONEncoder: json.Marshal,
		JSONDecoder: json.Unmarshal,
	})

	client, err := database.ConnectMongoDB()
	if err != nil {
		panic(err)
	}

	defer func() {
		if err := client.Disconnect(ctx); err != nil {
			panic(err)
		}
	}()

	routes.InitRoutes(app)
	log.Fatal(app.Listen(":8000"))
}
