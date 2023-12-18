package routes

import (
	"github.com/gofiber/fiber/v2"

	auth "github.com/fpmi-hci-2023/project13b-auth/api/auth/v1/gen"
	"github.com/fpmi-hci-2023/project13b-auth/internal/server/handlers"
)

func AuthRouter(app fiber.Router, authClient auth.AuthServiceClient) {

	handler := handlers.NewAuthHandler(app, authClient)

	app.Post("/login", handler.Login)

	app.Post("/validate", handler.Validate)

	app.Post("/logout", handler.Logout)

	app.Post("/reg", handler.Registration)

	app.Get("/i", handler.Info)

}
