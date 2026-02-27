package controllers

import (
	"bitbucket/api_format/db"
	"bitbucket/api_format/domain"
	"bitbucket/api_format/models"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

// --------------------- HELPERS ---------------------
func jsonError(c *fiber.Ctx, code int, message string, err error) error {
	resp := fiber.Map{"code": code, "message": message, "status": false}
	if err != nil {
		resp["error"] = err.Error()
	}
	return c.Status(code).JSON(resp)
}

func jsonSuccess(c *fiber.Ctx, code int, message string, details any) error {
	return c.Status(code).JSON(fiber.Map{"code": code, "message": message, "status": true, "details": details})
}

// // --------------------- GENERATE STOCK NUMBER ---------------------
// func GenerateStockNumber(tx *gorm.DB) (string, error) {
//     var lastSeq int
//     year := time.Now().Year()

//     res := tx.Exec("UPDATE stock_sequence SET last_seq = last_seq + 1 WHERE year = ?", year)
//     if res.RowsAffected == 0 {
//         if err := tx.Exec("INSERT INTO stock_sequence(year, last_seq) VALUES (?, 1)", year).Error; err != nil {
//             return "", err
//         }
//         lastSeq = 1
//     } else {
//         if err := tx.Raw("SELECT last_seq FROM stock_sequence WHERE year = ?", year).Scan(&lastSeq).Error; err != nil {
//             return "", err
//         }
//     }

//     return fmt.Sprintf("STK-%d_%03d", year, lastSeq), nil
// }package controllers

// --------------------- CREATE STOCK (TRANSACTIONAL & SEQUENTIAL STOCK NUMBER) ---------------------
func CreateStock(c *fiber.Ctx) error {

	id, _ := strconv.Atoi(c.Query("tenant_info_id"))
	var stock models.Stocks

	// Parse request payload
	if err := c.BodyParser(&stock); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"code":    http.StatusBadRequest,
			"message": "Invalid request payload",
			"status":  false,
			"error":   err.Error(),
		})
	}

	// Start transaction
	tx := db.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Check if stock already exists for same registration + date
	if existingID := domain.GetStockByRegistrationAndDateTx(tx, *stock.DatePurchased, *stock.Registration); existingID != 0 {
		tx.Rollback()
		return c.Status(http.StatusConflict).JSON(fiber.Map{
			"code":    http.StatusConflict,
			"message": "Stock already exists!",
			"status":  false,
		})
	}

	// Generate sequential stock number safely
	stockNumber, err := GenerateStockNumber(tx)
	if err != nil {
		tx.Rollback()
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"code":    http.StatusInternalServerError,
			"message": "Failed to generate stock number",
			"status":  false,
			"error":   err.Error(),
		})
	}
	stock.StockNumber = &stockNumber
	stock.TenantInfoID = &id

	// Create stock record
	newID, err := domain.CreateStockTx(tx, stock)
	if err != nil {
		tx.Rollback()
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"code":    http.StatusInternalServerError,
			"message": "Failed to create stock",
			"status":  false,
			"error":   err.Error(),
		})
	}
	stock.StockID = newID

	// Commit transaction
	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"code":    http.StatusInternalServerError,
			"message": "Transaction commit failed",
			"status":  false,
			"error":   err.Error(),
		})
	}

	// Return created stock details
	return c.Status(http.StatusCreated).JSON(fiber.Map{
		"code":    http.StatusCreated,
		"message": "Stock created successfully",
		"status":  true,
		"details": stock,
	})
}

// --------------------- HELPER: GENERATE SEQUENTIAL STOCK NUMBER ---------------------
func GenerateStockNumber(tx *gorm.DB) (string, error) {
	year := time.Now().Year()
	var maxSeq int

	// Get max sequence from existing stock numbers
	query := `
        SELECT IFNULL(MAX(CAST(SUBSTRING_INDEX(stock_number, '_', -1) AS UNSIGNED)), 0)
        FROM stocks
        WHERE stock_number LIKE ?`
	pattern := fmt.Sprintf("STK-%d_%%", year)
	if err := tx.Raw(query, pattern).Scan(&maxSeq).Error; err != nil {
		return "", err
	}

	nextSeq := maxSeq + 1

	// Update stock_sequence table to remain in sync
	if err := tx.Exec(`
        INSERT INTO stock_sequence(year, last_seq)
        VALUES (?, ?)
        ON DUPLICATE KEY UPDATE last_seq = ?`,
		year, nextSeq, nextSeq).Error; err != nil {
		return "", err
	}

	return fmt.Sprintf("STK-%d_%03d", year, nextSeq), nil
}

// // --------------------- UPDATE STOCK ---------------------
func UpdateStock(c *fiber.Ctx) error {
	var payload map[string]interface{}
	if err := c.BodyParser(&payload); err != nil {
		return jsonError(c, http.StatusBadRequest, "Invalid request payload", err)
	}

	registration, ok := payload["stock_id"]
	if !ok {
		return jsonError(c, http.StatusBadRequest, "`stock_id` is required", nil)
	}

	tx := db.DB.Begin()
	if err := tx.Table("stocks").Where("stock_id = ?", registration).Updates(payload).Error; err != nil {
		tx.Rollback()
		return jsonError(c, http.StatusInternalServerError, "Failed to update stock record", err)
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return jsonError(c, http.StatusInternalServerError, "Transaction commit failed", err)
	}

	return jsonSuccess(c, http.StatusOK, "Stock updated successfully", nil)
}

// --------------------- UPDATE STOCK ---------------------
// func UpdateStock(c *fiber.Ctx) error {
// 	var payload map[string]interface{}
// 	if err := c.BodyParser(&payload); err != nil {
// 		return jsonError(c, http.StatusBadRequest, "Invalid request payload", err)
// 	}

// 	// Ensure at least one identifier is provided
// 	registration, regOk := payload["registration"]
// 	stockNumber, stkOk := payload["stock_number"]

// 	if !regOk && !stkOk {
// 		return jsonError(c, http.StatusBadRequest, "Either 'registration' or 'stock_number' is required", nil)
// 	}

// 	tx := db.DB.Begin()
// 	var query *gorm.DB

// 	// Use whichever key is available
// 	if regOk {
// 		query = tx.Table("stocks").Where("registration = ?", registration)
// 	} else {
// 		query = tx.Table("stocks").Where("stock_number = ?", stockNumber)
// 	}

// 	// Perform the update
// 	if err := query.Updates(payload).Error; err != nil {
// 		tx.Rollback()
// 		return jsonError(c, http.StatusInternalServerError, "Failed to update stock record", err)
// 	}

// 	// Commit transaction
// 	if err := tx.Commit().Error; err != nil {
// 		tx.Rollback()
// 		return jsonError(c, http.StatusInternalServerError, "Transaction commit failed", err)
// 	}

// 	return jsonSuccess(c, http.StatusOK, "Stock updated successfully", nil)
// }

// --------------------- GET STOCK BY ID ---------------------

func GetStockByID(c *fiber.Ctx) error {
	id, _ := strconv.Atoi(c.Query("stock_id"))
	stock, err := domain.GetStockByID(id)
	if err != nil {
		return jsonError(c, http.StatusNotFound, "Stock not found", err)
	}
	if stock.StockID == 0 {
		return jsonError(c, http.StatusNotFound, "Stock not found", nil)
	}
	return jsonSuccess(c, http.StatusOK, "Success", stock)
}

func GetStockByNo(c *fiber.Ctx) error {
	id := c.Query("stock_no")
	stock, err := domain.GetStockByNo(id)
	if err != nil {
		return jsonError(c, http.StatusNotFound, "Stock not found", err)
	}
	if stock.StockID == 0 {
		return jsonError(c, http.StatusNotFound, "Stock not found", nil)
	}
	return jsonSuccess(c, http.StatusOK, "Success", stock)
}

func GetStockByStatus(c *fiber.Ctx) error {
	id, _ := strconv.Atoi(c.Query("status"))
	stock := domain.GetStockByStatus(id)

	if len(stock) == 0 {
		return jsonError(c, http.StatusNotFound, "Stock not found", nil)
	}
	return jsonSuccess(c, http.StatusOK, "Success", stock)
}

// --------------------- GET ALL STOCKS ---------------------
func GetAllStocks(c *fiber.Ctx) error {
	stocks := domain.GetAllStocks()
	if len(stocks) == 0 {
		return jsonSuccess(c, http.StatusOK, "No stocks found", []models.Stocks{})
	}
	return jsonSuccess(c, http.StatusOK, "Success", stocks)
}

func GetAllStocksByTid(c *fiber.Ctx) error {

	id, _ := strconv.Atoi(c.Query("tenant_info_id"))

	stocks := domain.GetAllStocksByTid(id)
	if len(stocks) == 0 {
		return jsonSuccess(c, http.StatusOK, "No stocks found", []models.Stocks{})
	}
	return jsonSuccess(c, http.StatusOK, "Success", stocks)
}

// --------------------- DELETE STOCK ---------------------
func DeleteStockByID(c *fiber.Ctx) error {
	id, _ := strconv.Atoi(c.Query("stock_id"))
	success, msg := domain.DeleteStockByID(id)
	if success {
		return jsonSuccess(c, http.StatusOK, "Stock record deleted successfully", nil)
	}
	return jsonError(c, http.StatusBadRequest, msg, nil)
}
