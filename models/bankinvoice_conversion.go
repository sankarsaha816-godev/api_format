package models

import "time"

// BankInvoiceConversion stores successful bank statement conversion responses.
type BankInvoiceConversion struct {
	ID               uint      `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	Environment      string    `gorm:"column:environment;type:varchar(30);index" json:"environment"`
	TenantInfoID     *int      `gorm:"column:tenant_info_id;index" json:"tenant_info_id,omitempty"`
	FileName         string    `gorm:"column:file_name;type:varchar(255)" json:"file_name"`
	RequestBank      string    `gorm:"column:request_bank;type:varchar(50)" json:"request_bank"`
	IncludeHeader    bool      `gorm:"column:include_header" json:"include_header"`
	Bank             string    `gorm:"column:bank;type:varchar(50);index" json:"bank"`
	AccountHolder    string    `gorm:"column:account_holder;type:varchar(255)" json:"account_holder"`
	AccountNumber    string    `gorm:"column:account_number;type:varchar(50)" json:"account_number"`
	SortCode         string    `gorm:"column:sort_code;type:varchar(50)" json:"sort_code"`
	StatementPeriod  string    `gorm:"column:statement_period;type:varchar(255)" json:"statement_period"`
	TransactionCount int       `gorm:"column:transaction_count" json:"transaction_count"`
	TransactionsJSON string    `gorm:"column:transactions_json;type:longtext" json:"transactions_json"`
	CSVOutput        string    `gorm:"column:csv_output;type:longtext" json:"csv_output"`
	ResponseJSON     string    `gorm:"column:response_json;type:longtext" json:"response_json"`
	CreatedAt        time.Time `gorm:"column:created_at;autoCreateTime" json:"created_at"`
	UpdatedAt        time.Time `gorm:"column:updated_at;autoUpdateTime" json:"updated_at"`
}

func (BankInvoiceConversion) TableName() string {
	return "bank_invoice_conversions"
}
