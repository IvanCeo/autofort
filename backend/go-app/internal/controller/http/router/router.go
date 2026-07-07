package router

import (
	"autofort/internal/controller/http"
	"os"

	"github.com/gofiber/swagger"

	"github.com/gofiber/fiber/v2"
	auth "github.com/gofiber/fiber/v2/middleware/basicauth"
)

func Route(app *fiber.App, handler *http.Handler) {
	app.Get("/swagger/*", swagger.HandlerDefault)
	api := app.Group("/api")
	api.Use(auth.New(auth.Config{
		Authorizer: func(name, pass string) bool {
			return name == os.Getenv("ADMIN_USER") && pass == os.Getenv("ADMIN_PASS")
		},
	}))

	// api.Post("/signin", handler.SignIn) нет пока ручки, может не надо будет
	api.Post("/customers", handler.CreateCustomerHandle)
	api.Post("/customers/:id/vehicles", handler.AddVehicleToCustomerHandle)
	api.Get("/customers", handler.ListCustomersHandle)
	api.Get("/customers/:id", handler.GetCustomerHandle)
	api.Patch("/customers/:id", handler.PatchCustomerHandle)

	api.Get("/vehicle-types", handler.GetVehicleTypesHandle)
	api.Patch("/vehicles/:id", handler.PatchVehicleHandle)

	api.Get("/vehicles/:id/work-order", handler.DownloadWorkOrderHandle)
}
