package middleware

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"bitbucket/api_format/db"
	"bitbucket/api_format/models"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"gorm.io/gorm"
)

// CreditUsageLogger logs valuation API hits to credit_usage_ledger
func CreditUsageLogger() fiber.Handler {
	return func(c *fiber.Ctx) error {
		endpoint := normalizeEndpointPath(c.Path())
		if !strings.HasPrefix(endpoint, "/valuation/") {
			return c.Next()
		}

		tenantID := extractTenantIDFromQuery(c)
		if tenantID == 0 {
			return c.Next()
		}

		userID := extractUserIDFromJWT(c)
		var userIDPtr *int64
		if userID != 0 {
			userIDPtr = &userID
		}

		creditsCost := getEndpointCreditCost(endpoint)

		var subscription models.TenantSubscription
		if creditsCost > 0 {
			if err := db.DB.Where("tenant_info_id = ? AND status = ?", tenantID, 1).
					Order("period_end DESC").First(&subscription).Error; err != nil {
				return c.Status(fiber.StatusPaymentRequired).JSON(fiber.Map{
					"code":    fiber.StatusPaymentRequired,
					"message": "No active subscription",
				})
			}

			if subscription.CreditsRemaining < creditsCost {
				return c.Status(fiber.StatusPaymentRequired).JSON(fiber.Map{
					"code":               fiber.StatusPaymentRequired,
					"message":            "Credit limit reached",
					"credits_remaining":  subscription.CreditsRemaining,
					"period_end":         subscription.PeriodEnd,
					"tenant_info_id":     subscription.TenantInfoID,
					"required_credits":   creditsCost,
				})
			}
		}

		err := c.Next()

		requestID := c.Get("X-Request-Id")
		if strings.TrimSpace(requestID) == "" {
			requestID = fmt.Sprintf("hit-%d", time.Now().UnixNano())
		}

		entryStatus := "success"
		if c.Response().StatusCode() >= 400 {
			entryStatus = "failed"
		}

		if creditsCost > 0 && entryStatus == "success" {
			_ = db.DB.Transaction(func(tx *gorm.DB) error {
				if subscription.ID == 0 {
					if err := tx.Where("tenant_info_id = ? AND status = ?", tenantID, 1).
						Order("period_end DESC").First(&subscription).Error; err != nil {
						return err
					}
				}

				subscription.CreditsUsed += creditsCost
				if subscription.CreditsRemaining >= creditsCost {
					subscription.CreditsRemaining -= creditsCost
				} else {
					subscription.CreditsRemaining = 0
					entryStatus = "exceeded"
				}

				if err := tx.Save(&subscription).Error; err != nil {
					return err
				}
				return nil
			})
		}

		entry := models.CreditUsageLedger{
			TenantInfoID:   tenantID,
			UserID:         userIDPtr,
			Endpoint:       endpoint,
			CreditsCharged: creditsCost,
			RequestID:      requestID,
			Status:         entryStatus,
			CreatedAt:      time.Now().UTC(),
		}

		_ = db.DB.Create(&entry).Error

		return err
	}
}

func normalizeEndpointPath(path string) string {
	trimmed := strings.TrimPrefix(path, "/")
	parts := strings.Split(trimmed, "/")
	if len(parts) >= 3 && parts[1] == "v1" {
		normalized := "/" + strings.Join(parts[2:], "/")
		if strings.HasPrefix(normalized, "/valuation/get/") {
			return "/valuation/get/:id"
		}
		return normalized
	}
	return path
}

func getEndpointCreditCost(endpoint string) int {
	var apiCost models.APICreditLedger
	if err := db.DB.Where("api_url = ? AND status = ?", endpoint, 1).First(&apiCost).Error; err == nil {
		if apiCost.CreditsRequired != nil {
			return *apiCost.CreditsRequired
		}
	}

	return 1
}

func extractTenantIDFromQuery(c *fiber.Ctx) int64 {
	keys := []string{"tenant_info_id", "tenantid", "tenant_id", "tid"}
	for _, key := range keys {
		val := strings.TrimSpace(c.Query(key))
		if val == "" {
			continue
		}
		if id, err := strconv.ParseInt(val, 10, 64); err == nil {
			return id
		}
	}
	return 0
}

func extractUserIDFromJWT(c *fiber.Ctx) int64 {
	authHeader := c.Get("Authorization")
	if authHeader == "" {
		return 0
	}

	tokenStr := strings.Replace(authHeader, "Bearer ", "", 1)
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		return []byte("secret_key"), nil
	})
	if err != nil || !token.Valid {
		return 0
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return 0
	}

	if val, ok := claims["sub"].(string); ok {
		if id, err := strconv.ParseInt(val, 10, 64); err == nil {
			return id
		}
	}
	if val, ok := claims["user_id"].(float64); ok {
		return int64(val)
	}
	if val, ok := claims["userid"].(float64); ok {
		return int64(val)
	}
	return 0
}
