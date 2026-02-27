package routes

import (
	"fmt"
	"net/http"

	"bitbucket/api_format/controllers"

	"github.com/gofiber/fiber/v2"
)

func Setup(app *fiber.App) {

	environments := []string{"admin", "dev", "live", "prod", "autodev", "autoprod"}

	for _, env := range environments {
		group := app.Group(fmt.Sprintf("/%s", env))
		v1 := group.Group("/v1")
		V1routes(v1)
		// db.CloseDB()
		// v2 := group.Group("/app")
		// V1routes(v2)
	}
}

func V1routes(group fiber.Router) {

	group.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"code":    http.StatusOK,
			"message": "Welcome to Dexworld",
			"status":  true,
		})
	})

	hpi := group.Group("/hpi")
	hpi.Get("/health", controllers.HPIHealth)
	hpi.Get("/products", controllers.HPIProducts)
	hpi.Post("/enquiry", controllers.HPIEnquiry)
	hpi.Post("/check", controllers.HPICheckVehicle)

	bankinvoice := group.Group("/bankinvoice")
	bankinvoice.Get("/health", controllers.BankInvoiceHealth)
	bankinvoice.Get("/history", controllers.GetBankInvoiceHistory)
	bankinvoice.Get("/history/tenant/:tenantid", controllers.GetBankInvoiceHistoryByTenantID)
	bankinvoice.Get("/history/:id", controllers.GetBankInvoiceHistoryByID)
	bankinvoice.Post("/convert", controllers.ConvertBankInvoicePDFToCSV)

	tenant := group.Group("/tenant")
	tenant.Post("/create", controllers.Createtenant)
	tenant.Put("/update", controllers.Updatetenant)
	tenant.Get("/getbyid", controllers.Gettenantbyid)
	tenant.Delete("delete", controllers.Deletetenantbyid)
	tenant.Get("/getall", controllers.Getalltenants)
	tenant.Get("/getallbyorgname", controllers.Getallbytenant)
	tenant.Get("/getavailableapps", controllers.Getavailableapps)
	tenant.Get("/gettenantapps", controllers.Gettenantapps)
	tenant.Get("/getallorganizations", controllers.Getallorganizations)

	stocks := group.Group("/stocks")
	stocks.Post("/create", controllers.CreateStock)
	stocks.Put("/update", controllers.UpdateStock)
	stocks.Get("/getbyid", controllers.GetStockByID)
	stocks.Get("/get/stockno", controllers.GetStockByNo)
	stocks.Get("/get/status", controllers.GetStockByStatus)
	stocks.Get("/getall", controllers.GetAllStocks)
	stocks.Get("/getall/tid", controllers.GetAllStocksByTid)
	stocks.Delete("/deletebyid", controllers.DeleteStockByID)

}
