package models

import "time"

type VehicleSale struct {
	SalesId             int     `gorm:"column:sales_id;primaryKey;autoIncrement" json:"sales_id"`
	TenantInfoID        *int    `gorm:"column:tenant_info_id" json:"tenant_info_id,omitempty"`
	VATType             *string `gorm:"column:vat_type;size:50" json:"vat_type"`
	VATSubmissionStatus *string `gorm:"column:vat_submission_status;size:50" json:"vat_submission_status"`
	StockNumber         *string `gorm:"column:stock_number;size:50" json:"stock_number"`
	Registration        *string `gorm:"column:registration;size:50" json:"registration"`
	CustomerID          *int    `gorm:"column:customer_id" json:"customer_id"`
	CustomerName        *string `gorm:"column:customer_name;size:100" json:"customer_name"`
	CustomerEmail       *string `gorm:"column:customer_email;size:100" json:"customer_email"`
	// New fields for margin VAT
	VehicleID     *int    `gorm:"column:vehicle_id" json:"vehicle_id,omitempty"`
	PurchasePrice *string `gorm:"column:purchase_price;size:50" json:"purchase_price,omitempty"`
	SellingPrice  *string `gorm:"column:selling_price;size:50" json:"selling_price,omitempty"`
	// margin and margin_vat are stored/generated in the DB; treat them as read-only
	Margin    *string `gorm:"column:margin;size:50;->" json:"margin,omitempty"`
	MarginVAT *string `gorm:"column:margin_vat;size:50;->" json:"margin_vat,omitempty"`
	// VatScheme     *string `gorm:"column:vat_scheme;size:50" json:"vat_scheme,omitempty"`
	SaleDate                 *string    `gorm:"column:sale_date;size:50" json:"sale_date"`
	SalesInvNo               *string    `gorm:"column:sales_inv_no;size:50" json:"sales_inv_no"`
	SaleType                 *string    `gorm:"column:sale_type;size:50" json:"sale_type"`
	SoldPrice                *string    `gorm:"column:sold_price;size:50" json:"sold_price"`
	FirstName                *string    `gorm:"column:first_name;size:100" json:"first_name"`
	LastName                 *string    `gorm:"column:last_name;size:100" json:"last_name"`
	AddressLine1             *string    `gorm:"column:address_line1;size:100" json:"address_line1"`
	AddressLine2             *string    `gorm:"column:address_line2;size:100" json:"address_line2"`
	City                     *string    `gorm:"column:city;size:100" json:"city"`
	State                    *string    `gorm:"column:state;size:100" json:"state"`
	Country                  *string    `gorm:"column:country;size:100" json:"country"`
	PostCode                 *string    `gorm:"column:post_code;size:100" json:"post_code"`
	Remarks                  *string    `gorm:"column:remarks;size:100" json:"remarks"`
	CalculationMode          *string    `gorm:"column:calculation_mode;size:50" json:"calculation_mode"`
	VATReceipt               *string    `gorm:"column:vat_receipt;size:50" json:"vat_receipt,omitempty"`
	VATReclaimable           *string    `gorm:"column:vat_reclaimable;size:50" json:"vat_reclaimable,omitempty"`
	VAT                      *string    `gorm:"column:vat;size:50" json:"vat,omitempty"`
	ServiceInvoicePending    *string    `gorm:"column:service_invoice_pending;size:50" json:"service_invoice_pending,omitempty"`
	PartsExchangePrice       *string    `gorm:"column:parts_exchange_price;size:50" json:"parts_exchange_price,omitempty"`
	PartsExchangeRegno       *string    `gorm:"column:parts_exchange_regno;size:50" json:"parts_exchange_regno,omitempty"`
	PartsExchangeStockNumber *string    `gorm:"column:parts_exchange_stock_number;size:50" json:"parts_exchange_stock_number,omitempty"`
	PartsExchangeVat         *string    `gorm:"column:parts_exchange_vat;size:50" json:"parts_exchange_vat,omitempty"`
	Finance                  *string    `gorm:"column:finance;size:50" json:"finance,omitempty"`
	CustomerPayment          *string    `gorm:"column:customer_payment;size:50" json:"customer_payment,omitempty"`
	ReserveDepositAmount     *string    `gorm:"column:reserve_deposit_amount;size:50" json:"reserve_deposit_amount,omitempty"`
	Status                   *int       `gorm:"column:status" json:"status,omitempty"` // "Sold"=1, "Stock"=2, "Reserve"=3, "On Hold"=4, "Others"=5
	CreatedBy                *string    `gorm:"column:created_by;size:50" json:"created_by,omitempty"`
	UpdatedBy                *string    `gorm:"column:updated_by;size:50" json:"updated_by,omitempty"`
	CreatedAt                *time.Time `gorm:"column:created_at;default:CURRENT_TIMESTAMP" json:"created_at,omitempty"`
	UpdatedAt                *time.Time `gorm:"column:updated_at;default:CURRENT_TIMESTAMP on update CURRENT_TIMESTAMP" json:"updated_at,omitempty"`
}

type Stocks struct {
	StockID      int  `gorm:"column:stock_id;primaryKey;autoIncrement" json:"stock_id"`
	TenantInfoID *int `gorm:"column:tenant_info_id;default:0" json:"tenant_info_id,omitempty"`

	StockNumber   *string `gorm:"column:stock_number;size:50" json:"stock_number,omitempty"`
	Registration  *string `gorm:"column:registration;size:50" json:"registration,omitempty"`
	DatePurchased *string `gorm:"column:date_purchased;size:50" json:"date_purchased,omitempty"`
	StockType     *string `gorm:"column:stock_type;size:50" json:"stock_type,omitempty"`

	SupplierName     *string `gorm:"column:supplier_name;size:50" json:"supplier_name,omitempty"`
	SupplierLocation *string `gorm:"column:supplier_location;size:50" json:"supplier_location,omitempty"`

	Make        *string `gorm:"column:make;size:50" json:"make,omitempty"`
	Model       *string `gorm:"column:model;size:50" json:"model,omitempty"`
	Description *string `gorm:"column:description;size:50" json:"description,omitempty"`
	Colour      *string `gorm:"column:colour;size:50" json:"colour,omitempty"`
	Mileage     *string `gorm:"column:mileage;size:50" json:"mileage,omitempty"`
	FuelType    *string `gorm:"column:fuel_type;size:50" json:"fuel_type,omitempty"`

	YearOfManufacturing *string `gorm:"column:year_of_manufacturing;size:50" json:"year_of_manufacturing,omitempty"`
	EngineCapacity      *string `gorm:"column:engine_capacity;size:50" json:"engine_capacity,omitempty"`

	PurchaseType   *string `gorm:"column:purchase_type;size:50" json:"purchase_type,omitempty"`
	PurchasePrice  *string `gorm:"column:purchase_price;size:50" json:"purchase_price,omitempty"`
	PurchaseInvNo  *string `gorm:"column:purchase_inv_no;size:50" json:"purchase_inv_no,omitempty"`
	ResellingPrice *string `gorm:"column:reselling_price;size:50" json:"reselling_price,omitempty"`

	DepositeAmount *string `gorm:"column:deposite_amount;size:50" json:"deposite_amount,omitempty"`

	ReserveDate   *string `gorm:"column:reserve_date;size:50" json:"reserve_date,omitempty"`
	ReserveRemark *string `gorm:"column:reserve_remark;size:50" json:"reserve_remark,omitempty"`

	OnholdDate   *string `gorm:"column:onhold_date;size:50" json:"onhold_date,omitempty"`
	OnholdRemark *string `gorm:"column:onhold_remark;size:50" json:"onhold_remark,omitempty"`

	PartsExchangeRegno       *string `gorm:"column:parts_exchange_regno;size:50" json:"parts_exchange_regno,omitempty"`
	PartsExchangeStockNumber *string `gorm:"column:parts_exchange_stock_number;size:50" json:"parts_exchange_stock_number,omitempty"`
	PartsExchangePrice       *string `gorm:"column:parts_exchange_price;size:50" json:"parts_exchange_price,omitempty"`
	VATType                  *string `gorm:"column:vat_type;size:50" json:"vat_type,omitempty"`
	VATReceipt               *string `gorm:"column:vat_receipt;size:50" json:"vat_receipt,omitempty"`
	VATReclaimable           *string `gorm:"column:vat_reclaimable;size:50" json:"vat_reclaimable,omitempty"`
	VAT                      *string `gorm:"column:vat;size:50" json:"vat,omitempty"`
	CalculationMode          *string `gorm:"column:calculation_mode;size:50" json:"calculation_mode,omitempty"`
	VATSubmissionStatus      *string `gorm:"column:vat_submission_status;size:50;default:'N'" json:"vat_submission_status,omitempty"`

	Status *int `gorm:"column:status;default:1" json:"status,omitempty"`
	// "Stock"=1, "Sold"=2, "Reserve"=3, "On Hold"=4, "Returned"=5, "Others"=6

	CreatedBy *string `gorm:"column:created_by;size:50" json:"created_by,omitempty"`
	UpdatedBy *string `gorm:"column:updated_by;size:50" json:"updated_by,omitempty"`

	CreatedAt           *time.Time `gorm:"column:created_at;default:CURRENT_TIMESTAMP" json:"created_at,omitempty"`
	UpdatedAt           *time.Time `gorm:"column:updated_at;default:CURRENT_TIMESTAMP on update CURRENT_TIMESTAMP" json:"updated_at,omitempty"`
	SupplierID          *int       `gorm:"column:supplier_id" json:"supplier_id,omitempty"`
	DepositCustomerName *string    `gorm:"column:deposit_customer_name;size:150" json:"deposit_customer_name,omitempty"`
	DepositCustomerID   *int       `gorm:"column:deposit_customer_id" json:"deposit_customer_id,omitempty"`
}

func (Stocks) TableName() string {
	return "stocks"
}

// type Overhead struct {
// 	OverheadID     int     `gorm:"column:overhead_id;primaryKey;autoIncrement" json:"overhead_id"`
// 	Period         *string `gorm:"column:period;size:50" json:"period,omitempty"`
// 	CostType       *string `gorm:"column:cost_type;size:50" json:"cost_type,omitempty"`
// 	CostDetails    *string `gorm:"column:cost_details;size:50" json:"cost_details,omitempty"`
// 	Cost           *string `gorm:"column:cost;size:50" json:"cost,omitempty"`
// 	VATReceipt     *string `gorm:"column:vat_receipt;size:50" json:"vat_receipt,omitempty"`
// 	VATReclaimable *string `gorm:"column:vat_reclaimable;size:50" json:"vat_reclaimable,omitempty"`
// 	VAT            *string `gorm:"column:vat;size:50" json:"vat,omitempty"`
// }

type Overhead struct {
	OverheadID          int        `gorm:"column:overhead_id;primaryKey;autoIncrement" json:"overhead_id"`
	OverheadKey         *string    `gorm:"column:overhead_key;size:50" json:"overhead_key,omitempty"`
	TenantInfoID        *int       `gorm:"column:tenant_info_id" json:"tenant_info_id,omitempty"`
	VATSubmissionStatus *string    `gorm:"column:vat_submission_status;size:50" json:"vat_submission_status"`
	OverheadTypeID      *int       `gorm:"column:overhead_type_id;size:50" json:"overhead_type_id,omitempty"`
	OverheadType        *string    `gorm:"column:overhead_type;size:50" json:"overhead_type,omitempty"`
	OverheadDetail      *string    `gorm:"column:overhead_detail;size:50" json:"overhead_detail,omitempty"`
	InvoiceNo           *string    `gorm:"column:invoice_no;size:50" json:"invoice_no,omitempty"`
	OverheadValue       *string    `gorm:"column:overhead_value;size:50" json:"overhead_value,omitempty"`
	CalculationMode     *string    `gorm:"column:calculation_mode;size:50" json:"calculation_mode,omitempty"`
	Date                *string    `gorm:"column:date;size:50" json:"date,omitempty"`
	Period              *string    `gorm:"column:period;size:50" json:"period,omitempty"`
	CostType            *string    `gorm:"column:cost_type;size:50" json:"cost_type,omitempty"`
	CostDetails         *string    `gorm:"column:cost_details;size:50" json:"cost_details,omitempty"`
	Cost                *string    `gorm:"column:cost;size:50" json:"cost,omitempty"`
	VATReceipt          *string    `gorm:"column:vat_receipt;size:50" json:"vat_receipt,omitempty"`
	VATReclaimable      *string    `gorm:"column:vat_reclaimable;size:50" json:"vat_reclaimable,omitempty"`
	VAT                 *string    `gorm:"column:vat;size:50" json:"vat,omitempty"`
	Status              *int       `gorm:"column:status" json:"status,omitempty"` // "Sold"=1, "Stock"=2, "Reserve"=3, "On Hold"=4, "Others"=5
	CreatedBy           *string    `gorm:"column:created_by;size:50" json:"created_by,omitempty"`
	UpdatedBy           *string    `gorm:"column:updated_by;size:50" json:"updated_by,omitempty"`
	CreatedAt           *time.Time `gorm:"column:created_at;default:CURRENT_TIMESTAMP" json:"created_at,omitempty"`
	UpdatedAt           *time.Time `gorm:"column:updated_at;default:CURRENT_TIMESTAMP on update CURRENT_TIMESTAMP" json:"updated_at,omitempty"`
}

type Cost struct {
	CostID              int        `gorm:"column:cost_id;primaryKey;autoIncrement" json:"cost_id"`
	TenantInfoID        *int       `gorm:"column:tenant_info_id" json:"tenant_info_id,omitempty"`
	VATSubmissionStatus *string    `gorm:"column:vat_submission_status;size:50" json:"vat_submission_status"`
	StockNumber         *string    `gorm:"column:stock_number" json:"stock_number,omitempty"`
	Registration        *string    `gorm:"column:registration" json:"registration,omitempty"`
	CostType            *string    `gorm:"column:cost_type" json:"cost_type,omitempty"`
	InvoiceNo           *string    `gorm:"column:invoice_no" json:"invoice_no,omitempty"`
	CostDetails         *string    `gorm:"column:cost_details" json:"cost_details,omitempty"`
	Cost                *string    `gorm:"column:cost" json:"cost,omitempty"`
	VATReceipt          *string    `gorm:"column:vat_receipt" json:"vat_receipt,omitempty"`
	CostDate            *string    `gorm:"column:cost_date" json:"cost_date,omitempty"`
	VATReclaimable      *string    `gorm:"column:vat_reclaimable" json:"vat_reclaimable,omitempty"`
	VAT                 *string    `gorm:"column:vat" json:"vat,omitempty"`
	Status              *int       `gorm:"column:status" json:"status,omitempty"` // "Sold"=1, "Stock"=2, "Reserve"=3, "On Hold"=4, "Others"=5
	CreatedBy           *string    `gorm:"column:created_by;size:50" json:"created_by,omitempty"`
	UpdatedBy           *string    `gorm:"column:updated_by;size:50" json:"updated_by,omitempty"`
	CreatedAt           *time.Time `gorm:"column:created_at;default:CURRENT_TIMESTAMP" json:"created_at,omitempty"`
	UpdatedAt           *time.Time `gorm:"column:updated_at;default:CURRENT_TIMESTAMP on update CURRENT_TIMESTAMP" json:"updated_at,omitempty"`
}

type OverheadDetails struct {
	OverheadDetailsID int        `gorm:"column:overhead_details_id;primaryKey;autoIncrement" json:"overhead_details_id"`
	OverheadID        *int       `gorm:"column:overhead_id;type:varchar(500);charset:utf8mb4;collation:utf8mb4_0900_ai_ci" json:"overhead_id,omitempty"`
	RegNo             *string    `gorm:"column:reg_no;type:varchar(50);charset:utf8mb4;collation:utf8mb4_0900_ai_ci" json:"reg_no,omitempty"`
	OverheadDetails   *string    `gorm:"column:overhead_details;type:varchar(500);charset:utf8mb4;collation:utf8mb4_0900_ai_ci" json:"overhead_details,omitempty"`
	Value             *string    `gorm:"column:value" json:"value,omitempty"`
	Overhead_type     *string    `gorm:"column:overhead_type;type:varchar(100);charset:utf8mb4;collation:utf8mb4_0900_ai_ci" json:"overhead_type,omitempty"`
	Comments          *string    `gorm:"column:comments;type:text;charset:utf8mb4;collation:utf8mb4_0900_ai_ci" json:"comments,omitempty"`
	Status            *int       `gorm:"column:status" json:"status,omitempty"` // "Sold"=1, "Stock"=2, "Reserve"=3, "On Hold"=4, "Others"=5
	CreatedBy         *string    `gorm:"column:created_by;size:50" json:"created_by,omitempty"`
	UpdatedBy         *string    `gorm:"column:updated_by;size:50" json:"updated_by,omitempty"`
	CreatedAt         *time.Time `gorm:"column:created_at;default:CURRENT_TIMESTAMP" json:"created_at,omitempty"`
	UpdatedAt         *time.Time `gorm:"column:updated_at;default:CURRENT_TIMESTAMP on update CURRENT_TIMESTAMP" json:"updated_at,omitempty"`
}
