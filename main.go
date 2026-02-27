package main

import (
	"bitbucket/api_format/db"
	"bitbucket/api_format/routes"
	"bitbucket/api_format/utils"
	"fmt"
	"log"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"gorm.io/gorm"
)

var router = fiber.New()

// NOTE: UniPDF License Setup
// UniPDF requires a license to extract text from PDFs.
// To enable PDF text extraction:
//   1. Get a free license from https://unidoc.io
//   2. Contact UniPDF support for the license file/key
//   3. Place license file in the working directory or set via environment
//
// Without a license, PDF text extraction will fail with:
//   "unipdf license code required"
//
// The bankinvoice/upload endpoint will return a 422 error with diagnostic details
// if text extraction is not available.

// DatabaseConnections holds all database connections
type DatabaseConnections struct {
	Dev      *gorm.DB
	Live     *gorm.DB
	Prod     *gorm.DB
	AutoDev  *gorm.DB
	AutoProd *gorm.DB
}

var dbConnections DatabaseConnections

// initAllDatabases connects to all environments at startup
func initAllDatabases() error {
	fmt.Println("Initializing all database connections...\n")

	// Connect to Dev
	connectDB("Dev", func() {
		db.DevConnect()
		dbConnections.Dev = db.DB
	})

	// Connect to Live
	connectDB("Live", func() {
		db.LiveConnect()
		dbConnections.Live = db.DB
	})

	// Connect to Prod
	connectDB("Prod", func() {
		db.ProdConnect()
		dbConnections.Prod = db.DB
	})

	// Connect to AutoDev
	connectDB("AutoDev", func() {
		db.AutoConnectDEV()
		dbConnections.AutoDev = db.DB
	})

	// Connect to AutoProd
	connectDB("AutoProd", func() {
		db.AutoConnect()
		dbConnections.AutoProd = db.DB
	})

	// Check if at least one database is connected
	if dbConnections.Dev == nil && dbConnections.Live == nil &&
		dbConnections.Prod == nil && dbConnections.AutoDev == nil &&
		dbConnections.AutoProd == nil {
		return fmt.Errorf("no databases could be connected")
	}

	return nil
}

// connectDB safely connects to a database and handles errors
func connectDB(name string, connectFunc func()) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("✗ %s database connection failed: %v\n", name, r)
		}
	}()

	connectFunc()
	fmt.Printf("✓ %s database connected\n", name)
}

// getDBFromEnvironment returns the appropriate database connection
func getDBFromEnvironment(environment string) *gorm.DB {
	switch strings.ToLower(environment) {
	case "live":
		if dbConnections.Live != nil {
			return dbConnections.Live
		}
		fmt.Printf("Warning: Live database not available, falling back to Dev\n")
		return dbConnections.Dev
	case "prod":
		if dbConnections.Prod != nil {
			return dbConnections.Prod
		}
		fmt.Printf("Warning: Prod database not available, falling back to Dev\n")
		return dbConnections.Dev
	case "autodev":
		if dbConnections.AutoDev != nil {
			return dbConnections.AutoDev
		}
		fmt.Printf("Warning: AutoDev database not available, falling back to Dev\n")
		return dbConnections.Dev
	case "autoprod":
		if dbConnections.AutoProd != nil {
			return dbConnections.AutoProd
		}
		fmt.Printf("Warning: AutoProd database not available, falling back to Dev\n")
		return dbConnections.Dev
	default:
		if dbConnections.Dev != nil {
			return dbConnections.Dev
		}
		return nil
	}
}

// extractEnvironmentFromPath extracts environment from URL path
// Example: /autodev/v1/service/sales/getall -> autodev
func extractEnvironmentFromPath(path string) string {
	parts := strings.Split(strings.TrimPrefix(path, "/"), "/")
	if len(parts) > 0 {
		env := parts[0]
		// Validate environment
		validEnvs := map[string]bool{
			"dev":      true,
			"live":     true,
			"prod":     true,
			"autodev":  true,
			"autoprod": true,
		}
		if validEnvs[strings.ToLower(env)] {
			return env
		}
	}
	return "dev" // default to dev
}

// environmentMiddleware extracts environment from URL path and sets the appropriate DB
func environmentMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Extract environment from URL path
		// Example: /autodev/v1/service/sales/getall
		environment := extractEnvironmentFromPath(c.Path())

		// Get the appropriate database connection
		dbConnection := getDBFromEnvironment(environment)

		if dbConnection == nil {
			return c.Status(fiber.StatusServiceUnavailable).JSON(fiber.Map{
				"error": "database not available",
			})
		}

		// Update the global db.DB to use the environment-specific connection
		db.DB = dbConnection

		// Store in context for reference
		c.Locals("db", dbConnection)
		c.Locals("environment", environment)

		fmt.Printf("[%s] Request to %s\n", strings.ToUpper(environment), c.Path())
		return c.Next()
	}
}

func main() {
	// Initialize all database connections ONCE at startup
	if err := initAllDatabases(); err != nil {
		log.Fatalf("Failed to initialize databases: %v", err)
	}

	fmt.Println("\nAll available databases initialized successfully!")

	// Initialize database models (auto-migrate) on the preferred connection
	// Prefer AutoDev (qea_auto_dev) when available, otherwise fall back to Dev
	var migrateDB *gorm.DB
	if dbConnections.AutoDev != nil {
		migrateDB = dbConnections.AutoDev
	} else if dbConnections.AutoProd != nil {
		migrateDB = dbConnections.AutoProd
	} else {
		migrateDB = dbConnections.Live
	}

	if migrateDB != nil {
		if err := db.InitializeDatabase(migrateDB); err != nil {
			fmt.Printf("Warning: Database initialization incomplete: %v\n", err)
		}
	} else {
		fmt.Printf("Warning: No database connection available for auto-migrate\n")
	}

	// Load R2 configuration from config.yaml
	utils.LoadR2Config()

	// R2 and other config should be loaded via existing config utilities in `utils/config.go`

	// Create Fiber app with increased body size limit
	router = fiber.New(fiber.Config{
		BodyLimit: 2 * 1024 * 1024 * 1024, // 2 GB max upload size (increase as needed)
	})

	// CORS middleware
	router.Use(cors.New(cors.Config{
		AllowHeaders:     "Origin,Content-Type,Accept,Content-Length,Accept-Language,Accept-Encoding,Connection,Access-Control-Allow-Origin",
		AllowOrigins:     "*",
		AllowCredentials: false,
		AllowMethods:     "GET,POST,HEAD,PUT,DELETE,PATCH,OPTIONS",
	}))

	// Environment middleware - MUST be before routes to set db.DB
	router.Use(environmentMiddleware())

	// Setup routes
	routes.Setup(router)

	// go jobs.StartCronJobs()

	router.Listen(":1033")

	// port, cert, key, _ := utils.Urls()

	// err := router.ListenTLS(port, cert, key)
	// if err != nil {
	// 	log.Fatalf("Error starting server: %v", err)
	// }
}
