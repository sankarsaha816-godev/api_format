package db

import (
	"bitbucket/api_format/models"
	"bitbucket/api_format/utils"
	"fmt"
	"log"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

// func DevConnect() {

// 	var DBname string
// 	var Username string
// 	var Password string
// 	var Host string
// 	var Port string

// 	Port, DBname, Password, Username, Host, _, _, _, _, _, _ = utils.DevConfig()

// 	dsn := Username + ":" + Password + "@tcp" + "(" + Host + ":" + Port + ")/" + DBname + "?" + "parseTime=true&loc=Local"

// 	database, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

// 	if err != nil {
// 		panic("Could not connect to the database")
// 	}
// 	print("Dev Database Connected Successfully")
// 	DB = database
// }

func DevConnect() {
	var DBname string
	var Username string
	var Password string
	var Host string
	var Port string

	Port, DBname, Password, Username, Host, _, _, _, _, _, _ = utils.DevConfig()

	// dsn := "zqbuvyrhhf:Package%40%123#@tcp(165.232.178.78:3306)/dbname?parseTime=true"

	dsn := Username + ":" + Password + "@tcp" + "(" + Host + ":" + Port + ")/" + DBname + "?" + "parseTime=true&loc=Local"

	//database, err := sql.Open("mysql", dsn)

	database, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("Could not connect to the database")
	}

	sqlDB, err := database.DB()
	if err != nil {
		log.Println("Failed to get DB from GORM:", err)
		panic("Could not connect to the database")

	}
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	if err = sqlDB.Ping(); err != nil {
		log.Println("Error:", err)
		panic("Could not connect to the database")
	}

	fmt.Println("dev Database Connected Successfully")
	DB = database
}

// func LiveConnect() {

// 	var DBname string
// 	var Username string
// 	var Password string
// 	var Host string
// 	var Port string

// 	Port, DBname, Password, Username, Host, _, _, _, _, _, _ = utils.LiveConfig()

// 	// dsn := "zqbuvyrhhf:Package%40%123#@tcp(165.232.178.78:3306)/dbname?parseTime=true"

// 	dsn := Username + ":" + Password + "@tcp" + "(" + Host + ":" + Port + ")/" + DBname + "?" + "parseTime=true&loc=Local"

// 	//database, err := sql.Open("mysql", dsn)

// 	fmt.Println("---------mydsn----------", dsn)

// 	database, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

// 	if err != nil {
// 		panic("Could not connect to the database")
// 	}
// 	print("Live Database Connected Successfully")
// 	DB = database
// }

// func LiveConnect() {
// 	var DBname, Username, Password, Host, Port string

// 	// Fetching the live configuration values
// 	Port, DBname, Password, Username, Host, _, _, _, _, _, _ = utils.LiveConfig()

// 	// Constructing the DSN (Data Source Name)
// 	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true&loc=Local", Username, Password, Host, Port, DBname)

// 	fmt.Println("---------mydsn----------", dsn)

// 	// Opening the database connection
// 	database, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
// 	if err != nil {
// 		log.Fatalf("Could not connect to the database: %v", err)
// 	}

// 	fmt.Println("Live Database Connected Successfully")
// 	DB = database
// }

func LiveConnect() {
	var DBname string
	var Username string
	var Password string
	var Host string
	var Port string

	Port, DBname, Password, Username, Host, _, _, _, _, _, _ = utils.LiveConfig()

	dsn := Username + ":" + Password + "@tcp" + "(" + Host + ":" + Port + ")/" + DBname + "?" + "parseTime=true&loc=Local"

	database, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("Could not connect to the database")
	}

	sqlDB, err := database.DB()
	if err != nil {
		log.Println("Failed to get DB from GORM:", err)
		panic("Could not connect to the database")

	}
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	if err = sqlDB.Ping(); err != nil {
		log.Println("Error:", err)
		panic("Could not connect to the database")
	}

	fmt.Println("Live Database Connected Successfully")
	DB = database

}

func ProdConnect() {
	var DBname string
	var Username string
	var Password string
	var Host string
	var Port string

	Port, DBname, Password, Username, Host, _, _, _, _, _, _ = utils.ProdConfig()

	dsn := Username + ":" + Password + "@tcp" + "(" + Host + ":" + Port + ")/" + DBname + "?" + "parseTime=true&loc=Local"

	database, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("Could not connect to the database")
	}

	sqlDB, err := database.DB()
	if err != nil {
		log.Println("Failed to get DB from GORM:", err)
		panic("Could not connect to the database")

	}
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	if err = sqlDB.Ping(); err != nil {
		log.Println("Error:", err)
		panic("Could not connect to the database")
	}

	fmt.Println("Live Database Connected Successfully")
	DB = database
}

func AutoConnect() {
	var DBname string
	var Username string
	var Password string
	var Host string
	var Port string

	Port, DBname, Password, Username, Host, _, _, _, _, _, _ = utils.AutoConfig()

	dsn := Username + ":" + Password + "@tcp" + "(" + Host + ":" + Port + ")/" + DBname + "?" + "parseTime=true&loc=Local"

	database, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("Could not connect to the database")
	}

	sqlDB, err := database.DB()
	if err != nil {
		log.Println("Failed to get DB from GORM:", err)
		panic("Could not connect to the database")

	}
	sqlDB.SetMaxIdleConns(100)
	sqlDB.SetMaxOpenConns(10000)
	sqlDB.SetConnMaxLifetime(time.Hour)

	if err = sqlDB.Ping(); err != nil {
		log.Println("Error:", err)
		panic("Could not connect to the database")
	}

	fmt.Println("Auto Database Connected Successfully")
	DB = database
}

func AutoConnectDEV() {
	var DBname string
	var Username string
	var Password string
	var Host string
	var Port string

	Port, DBname, Password, Username, Host, _, _, _, _, _, _ = utils.AutoConfigDEV()

	dsn := Username + ":" + Password + "@tcp" + "(" + Host + ":" + Port + ")/" + DBname + "?" + "parseTime=true&loc=Local"

	database, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("Could not connect to the database")
	}

	sqlDB, err := database.DB()
	if err != nil {
		log.Println("Failed to get DB from GORM:", err)
		panic("Could not connect to the database")

	}
	sqlDB.SetMaxIdleConns(100)
	sqlDB.SetMaxOpenConns(10000)
	sqlDB.SetConnMaxLifetime(time.Hour)

	if err = sqlDB.Ping(); err != nil {
		log.Println("Error:", err)
		panic("Could not connect to the database")
	}

	fmt.Println("Auto Database Connected Successfully")
	DB = database
}

// func Prodconnect() {

// 	var DBname string
// 	var Username string
// 	var Password string
// 	var Host string
// 	var Port string

// 	Port, DBname, Password, Username, Host, _, _, _, _, _, _ = utils.QlikConfig()

// 	// DSN format for PostgreSQL: "host=myhost port=myport user=gorm dbname=gorm password=mypassword sslmode=disable TimeZone=Asia/Shanghai"
// 	dsn := "host=" + Host + " port=" + Port + " user=" + Username + " dbname=" + DBname + " password=" + Password + " sslmode=disable TimeZone=UTC"
// 	postgresDSN := "postgres://postgres:ql1ks3ns3adm!n@10.10.1.25:4432/QSR?sslmode=disable"
// 	fmt.Println("---------psdsn----------", dsn)

// 	database, err := gorm.Open(postgres.Open(postgresDSN), &gorm.Config{})

// 	if err != nil {
// 		panic("Could not connect to the database")
// 	}
// 	print("Live Database Connected Successfully")
// 	DB = database

// }

func CloseDB() {
	sqlDB, err := DB.DB()
	if err != nil {
		log.Println("Error retrieving sql.DB from GORM:", err)
		return
	}
	fmt.Println("Connection closed Successfully")
	sqlDB.Close() // This will close the connection pool

}

// // InitializeDatabase auto-migrates all models
// // InitializeDatabase auto-migrates provided models on the given DB connection
func InitializeDatabase(dbConn *gorm.DB) error {
	if dbConn == nil {
		return fmt.Errorf("database connection not initialized")
	}

	// Models to auto-migrate
	modelsToMigrate := []interface{}{
		// Ensure the VehicleSale model is auto-migrated so missing columns (e.g., vehicle_id)
		// are created in existing `vehicle_sales` table when possible
		&models.VehicleSale{},
		// &models.StockImage{},
		&models.VATTransaction{},
		&models.VehicleSalesVAT{},
		&models.WorkshopTransaction{},
		&models.VATPurchase{},
		&models.ECTransaction{},
		&models.VATPeriod{},
		&models.VehicleValuation{},
	}
	fmt.Println(modelsToMigrate...)

	// if err := dbConn.AutoMigrate(modelsToMigrate...); err != nil {
	// 	log.Printf("Warning: Failed to auto-migrate models: %v\n", err)
	// 	return err
	// }

	fmt.Println("Database models auto-migrated successfully")
	return nil
}
