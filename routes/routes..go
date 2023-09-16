package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"sahabatgolkar.com/controllers"
	"sahabatgolkar.com/middlewares"
)

func InitRoutes(app *fiber.App) {
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowMethods: "GET,POST,HEAD,PUT,DELETE,PATCH",
	}))
	api := app.Group("/api")

	v1 := api.Group("/v1")
	v1.Post("/test", controllers.Test)
	v1.Post("/auth/local", controllers.Login)
	v1.Get("/auth/create-admin", controllers.CreateAdmin)

	auth := v1.Group("", middlewares.AuthHandler)
	auth.Get("/assets", controllers.ListAssets)
	auth.Get("/assets/pull/:uid", controllers.PullData)
	auth.Get("/assets/form/:uid", controllers.ReadForm)
	auth.Get("/assets/total/:uid", controllers.TotalDetail)
	auth.Post("/assets/data/:uid", controllers.ReadDetail)
}
