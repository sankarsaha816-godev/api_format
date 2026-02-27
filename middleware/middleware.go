package middleware

import (
	"fmt"
	"strings"
	"time"

	"bitbucket/api_format/db"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

type AuthLog struct {
	ID          uint      `gorm:"primaryKey;autoIncrement"`
	TenantID    int       `gorm:"column:tenant_id"`
	UserType    string    `gorm:"column:user_type"`
	TenantName  string    `gorm:"column:tenant_name"`
	IPAddress   string    `gorm:"column:ip_address"`
	RequestPath string    `gorm:"column:request_path"`
	Status      string    `gorm:"column:status"` // success / error
	Message     string    `gorm:"column:message"`
	CreatedAt   time.Time `gorm:"autoCreateTime"`
}

func JWTAuth(c *fiber.Ctx) error {
	authHeader := c.Get("Authorization")
	if authHeader == "" {
		logAuthAttempt(0, "", "", c.IP(), c.OriginalURL(), "error", "Missing Authorization header")
		return unauthorizedResponse(c, "Missing Authorization header")
	}

	tokenStr := strings.Replace(authHeader, "Bearer ", "", 1)
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		return []byte("secret_key"), nil
	})

	if err != nil || !token.Valid {
		logAuthAttempt(0, "", "", c.IP(), c.OriginalURL(), "error", "Invalid or expired token")
		return unauthorizedResponse(c, "Invalid or expired token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		logAuthAttempt(0, "", "", c.IP(), c.OriginalURL(), "error", "Unable to parse token claims")
		return unauthorizedResponse(c, "Unable to parse token claims")
	}

	// Safe claim extraction with validation
	var tenantID int
	var userType, tenantName string
	var message string

	// Validate tenant_id
	valTenantID, ok := claims["tenant_id"].(float64)
	if !ok {
		message = "tenant_id missing or invalid in JWT"
		logAuthAttempt(0, "", "", c.IP(), c.OriginalURL(), "error", message)
		return unauthorizedResponse(c, message)
	}
	tenantID = int(valTenantID)

	// Validate user_type
	valUserType, ok := claims["user_type"].(string)
	if !ok || valUserType == "" {
		message = "user_type missing or invalid in JWT"
		logAuthAttempt(tenantID, "", "", c.IP(), c.OriginalURL(), "error", message)
		return unauthorizedResponse(c, message)
	}
	userType = valUserType

	// Validate tenant_name
	valTenantName, ok := claims["tenant_name"].(string)
	if !ok || valTenantName == "" {
		message = "tenant_name missing or invalid in JWT"
		logAuthAttempt(tenantID, userType, "", c.IP(), c.OriginalURL(), "error", message)
		return unauthorizedResponse(c, message)
	}
	tenantName = valTenantName

	// All claims valid — store and continue
	c.Locals("tenant_id", tenantID)
	c.Locals("user_type", userType)
	c.Locals("tenant_name", tenantName)

	logAuthAttempt(tenantID, userType, tenantName, c.IP(), c.OriginalURL(), "success", "JWT validated successfully")

	return c.Next()
}

// unauthorizedResponse returns 401 and stops further execution
func unauthorizedResponse(c *fiber.Ctx, message string) error {
	return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
		"status":  false,
		"message": message,
	})
}

// Logs all auth attempts
func logAuthAttempt(tenantID int, userType, tenantName, ip, path, status, message string) {
	log := AuthLog{
		TenantID:    tenantID,
		UserType:    userType,
		TenantName:  tenantName,
		IPAddress:   ip,
		RequestPath: path,
		Status:      status,
		Message:     message,
		CreatedAt:   time.Now(),
	}

	go func(l AuthLog) {
		defer func() {
			if r := recover(); r != nil {
				fmt.Printf(" Failed to log auth attempt (panic): %v\n", r)
			}
		}()
		if err := db.DB.Create(&l).Error; err != nil {
			fmt.Printf(" Failed to log auth attempt: %v\n", err)
		}
	}(log)
}


// package middleware

// import (
// 	"fmt"
// 	"strings"
// 	"time"

// 	"bitbucket/api_format/db"

// 	"github.com/gofiber/fiber/v2"
// 	"github.com/golang-jwt/jwt/v5"
// )

// type AuthLog struct {
// 	ID          uint      `gorm:"primaryKey;autoIncrement"`
// 	TenantID    int       `gorm:"column:tenant_id"`
// 	UserType    string    `gorm:"column:user_type"`
// 	TenantName  string    `gorm:"column:tenant_name"`
// 	IPAddress   string    `gorm:"column:ip_address"`
// 	RequestPath string    `gorm:"column:request_path"`
// 	Status      string    `gorm:"column:status"` // success / failed
// 	Message     string    `gorm:"column:message"`
// 	CreatedAt   time.Time `gorm:"autoCreateTime"`
// }

// func JWTAuth(c *fiber.Ctx) error {
// 	authHeader := c.Get("Authorization")
// 	if authHeader == "" {
// 		logAuthAttempt(0, "", "", c.IP(), c.OriginalURL(), "error", "Missing Authorization header")
// 		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
// 			"status":  false,
// 			"message": "Missing Authorization header",
// 		})
// 	}

// 	tokenStr := strings.Replace(authHeader, "Bearer ", "", 1)
// 	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
// 		return []byte("secret_key"), nil
// 	})

// 	if err != nil || !token.Valid {
// 		logAuthAttempt(0, "", "", c.IP(), c.OriginalURL(), "error", "Invalid or expired token")
// 		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
// 			"status":  false,
// 			"message": "Invalid or expired token",
// 		})
// 	}

// 	claims, ok := token.Claims.(jwt.MapClaims)
// 	if !ok {
// 		logAuthAttempt(0, "", "", c.IP(), c.OriginalURL(), "error", "Unable to parse token claims")
// 		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
// 			"status":  false,
// 			"message": "Unable to parse token claims",
// 		})
// 	}

// 	// Safe claim extraction
// 	var tenantID int
// 	var userType, tenantName string
// 	status := "success"
// 	message := "JWT validated"

// 	if val, ok := claims["tenant_id"].(float64); ok {
// 		tenantID = int(val)
// 	} else {
// 		status = "warning"
// 		message = "tenant_id missing in JWT"
// 	}

// 	if val, ok := claims["user_type"].(string); ok {
// 		userType = val
// 	} else {
// 		status = "warning"
// 		message = "user_type missing or invalid in JWT"
// 	}

// 	if val, ok := claims["tenant_name"].(string); ok {
// 		tenantName = val
// 	} else {
// 		status = "warning"
// 		message = "tenant_name missing or invalid in JWT"
// 	}

// 	// Store context
// 	c.Locals("tenant_id", tenantID)
// 	c.Locals("user_type", userType)
// 	c.Locals("tenant_name", tenantName)

// 	// Log authentication attempt
// 	logAuthAttempt(tenantID, userType, tenantName, c.IP(), c.OriginalURL(), status, message)

// 	return c.Next()
// }

// func logAuthAttempt(tenantID int, userType, tenantName, ip, path, status, message string) {
// 	log := AuthLog{
// 		TenantID:    tenantID,
// 		UserType:    userType,
// 		TenantName:  tenantName,
// 		IPAddress:   ip,
// 		RequestPath: path,
// 		Status:      status,
// 		Message:     message,
// 		CreatedAt:   time.Now(),
// 	}

// 	if err := db.DB.Create(&log).Error; err != nil {
// 		fmt.Printf(" Failed to log auth attempt: %v\n", err)
// 	}
// }
