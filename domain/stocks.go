package domain

import (
	"bitbucket/api_format/db"
	"bitbucket/api_format/models"
	"log"

	"gorm.io/gorm"
)

// --------------------- CREATE STOCK (TRANSACTIONAL) ---------------------
func CreateStockTx(tx *gorm.DB, stock models.Stocks) (int, error) {
	if err := tx.Table("stocks").Create(&stock).Error; err != nil {
		log.Println("Error creating stock:", err)
		return 0, err
	}
	return stock.StockID, nil
}

// --------------------- GET STOCK BY REGISTRATION (NON-transactional) ---------------------
func GetStockByRegistration(registration string) int {
	var stock models.Stocks
	if err := db.DB.Table("stocks").
		Select("stock_id").
		Where("registration = ?", registration).
		First(&stock).Error; err != nil {
		return 0
	}
	return stock.StockID
}

// --------------------- GET STOCK BY REGISTRATION AND DATE (TRANSACTIONAL) ---------------------
func GetStockByRegistrationAndDateTx(tx *gorm.DB, datePurchased, registration string) int {
	var stockID int
	if err := tx.Table("stocks").
		Select("stock_id").
		Where("date_purchased = ? AND registration = ? and status = 1", datePurchased, registration).
		Limit(1).
		Scan(&stockID).Error; err != nil {
		return 0
	}
	return stockID
}

// --------------------- GET STOCK BY ID ---------------------
// func GetStockByID(id int) models.Stocks {
// 	var stock models.Stocks
// 	if err := db.DB.Table("stocks").
// 		Where("stock_id = ? and status = 1", id).
// 		First(&stock).Error; err != nil {
// 		log.Println("Stock not found:", err)
// 	}
// 	return stock
// }

func GetStockByID(id int) (models.Stocks, error) {
	var stock models.Stocks

	err := db.DB.Model(&models.Stocks{}).
		Where("stock_id = ? AND status = 1", id).
		First(&stock).Error

	if err != nil {
		log.Println("Error fetching stock:", err)
		return stock, err
	}

	return stock, nil
}

func GetStockByNo(id string) (models.Stocks, error) {
	var stock models.Stocks

	err := db.DB.Model(&models.Stocks{}).
		Where("stock_number = ?", id).
		First(&stock).Error

	if err != nil {
		log.Println("Error fetching stock:", err)
		return stock, err
	}

	return stock, nil
}

// --------------------- GET STOCK BY ID ---------------------
func GetStockByStatus(id int) []models.Stocks {
	var stock []models.Stocks
	if err := db.DB.Table("stocks").
		Where("status = ?", id).
		Find(&stock).Error; err != nil {
		log.Println("Stock not found:", err)
	}
	return stock
}

// --------------------- GET ALL STOCKS ---------------------

func GetAllStocks() []models.Stocks {
	var stocks []models.Stocks

	if err := db.DB.Table("stocks").
		// Where("status = 1").
		Find(&stocks).Error; err != nil {
		log.Println("Error retrieving active stocks:", err)
	}

	return stocks
}

func GetAllStocksByTid(id int) []models.Stocks {
	var stocks []models.Stocks

	if err := db.DB.Table("stocks").
		Where("tenant_info_id = ? and status != 0", id).
		Find(&stocks).Error; err != nil {
		log.Println("Error retrieving active stocks:", err)
	}

	return stocks
}

func GetAllStocksByTidwithstocks(id int) []models.Stocks {
	var stocks []models.Stocks

	if err := db.DB.Table("stocks").
		Where("tenant_info_id = ? and status = 1", id).
		Find(&stocks).Error; err != nil {
		log.Println("Error retrieving active stocks:", err)
	}

	return stocks
}

// func GetAllStocks() []models.Stocks {
// 	var stocks []models.Stocks
// 	if err := db.DB.Table("stocks").Where("status = 1").Find(&stocks).Error; err != nil {
// 		log.Println("Error retrieving stocks:", err)
// 	}
// 	return stocks
// }

// --------------------- DELETE STOCK ---------------------

func DeleteStockByID(id int) (bool, string) {
	tx := db.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Update status to 0 instead of deleting the record
	if err := tx.Table("stocks").
		Where("stock_id = ?", id).
		Update("status", 0).Error; err != nil {
		tx.Rollback()
		log.Println("Error updating stock status:", err)
		return false, "Failed to update stock status"
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		log.Println("Commit error while updating stock status:", err)
		return false, "Transaction commit failed"
	}

	return true, "Stock marked as inactive successfully"
}

// func DeleteStockByID(id int) (bool, string) {
// 	tx := db.DB.Begin()
// 	defer func() {
// 		if r := recover(); r != nil {
// 			tx.Rollback()
// 		}
// 	}()

// 	if err := tx.Table("stocks").
// 		Where("stock_id = ?", id).
// 		Delete(nil).Error; err != nil {
// 		tx.Rollback()
// 		log.Println("Error deleting stock:", err)
// 		return false, "Failed to delete stock"
// 	}

// 	if err := tx.Commit().Error; err != nil {
// 		tx.Rollback()
// 		log.Println("Commit error while deleting stock:", err)
// 		return false, "Transaction commit failed"
// 	}

// 	return true, ""
// }
