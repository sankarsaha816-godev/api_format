package models

import "time"

// VehicleValuation stores vehicle valuation data from external API
type VehicleValuation struct {
	ID                    int64      `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	TenantInfoID          *int64     `gorm:"column:tenant_info_id;index" json:"tenant_info_id,omitempty"`
	VehicleRegistration   *string    `gorm:"column:vehicle_registration;size:50;index" json:"vehicle_registration,omitempty"`
	CurrentMileage        *int       `gorm:"column:current_mileage" json:"current_mileage,omitempty"`
	
	// Valuation response fields
	ValuationID           *string    `gorm:"column:valuation_id;size:100" json:"valuation_id,omitempty"`
	Make                  *string    `gorm:"column:make;size:100" json:"make,omitempty"`
	Model                 *string    `gorm:"column:model;size:100" json:"model,omitempty"`
	// Colour                *string    `gorm:"column:colour;size:100" json:"colour,omitempty"`
	// TransmissionType      *string    `gorm:"column:transmission_type;size:50" json:"transmission_type,omitempty"`
	// FuelType              *string    `gorm:"column:fuel_type;size:50" json:"fuel_type,omitempty"`
	// EngineSize            *string    `gorm:"column:engine_size;size:50" json:"engine_size,omitempty"`
	// BodyType              *string    `gorm:"column:body_type;size:100" json:"body_type,omitempty"`
	Doors                 *int       `gorm:"column:doors" json:"doors,omitempty"`
	Seats                 *int       `gorm:"column:seats" json:"seats,omitempty"`
	YearOfManufacture     *int       `gorm:"column:year_of_manufacture" json:"year_of_manufacture,omitempty"`
	
	// Valuation amounts
	// EstimatedValue        *float64   `gorm:"column:estimated_value;type:decimal(10,2)" json:"estimated_value,omitempty"`
	// ValuationLow          *float64   `gorm:"column:valuation_low;type:decimal(10,2)" json:"valuation_low,omitempty"`
	// ValuationHigh         *float64   `gorm:"column:valuation_high;type:decimal(10,2)" json:"valuation_high,omitempty"`

	// Retail valuation amounts
	RetailLowValuation     *float64  `gorm:"column:retail_low_valuation;type:decimal(10,2)" json:"retail_low_valuation,omitempty"`
	RetailAverageValuation *float64  `gorm:"column:retail_average_valuation;type:decimal(10,2)" json:"retail_average_valuation,omitempty"`
	RetailHighValuation    *float64  `gorm:"column:retail_high_valuation;type:decimal(10,2)" json:"retail_high_valuation,omitempty"`

	// Trade valuation amounts
	TradeLowValuation      *float64  `gorm:"column:trade_low_valuation;type:decimal(10,2)" json:"trade_low_valuation,omitempty"`
	TradeAverageValuation  *float64  `gorm:"column:trade_average_valuation;type:decimal(10,2)" json:"trade_average_valuation,omitempty"`
	TradeHighValuation     *float64  `gorm:"column:trade_high_valuation;type:decimal(10,2)" json:"trade_high_valuation,omitempty"`
	
	// Additional data
	Depreciation          *float64   `gorm:"column:depreciation;type:decimal(10,2)" json:"depreciation,omitempty"`
	Mileage               *int       `gorm:"column:mileage" json:"mileage,omitempty"`
	BregoDerivativeID     *string    `gorm:"column:brego_derivative_id;size:100" json:"brego_derivative_id,omitempty"`
	VehicleIdentificationNumber *string `gorm:"column:vehicle_identification_number;size:100" json:"vehicle_identification_number,omitempty"`
	CurrencyUnit          *string    `gorm:"column:currency_unit;size:10" json:"currency_unit,omitempty"`
	MileageUnit           *string    `gorm:"column:mileage_unit;size:10" json:"mileage_unit,omitempty"`
	Version               *string    `gorm:"column:version;size:50" json:"version,omitempty"`
	VehicleDesc           *string    `gorm:"column:vehicle_desc;size:255" json:"vehicle_desc,omitempty"`
	VehicleType           *string    `gorm:"column:vehicle_type;size:50" json:"vehicle_type,omitempty"`
	ApiResponseRaw        *string    `gorm:"column:api_response_raw;type:longtext" json:"api_response_raw,omitempty"`
	
	CreatedAt             *time.Time `gorm:"column:created_at;default:CURRENT_TIMESTAMP" json:"created_at,omitempty"`
	UpdatedAt             *time.Time `gorm:"column:updated_at;default:CURRENT_TIMESTAMP on update CURRENT_TIMESTAMP" json:"updated_at,omitempty"`
}

// TableName specifies the table name for VehicleValuation
func (VehicleValuation) TableName() string {
	return "vehicle_valuations"
}

// VehicleDesirability stores vehicle desirability rating data from external API
type VehicleDesirability struct {
	ID                           int64      `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	TenantInfoID                 *int64     `gorm:"column:tenant_info_id;index" json:"tenant_info_id,omitempty"`
	VehicleRegistration          *string    `gorm:"column:vehicle_registration;size:50;index" json:"vehicle_registration,omitempty"`
	
	// Desirability ratings
	DesirabilityRating           *int       `gorm:"column:desirability_rating" json:"desirability_rating,omitempty"`
	DesirabilityRatingLastMonth  *int       `gorm:"column:desirability_rating_last_month" json:"desirability_rating_last_month,omitempty"`
	
	// Regional data (national)
	Region                       *string    `gorm:"column:region;size:50" json:"region,omitempty"`
	AverageAdvertDays            *int       `gorm:"column:average_advert_days" json:"average_advert_days,omitempty"`
	AverageAdvertMilesPerMonth   *int       `gorm:"column:average_advert_miles_per_month" json:"average_advert_miles_per_month,omitempty"`
	AverageAdvertPriceChange     *float64   `gorm:"column:average_advert_price_change;type:decimal(10,2)" json:"average_advert_price_change,omitempty"`
	AverageManufacturerMilesPerMonth *int   `gorm:"column:average_manufacturer_miles_per_month" json:"average_manufacturer_miles_per_month,omitempty"`
	AverageModelRangeMilesPerMonth *int     `gorm:"column:average_model_range_miles_per_month" json:"average_model_range_miles_per_month,omitempty"`
	AverageModelMilesPerMonth    *int       `gorm:"column:average_model_miles_per_month" json:"average_model_miles_per_month,omitempty"`
	
	// Price and mileage ranges
	MinAdvertPrice               *float64   `gorm:"column:min_advert_price;type:decimal(10,2)" json:"min_advert_price,omitempty"`
	MaxAdvertPrice               *float64   `gorm:"column:max_advert_price;type:decimal(10,2)" json:"max_advert_price,omitempty"`
	MinAdvertMileage             *int       `gorm:"column:min_advert_mileage" json:"min_advert_mileage,omitempty"`
	MaxAdvertMileage             *int       `gorm:"column:max_advert_mileage" json:"max_advert_mileage,omitempty"`
	
	// Market data
	WeeklyAdvertQtyChange        *int       `gorm:"column:weekly_advert_qty_change" json:"weekly_advert_qty_change,omitempty"`
	AdvertQty                    *int       `gorm:"column:advert_qty" json:"advert_qty,omitempty"`
	AverageAdvertOwnersPerAnnum  *float64   `gorm:"column:average_advert_owners_per_annum;type:decimal(5,2)" json:"average_advert_owners_per_annum,omitempty"`
	MinAdvertPreviousOwners      *int       `gorm:"column:min_advert_previous_owners" json:"min_advert_previous_owners,omitempty"`
	MaxAdvertPreviousOwners      *int       `gorm:"column:max_advert_previous_owners" json:"max_advert_previous_owners,omitempty"`
	
	// Colour distribution (JSON array stored as text)
	VehicleColours               *string    `gorm:"column:vehicle_colours;type:longtext" json:"vehicle_colours,omitempty"`
	
	// Vehicle information
	VehicleDesc                  *string    `gorm:"column:vehicle_desc;size:255" json:"vehicle_desc,omitempty"`
	VehicleType                  *string    `gorm:"column:vehicle_type;size:50" json:"vehicle_type,omitempty"`
	
	// Raw response
	ApiResponseRaw               *string    `gorm:"column:api_response_raw;type:longtext" json:"api_response_raw,omitempty"`
	
	CreatedAt                    *time.Time `gorm:"column:created_at;default:CURRENT_TIMESTAMP" json:"created_at,omitempty"`
	UpdatedAt                    *time.Time `gorm:"column:updated_at;default:CURRENT_TIMESTAMP on update CURRENT_TIMESTAMP" json:"updated_at,omitempty"`
}

// TableName specifies the table name for VehicleDesirability
func (VehicleDesirability) TableName() string {
	return "vehicle_desirability"
}

// VehicleMarketValue stores vehicle fair market value data from external API
type VehicleMarketValue struct {
	ID                          int64      `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	TenantInfoID                *int64     `gorm:"column:tenant_info_id;index" json:"tenant_info_id,omitempty"`
	VehicleRegistration         *string    `gorm:"column:vehicle_registration;size:50;index" json:"vehicle_registration,omitempty"`
	CurrentMileage              *int       `gorm:"column:current_mileage" json:"current_mileage,omitempty"`
	PostCode                    *string    `gorm:"column:post_code;size:20" json:"post_code,omitempty"`
	RadiusMiles                 *int       `gorm:"column:radius_miles" json:"radius_miles,omitempty"`

	// Market value data
	LocalFairMarketValue        *int       `gorm:"column:local_fair_market_value" json:"local_fair_market_value,omitempty"`
	NationalFairMarketValue     *int       `gorm:"column:national_fair_market_value" json:"national_fair_market_value,omitempty"`
	LocalAverageDaysToSold      *int       `gorm:"column:local_average_days_to_sold" json:"local_average_days_to_sold,omitempty"`
	NationalAverageDaysToSold   *int       `gorm:"column:national_average_days_to_sold" json:"national_average_days_to_sold,omitempty"`

	// Raw response
	ApiResponseRaw              *string    `gorm:"column:api_response_raw;type:longtext" json:"api_response_raw,omitempty"`

	CreatedAt                   *time.Time `gorm:"column:created_at;default:CURRENT_TIMESTAMP" json:"created_at,omitempty"`
	UpdatedAt                   *time.Time `gorm:"column:updated_at;default:CURRENT_TIMESTAMP on update CURRENT_TIMESTAMP" json:"updated_at,omitempty"`
}

// TableName specifies the table name for VehicleMarketValue
func (VehicleMarketValue) TableName() string {
	return "vehicle_market_values"
}
