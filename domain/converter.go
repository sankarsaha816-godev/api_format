package domain

import (
	"fmt"
	"strings"

	"bitbucket/api_format/internal/extractor"
	internalmodels "bitbucket/api_format/internal/models"
	"bitbucket/api_format/internal/writer"
	"bitbucket/api_format/models"
	"bitbucket/api_format/parser"
)

// ConverterService orchestrates PDF extraction, parsing, and CSV generation.
type ConverterService struct {
	version string
}

func NewConverterService(version string) *ConverterService {
	return &ConverterService{version: version}
}

func (s *ConverterService) Version() string {
	return s.version
}

// ConvertPDF converts a bank statement PDF into structured transactions and CSV.
func (s *ConverterService) ConvertPDF(filePath, bankHint string, includeMetadata bool) (*models.BankStatementConvertResult, error) {
	pages, err := extractor.ExtractPages(filePath)
	if err != nil {
		return nil, fmt.Errorf("extract pages: %w", err)
	}

	bankType, err := resolveBankType(pages, bankHint)
	if err != nil {
		return nil, err
	}

	statementParser, err := parser.New(bankType)
	if err != nil {
		return nil, err
	}

	parsed, err := statementParser.Parse(pages)
	if err != nil {
		return nil, fmt.Errorf("parse statement: %w", err)
	}
	if len(parsed.Transactions) == 0 {
		return nil, fmt.Errorf("no transactions found in statement")
	}

	csvOutput, err := writer.BuildCSV(parsed, includeMetadata)
	if err != nil {
		return nil, fmt.Errorf("build csv: %w", err)
	}

	totalDebits, totalCredits := summarizeTotals(parsed.Transactions)

	return &models.BankStatementConvertResult{
		Bank:            string(parsed.Bank),
		AccountHolder:   parsed.AccountHolder,
		AccountNumber:   parsed.AccountNumber,
		SortCode:        parsed.SortCode,
		StatementPeriod: parsed.StatementPeriod,
		// Transactions:     parsed.Transactions,
		TransactionCount: len(parsed.Transactions),
		TotalDebits:      totalDebits,
		TotalCredits:     totalCredits,
		Net:              totalCredits - totalDebits,
		RawText:          strings.Join(pages, "\n--- PAGE BREAK ---\n"),
		Version:          s.version,
		DebugLines:       parsed.DebugLines,
		CSV:              csvOutput,
	}, nil
}

func summarizeTotals(transactions []internalmodels.Transaction) (float64, float64) {
	var totalDebits float64
	var totalCredits float64

	for _, transaction := range transactions {
		if transaction.Amount <= 0 {
			continue
		}

		switch strings.ToUpper(strings.TrimSpace(transaction.Type)) {
		case "DEBIT":
			totalDebits += transaction.Amount
		case "CREDIT":
			totalCredits += transaction.Amount
		}
	}

	return totalDebits, totalCredits
}

func resolveBankType(pages []string, bankHint string) (internalmodels.BankType, error) {
	normalized := strings.TrimSpace(strings.ToLower(bankHint))
	if normalized == "" {
		return parser.AutoDetect(pages)
	}

	bankType := internalmodels.BankType(normalized)
	if _, err := parser.New(bankType); err != nil {
		return "", fmt.Errorf("unsupported bank value %q, expected one of: metro, hsbc, barclays", bankHint)
	}
	return bankType, nil
}
