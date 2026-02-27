package models

import internalmodels "bitbucket/api_format/internal/models"

// BankStatementConvertResult is the response payload for bank statement conversion.
type BankStatementConvertResult struct {
	Bank             string                       `json:"bank"`
	AccountHolder    string                       `json:"accountHolder,omitempty"`
	AccountNumber    string                       `json:"accountNumber,omitempty"`
	SortCode         string                       `json:"sortCode,omitempty"`
	StatementPeriod  string                       `json:"statementPeriod,omitempty"`
	Transactions     []internalmodels.Transaction `json:"transactions"`
	TransactionCount int                          `json:"transactionCount"`
	TotalDebits      float64                      `json:"totalDebits"`
	TotalCredits     float64                      `json:"totalCredits"`
	Net              float64                      `json:"net"`
	RawText          string                       `json:"rawText,omitempty"`
	Version          string                       `json:"version,omitempty"`
	DebugLines       []internalmodels.DebugLine   `json:"debugLines,omitempty"`
	CSV              string                       `json:"csv"`
}
