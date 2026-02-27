package models

import "time"

// VATTransaction represents a transaction for VAT purposes
type VATTransaction struct {
	ID              int64     `gorm:"primaryKey;column:id" json:"id"`
	TenantInfoID    int64     `gorm:"column:tenant_info_id;index" json:"tenant_info_id"`
	Box             int       `gorm:"column:box" json:"box"` // VAT return box number (1-9)
	GrossAmount     float64   `gorm:"column:gross_amount" json:"gross_amount"`
	VATScheme       string    `gorm:"column:vat_scheme" json:"vat_scheme"` // standard, margin, zero, exempt
	VATAmount       float64   `gorm:"column:vat_amount" json:"vat_amount"`
	CalculationMode string    `gorm:"column:calculation_mode" json:"calculation_mode"` // sales, purchase, overhead, cost
	Transaction     string    `gorm:"column:transaction" json:"transaction"`           // sales, purchase, overhead, cost (mirrors calculation_mode)
	NetAmount       float64   `gorm:"column:net_amount" json:"net_amount"`
	Name            string    `gorm:"column:name" json:"name"`               // regno, overhead type, etc
	VatReceipt      *string   `gorm:"column:vat_receipt" json:"vat_receipt"` // VAT receipt number or reference
	Details         string    `gorm:"column:details" json:"details"`         // stock num, detail, etc
	Description     string    `gorm:"column:description" json:"description"` // Additional description
	CustomerID      *int64    `gorm:"column:customer_id" json:"customer_id"` // Link to customer (optional)
	SupplierID      *int64    `gorm:"column:supplier_id" json:"supplier_id"` // Link to supplier (optional)
	Label           string    `gorm:"column:label;index" json:"label"`       // Quarter label like Q12025
	Date            time.Time `gorm:"column:date;index" json:"date"`
	CreatedAt       time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt       time.Time `gorm:"column:updated_at" json:"updated_at"`
}

func (VATTransaction) TableName() string {
	return "vat_transactions"
}

// VehicleSalesVAT represents vehicle sales for margin or VAT-qualifying scheme
// Note: `selling_price` is stored as GROSS (includes VAT). Use division by 1.20 to derive net when needed.
type VehicleSalesVAT struct {
	ID            int64     `gorm:"primaryKey;column:id" json:"id"`
	TenantInfoID  int64     `gorm:"column:tenant_info_id;index" json:"tenant_info_id"`
	VehicleID     int64     `gorm:"column:vehicle_id;index" json:"vehicle_id"`
	PurchasePrice float64   `gorm:"column:purchase_price" json:"purchase_price"`
	SellingPrice  float64   `gorm:"column:selling_price" json:"selling_price"`
	Margin        float64   `gorm:"column:margin" json:"margin"`
	MarginVAT     float64   `gorm:"column:margin_vat" json:"margin_vat"`
	VATScheme     string    `gorm:"column:vat_scheme" json:"vat_scheme"` // margin, vat_qualifying, zero
	SaleDate      time.Time `gorm:"column:sale_date;index" json:"sale_date"`
	CustomerID    *int64    `gorm:"column:customer_id" json:"customer_id"`
	InvoiceID     *int64    `gorm:"column:invoice_id" json:"invoice_id"`
	CreatedAt     time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt     time.Time `gorm:"column:updated_at" json:"updated_at"`
}

func (VehicleSalesVAT) TableName() string {
	return "vehicle_sales_vat"
}

// WorkshopTransaction represents workshop labour and parts transactions
type WorkshopTransaction struct {
	ID           int64     `gorm:"primaryKey;column:id" json:"id"`
	TenantInfoID int64     `gorm:"column:tenant_info_id;index" json:"tenant_info_id"`
	Date         time.Time `gorm:"column:date;index" json:"date"`
	LabourNet    float64   `gorm:"column:labour_net" json:"labour_net"`
	LabourVAT    float64   `gorm:"column:labour_vat" json:"labour_vat"`
	PartsNet     float64   `gorm:"column:parts_net" json:"parts_net"`
	PartsVAT     float64   `gorm:"column:parts_vat" json:"parts_vat"`
	ZeroRatedNet float64   `gorm:"column:zero_rated_net" json:"zero_rated_net"`
	CustomerID   *int64    `gorm:"column:customer_id" json:"customer_id"`
	InvoiceID    *int64    `gorm:"column:invoice_id" json:"invoice_id"`
	Description  string    `gorm:"column:description" json:"description"`
	CreatedAt    time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt    time.Time `gorm:"column:updated_at" json:"updated_at"`
}

func (WorkshopTransaction) TableName() string {
	return "workshop_transactions"
}

// VATPurchase represents a purchase transaction for VAT reclaim
type VATPurchase struct {
	ID            int64     `gorm:"primaryKey;column:id" json:"id"`
	TenantInfoID  int64     `gorm:"column:tenant_info_id;index" json:"tenant_info_id"`
	Date          time.Time `gorm:"column:date;index" json:"date"`
	SupplierID    *int64    `gorm:"column:supplier_id" json:"supplier_id"`
	SupplierName  string    `gorm:"column:supplier_name" json:"supplier_name"`
	NetAmount     float64   `gorm:"column:net_amount" json:"net_amount"`
	VATAmount     float64   `gorm:"column:vat_amount" json:"vat_amount"`
	GrossAmount   float64   `gorm:"column:gross_amount" json:"gross_amount"`
	Reclaimable   bool      `gorm:"column:reclaimable;index" json:"reclaimable"`
	Category      string    `gorm:"column:category" json:"category"` // car_purchase, parts, supplies, overheads, other
	Description   string    `gorm:"column:description" json:"description"`
	InvoiceNumber string    `gorm:"column:invoice_number" json:"invoice_number"`
	CreatedAt     time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt     time.Time `gorm:"column:updated_at" json:"updated_at"`
}

func (VATPurchase) TableName() string {
	return "vat_purchases"
}

// ECTransaction represents EU acquisitions and supplies
type ECTransaction struct {
	ID                     int64     `gorm:"primaryKey;column:id" json:"id"`
	TenantInfoID           int64     `gorm:"column:tenant_info_id;index" json:"tenant_info_id"`
	Date                   time.Time `gorm:"column:date;index" json:"date"`
	Type                   string    `gorm:"column:type;index" json:"type"` // acquisition, supply
	NetAmount              float64   `gorm:"column:net_amount" json:"net_amount"`
	VATDue                 float64   `gorm:"column:vat_due" json:"vat_due"`
	CustomerOrSupplierID   *int64    `gorm:"column:customer_or_supplier" json:"customer_or_supplier_id"`
	CustomerOrSupplierName string    `gorm:"column:customer_or_supplier_name" json:"customer_or_supplier_name"`
	CountryCode            string    `gorm:"column:country_code" json:"country_code"`
	Description            string    `gorm:"column:description" json:"description"`
	CreatedAt              time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt              time.Time `gorm:"column:updated_at" json:"updated_at"`
}

func (ECTransaction) TableName() string {
	return "ec_transactions"
}

// VATPeriod represents a VAT return period (MTD submission)
type VATPeriod struct {
	ID              int64      `gorm:"primaryKey;column:id" json:"id"`
	TenantInfoID    int64      `gorm:"column:tenant_info_id;index" json:"tenant_info_id"`
	PeriodKey       string     `gorm:"column:period_key;index" json:"period_key"` // e.g., "23A1"
	StartDate       time.Time  `gorm:"column:start_date;index" json:"start_date"`
	EndDate         time.Time  `gorm:"column:end_date" json:"end_date"`
	Box1            float64    `gorm:"column:box1" json:"box1"`
	Box2            float64    `gorm:"column:box2" json:"box2"`
	Box3            float64    `gorm:"column:box3" json:"box3"`
	Box4            float64    `gorm:"column:box4" json:"box4"`
	Box5            float64    `gorm:"column:box5" json:"box5"`
	Box6            float64    `gorm:"column:box6" json:"box6"`
	Box7            float64    `gorm:"column:box7" json:"box7"`
	Box8            float64    `gorm:"column:box8" json:"box8"`
	Box9            float64    `gorm:"column:box9" json:"box9"`
	Finalised       bool       `gorm:"column:finalised" json:"finalised"`
	SubmittedToHMRC bool       `gorm:"column:submitted_to_hmrc" json:"submitted_to_hmrc"`
	SubmissionDate  *time.Time `gorm:"column:submission_date" json:"submission_date"`
	CreatedAt       time.Time  `gorm:"column:created_at" json:"created_at"`
	UpdatedAt       time.Time  `gorm:"column:updated_at" json:"updated_at"`
}

func (VATPeriod) TableName() string {
	return "vat_periods"
}

// TenantVATConfig represents VAT configuration for a tenant (quarters and submission tracking)
type TenantVATConfig struct {
	TenantVATConfigID   int64     `gorm:"primaryKey;column:tenant_vat_config_id" json:"tenant_vat_config_id"`
	TenantInfoID        *int64    `gorm:"column:tenant_info_id;index" json:"tenant_info_id"`
	TenantName          *string   `gorm:"column:tenant_name" json:"tenant_name"`
	Quarter             *string   `gorm:"column:quater" json:"quarter"`
	VATQuarterMonths    *string   `gorm:"column:vat_quater_months" json:"vat_quarter_months"`
	VATQuarterStartDate *string   `gorm:"column:vat_quater_start_date" json:"vat_quarter_start_date"`
	VATQuarterEndDate   *string   `gorm:"column:vat_quater_end_date" json:"vat_quarter_end_date"`
	Year                *string   `gorm:"column:year" json:"year"`
	Label               *string   `gorm:"column:label" json:"label"`                                 // Format: Q12025 (Quarter + Year)
	VATSubmissionStatus *string   `gorm:"column:vat_submission_status" json:"vat_submission_status"` // Y/N
	CreatedAt           time.Time `gorm:"column:created;autoCreateTime:false" json:"created_at"`
	UpdatedAt           time.Time `gorm:"column:updated;autoUpdateTime:false" json:"updated_at"`
}

func (TenantVATConfig) TableName() string {
	return "tenant_vat_configs"
}
