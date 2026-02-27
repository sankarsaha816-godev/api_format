package controllers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path"
	"strconv"
	"strings"
	"sync"
	"time"

	"bitbucket/api_format/db"
	"bitbucket/api_format/domain"
	"bitbucket/api_format/internal/extractor"
	"bitbucket/api_format/models"

	"github.com/gofiber/fiber/v2"
)

var converterService = domain.NewConverterService("v1")
var bankInvoiceMigrationOnce sync.Once

func BankInvoiceHealth(c *fiber.Ctx) error {
	pdfToTextAvailable, pdfToTextBinary, pdfToTextErr := extractor.PdfToTextStatus()
	pdfToTextError := ""
	if pdfToTextErr != nil {
		pdfToTextError = pdfToTextErr.Error()
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"code":    fiber.StatusOK,
		"status":  true,
		"message": "ok",
		"version": converterService.Version(),
		"extractor": fiber.Map{
			"pdftotext_available": pdfToTextAvailable,
			"pdftotext_binary":    pdfToTextBinary,
			"pdftotext_error":     pdfToTextError,
		},
	})
}

func GetBankInvoiceHistory(c *fiber.Ctx) error {
	if err := ensureBankInvoiceTable(); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"code":    fiber.StatusInternalServerError,
			"status":  false,
			"message": "failed to initialize bank invoice table: " + err.Error(),
		})
	}

	limit := c.QueryInt("limit", 50)
	if limit <= 0 {
		limit = 50
	}
	if limit > 500 {
		limit = 500
	}

	offset := c.QueryInt("offset", 0)
	if offset < 0 {
		offset = 0
	}

	query := db.DB.Model(&models.BankInvoiceConversion{})

	if environment := strings.TrimSpace(c.Query("environment")); environment != "" {
		query = query.Where("environment = ?", environment)
	}
	if bank := strings.TrimSpace(c.Query("bank")); bank != "" {
		query = query.Where("bank = ?", bank)
	}
	if accountNumber := strings.TrimSpace(c.Query("accountNumber")); accountNumber != "" {
		query = query.Where("account_number = ?", accountNumber)
	}
	tenantInfoID := strings.TrimSpace(c.Query("tenant_info_id"))
	if tenantInfoID == "" {
		tenantInfoID = strings.TrimSpace(c.Query("tenantid"))
	}
	if tenantInfoID != "" {
		query = query.Where("tenant_info_id = ?", tenantInfoID)
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"code":    fiber.StatusInternalServerError,
			"status":  false,
			"message": "failed to count bank invoice history: " + err.Error(),
		})
	}

	var records []models.BankInvoiceConversion
	if err := query.Order("id desc").Limit(limit).Offset(offset).Find(&records).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"code":    fiber.StatusInternalServerError,
			"status":  false,
			"message": "failed to fetch bank invoice history: " + err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"code":    fiber.StatusOK,
		"status":  true,
		"message": "Bank invoice history fetched successfully",
		"details": fiber.Map{
			"total":   total,
			"limit":   limit,
			"offset":  offset,
			"records": records,
		},
	})
}

func GetBankInvoiceHistoryByTenantID(c *fiber.Ctx) error {
	if err := ensureBankInvoiceTable(); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"code":    fiber.StatusInternalServerError,
			"status":  false,
			"message": "failed to initialize bank invoice table: " + err.Error(),
		})
	}

	tenantID := strings.TrimSpace(c.Params("tenantid"))
	if tenantID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"code":    fiber.StatusBadRequest,
			"status":  false,
			"message": "tenantid is required",
		})
	}

	limit := c.QueryInt("limit", 50)
	if limit <= 0 {
		limit = 50
	}
	if limit > 500 {
		limit = 500
	}

	offset := c.QueryInt("offset", 0)
	if offset < 0 {
		offset = 0
	}

	query := db.DB.Model(&models.BankInvoiceConversion{}).Where("tenant_info_id = ?", tenantID)

	if bank := strings.TrimSpace(c.Query("bank")); bank != "" {
		query = query.Where("bank = ?", bank)
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"code":    fiber.StatusInternalServerError,
			"status":  false,
			"message": "failed to count bank invoice history: " + err.Error(),
		})
	}

	var records []models.BankInvoiceConversion
	if err := query.Order("id desc").Limit(limit).Offset(offset).Find(&records).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"code":    fiber.StatusInternalServerError,
			"status":  false,
			"message": "failed to fetch bank invoice history: " + err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"code":    fiber.StatusOK,
		"status":  true,
		"message": "Bank invoice history fetched successfully",
		"details": fiber.Map{
			"tenantid": tenantID,
			"total":    total,
			"limit":    limit,
			"offset":   offset,
			"records":  records,
		},
	})
}

func GetBankInvoiceHistoryByID(c *fiber.Ctx) error {
	if err := ensureBankInvoiceTable(); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"code":    fiber.StatusInternalServerError,
			"status":  false,
			"message": "failed to initialize bank invoice table: " + err.Error(),
		})
	}

	id := c.Params("id")
	var record models.BankInvoiceConversion
	if err := db.DB.First(&record, "id = ?", id).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"code":    fiber.StatusNotFound,
			"status":  false,
			"message": "bank invoice history record not found",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"code":    fiber.StatusOK,
		"status":  true,
		"message": "Bank invoice history record fetched successfully",
		"details": record,
	})
}

// ConvertBankInvoicePDFToCSV converts uploaded bank statement PDF into structured JSON + CSV.
func ConvertBankInvoicePDFToCSV(c *fiber.Ctx) error {
	pdfPath, originalName, cleanup, err := resolveBankInvoicePDFSource(c)
	if err != nil {
		status := mapConvertErrorToStatus(err)
		if status == fiber.StatusInternalServerError {
			status = fiber.StatusBadRequest
		}
		return c.Status(status).JSON(fiber.Map{
			"code":    status,
			"status":  false,
			"message": err.Error(),
		})
	}
	defer cleanup()

	bankHint := strings.TrimSpace(firstNonEmpty(c.FormValue("bank"), c.Query("bank")))
	includeHeader := parseHeaderFlag(firstNonEmpty(c.FormValue("header"), c.Query("header")))

	result, err := converterService.ConvertPDF(pdfPath, bankHint, includeHeader)
	if err != nil {
		status := mapConvertErrorToStatus(err)
		return c.Status(status).JSON(fiber.Map{
			"code":    status,
			"status":  false,
			"message": err.Error(),
		})
	}

	if err := persistBankInvoiceConversion(c, originalName, bankHint, includeHeader, result); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"code":    fiber.StatusInternalServerError,
			"status":  false,
			"message": "conversion succeeded but failed to save response: " + err.Error(),
		})
	}

	if strings.EqualFold(c.Query("format"), "csv") {
		c.Set("Content-Type", "text/csv; charset=utf-8")
		c.Set("Content-Disposition", "attachment; filename=bank_statement.csv")
		return c.SendString(result.CSV)
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"bank":    result.Bank,
		"accountInfo": fiber.Map{
			"holder":   result.AccountHolder,
			"number":   result.AccountNumber,
			"sortCode": result.SortCode,
		},
		"transactions": result.Transactions,
		"csv":          result.CSV,
		"totalDebit":   result.TotalDebits,
		"totalCredit":  result.TotalCredits,
		"count":        result.TransactionCount,
		"rawText":      result.RawText,
		"version":      result.Version,
		"debugLines":   result.DebugLines,
	})
}

func resolveBankInvoicePDFSource(c *fiber.Ctx) (string, string, func(), error) {
	if file, err := c.FormFile("file"); err == nil && file != nil {
		tempFile, createErr := os.CreateTemp("", "bank-statement-*.pdf")
		if createErr != nil {
			return "", "", func() {}, fmt.Errorf("failed to prepare temporary file")
		}
		tempPath := tempFile.Name()
		_ = tempFile.Close()

		if saveErr := c.SaveFile(file, tempPath); saveErr != nil {
			_ = os.Remove(tempPath)
			return "", "", func() {}, fmt.Errorf("failed to save uploaded PDF")
		}

		return tempPath, file.Filename, func() { _ = os.Remove(tempPath) }, nil
	}

	fileURL := strings.TrimSpace(firstNonEmpty(c.FormValue("file_url"), c.Query("file_url"), extractFileURLFromJSONBody(c)))
	if fileURL == "" {
		return "", "", func() {}, fmt.Errorf("PDF input required: send multipart form-data field 'file' or provide 'file_url'")
	}

	tempPath, fileName, downloadErr := downloadPDFToTemp(fileURL)
	if downloadErr != nil {
		return "", "", func() {}, downloadErr
	}

	return tempPath, fileName, func() { _ = os.Remove(tempPath) }, nil
}

func extractFileURLFromJSONBody(c *fiber.Ctx) string {
	contentType := strings.ToLower(c.Get("Content-Type"))
	if !strings.Contains(contentType, "application/json") {
		return ""
	}

	var payload map[string]interface{}
	if err := json.Unmarshal(c.Body(), &payload); err != nil {
		return ""
	}

	value, _ := payload["file_url"].(string)
	return strings.TrimSpace(value)
}

func downloadPDFToTemp(fileURL string) (string, string, error) {
	parsedURL, err := url.Parse(fileURL)
	if err != nil || parsedURL.Scheme == "" || parsedURL.Host == "" {
		return "", "", fmt.Errorf("invalid file_url")
	}

	client := &http.Client{Timeout: 60 * time.Second}
	resp, err := client.Get(fileURL)
	if err != nil {
		return "", "", fmt.Errorf("failed to download file_url: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		return "", "", fmt.Errorf("file_url returned non-success status: %d", resp.StatusCode)
	}

	contentType := strings.ToLower(strings.TrimSpace(resp.Header.Get("Content-Type")))

	tempFile, err := os.CreateTemp("", "bank-statement-url-*.pdf")
	if err != nil {
		return "", "", fmt.Errorf("failed to prepare temporary file")
	}
	tempPath := tempFile.Name()

	const maxPDFBytes = 50 * 1024 * 1024
	limitedBody, readErr := io.ReadAll(io.LimitReader(resp.Body, maxPDFBytes+1))
	if readErr != nil {
		_ = os.Remove(tempPath)
		return "", "", fmt.Errorf("failed to download PDF content")
	}

	written, writeErr := tempFile.Write(limitedBody)
	closeErr := tempFile.Close()
	if writeErr != nil {
		_ = os.Remove(tempPath)
		return "", "", fmt.Errorf("failed to save downloaded PDF content")
	}
	if closeErr != nil {
		_ = os.Remove(tempPath)
		return "", "", fmt.Errorf("failed to finalize temporary file")
	}
	if written > maxPDFBytes {
		_ = os.Remove(tempPath)
		return "", "", fmt.Errorf("file_url exceeds 50MB limit")
	}

	preview := limitedBody
	if len(preview) > 2048 {
		preview = preview[:2048]
	}

	if !looksLikePDF(preview) {
		_ = os.Remove(tempPath)
		trimmedPreview := abbreviateText(strings.TrimSpace(string(preview)), 180)
		if strings.Contains(contentType, "text/html") || strings.Contains(contentType, "application/json") {
			return "", "", fmt.Errorf("file_url did not return a valid PDF (content-type: %s, preview: %q)", contentType, trimmedPreview)
		}
		return "", "", fmt.Errorf("file_url does not appear to be a valid PDF (content-type: %s)", contentType)
	}

	name := path.Base(parsedURL.Path)
	if strings.TrimSpace(name) == "" || name == "/" || name == "." {
		name = "statement.pdf"
	}

	return tempPath, name, nil
}

func looksLikePDF(preview []byte) bool {
	if len(preview) == 0 {
		return false
	}
	trimmed := bytes.TrimSpace(preview)
	return bytes.HasPrefix(trimmed, []byte("%PDF-"))
}

func abbreviateText(value string, maxLen int) string {
	if maxLen <= 0 || len(value) <= maxLen {
		return value
	}
	return value[:maxLen] + "..."
}

func firstNonEmpty(values ...string) string {
	for _, value := range values {
		trimmed := strings.TrimSpace(value)
		if trimmed != "" {
			return trimmed
		}
	}
	return ""
}

func parseHeaderFlag(raw string) bool {
	if strings.TrimSpace(raw) == "" {
		return true
	}

	switch strings.ToLower(strings.TrimSpace(raw)) {
	case "false", "0", "no", "off":
		return false
	default:
		return true
	}
}

func mapConvertErrorToStatus(err error) int {
	message := strings.ToLower(err.Error())
	if strings.Contains(message, "not return a valid pdf") ||
		strings.Contains(message, "does not appear to be a valid pdf") ||
		strings.Contains(message, "non-success status") ||
		strings.Contains(message, "invalid file_url") ||
		strings.Contains(message, "unsupported or malformed content streams") {
		return fiber.StatusBadRequest
	}
	if strings.Contains(message, "extract pages") || strings.Contains(message, "extract text failed") {
		return fiber.StatusUnprocessableEntity
	}
	if strings.Contains(message, "no transactions") {
		return fiber.StatusUnprocessableEntity
	}
	if strings.Contains(message, "parse statement") {
		return fiber.StatusUnprocessableEntity
	}
	if strings.Contains(message, "unsupported bank") || strings.Contains(message, "auto-detect") || strings.Contains(message, "required") {
		return fiber.StatusBadRequest
	}

	fmt.Println("bankinvoice conversion error:", err)
	return fiber.StatusInternalServerError
}

func persistBankInvoiceConversion(c *fiber.Ctx, fileName, requestBank string, includeHeader bool, result *models.BankStatementConvertResult) error {
	if err := ensureBankInvoiceTable(); err != nil {
		return err
	}

	tenantInfoID, err := parseTenantInfoID(c)
	if err != nil {
		return err
	}

	transactionsJSONBytes, err := json.Marshal(result.Transactions)
	if err != nil {
		return fmt.Errorf("marshal transactions: %w", err)
	}

	responseJSONBytes, err := json.Marshal(result)
	if err != nil {
		return fmt.Errorf("marshal response: %w", err)
	}

	environment, _ := c.Locals("environment").(string)

	record := models.BankInvoiceConversion{
		Environment:      environment,
		TenantInfoID:     tenantInfoID,
		FileName:         fileName,
		RequestBank:      requestBank,
		IncludeHeader:    includeHeader,
		Bank:             result.Bank,
		AccountHolder:    result.AccountHolder,
		AccountNumber:    result.AccountNumber,
		SortCode:         result.SortCode,
		StatementPeriod:  result.StatementPeriod,
		TransactionCount: result.TransactionCount,
		TransactionsJSON: string(transactionsJSONBytes),
		CSVOutput:        result.CSV,
		ResponseJSON:     string(responseJSONBytes),
	}

	if err := db.DB.Create(&record).Error; err != nil {
		return fmt.Errorf("insert bank invoice conversion: %w", err)
	}

	return nil
}

func parseTenantInfoID(c *fiber.Ctx) (*int, error) {
	raw := strings.TrimSpace(c.FormValue("tenant_info_id"))
	if raw == "" {
		raw = strings.TrimSpace(c.Query("tenant_info_id"))
	}
	if raw == "" {
		return nil, nil
	}

	parsed, err := strconv.Atoi(raw)
	if err != nil {
		return nil, fmt.Errorf("invalid tenant_info_id value %q", raw)
	}

	return &parsed, nil
}

func ensureBankInvoiceTable() error {
	if db.DB == nil {
		return fmt.Errorf("database connection not available")
	}

	var migrateErr error
	bankInvoiceMigrationOnce.Do(func() {
		migrateErr = db.DB.AutoMigrate(&models.BankInvoiceConversion{})
	})
	if migrateErr != nil {
		return fmt.Errorf("auto-migrate bank_invoice_conversions: %w", migrateErr)
	}

	return nil
}
