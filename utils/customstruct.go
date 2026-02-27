package utils

import (
	"time"
)

// type App_roles struct {
// 	Approleid             int                     `gorm:"column:approleid;primaryKey"`
// 	Categoryname          string                  `gorm:"column:categoryname"`
// 	Rolename              string                  `gorm:"column:rolename"`
// 	Status                *int                    `gorm:"column:status"`
// 	App_permissions_roles []App_permissions_roles `gorm:"foreignKey:Approleid;references:Approleid"`
// }

// type App_permissions_roles struct {
// 	Apppermissionroleid      int                               `gorm:"column:apppermissionroleid"`
// 	Approleid                int                               `gorm:"column:approleid"`
// 	Rolename                 string                            `gorm:"column:rolename"`
// 	Appmoduleid              int                               `gorm:"column:appmoduleid"`
// 	Modulename               string                            `gorm:"column:modulename"`
// 	Taskid                   string                            `gorm:"column:taskid"`
// 	Qsfintegratortaskid      string                            `gorm:"column:qsfintegratortaskid"`
// 	Qsfscenarioruntaskid     string                            `gorm:"column:qsfscenarioruntaskid"`
// 	Status                   *int                              `gorm:"column:status"`
// 	App_permissions_features []models.App_permissions_features `gorm:"foreignKey:Appmoduleid;references:Appmoduleid"`
// }

// type Userinfos struct {
// 	Userid                int                     `json:"userid" gorm:"Primary_Key"`
// 	Authname              string                  `json:"authname"`
// 	Firstname             string                  `json:"firstname"`
// 	Lastname              string                  `json:"lastname"`
// 	Password              string                  `json:"password"`
// 	Email                 string                  `json:"email"`
// 	Dialcode              string                  `json:"dialcode"`
// 	Contactno             string                  `json:"contactno"`
// 	Configid              int                     `json:"configid"`
// 	Address               string                  `json:"address"`
// 	Suburb                string                  `json:"suburb"`
// 	City                  string                  `json:"city"`
// 	State                 string                  `json:"state"`
// 	Postcode              string                  `json:"postcode"`
// 	Tenantid              string                  `json:"tenantid"`
// 	Windowsid             string                  `json:"windowsid"`
// 	Partnerid             int                     `json:"partnerid"`
// 	Referenceid           int                     `json:"referenceid"`
// 	Approleid             int                     `json:"approleid"` // This is the actual column in the database
// 	Categoryname          string                  `gorm:"column:categoryname"`
// 	Rolename              string                  `gorm:"column:rolename"`
// 	Status                *int                    `json:"status"`
// 	App_permissions_roles []App_permissions_roles `gorm:"foreignKey:Approleid;references:Approleid"` // Maps roleid to approleid
// }

type Distinctstruct struct {
	DataType    string `json:"Datatype"`
	ActualFlag  string `json:"Actualflag"`
	Year        string `json:"Year"`
	Groupings   string `json:"Groupings"`
	BaseVersion string `json:"BaseVersion"`
	Type1       string `json:"Type1"`
	Type2       string `json:"Type2"`
	Type3       string `json:"Type3"`
	Type4       string `json:"Type4"`
	Type5       string `json:"Type5"`
	Type6       string `json:"Type6"`
	Region      string `json:"Region"`
	Currency    string `json:"Currency"`
	Class_Type  string `json:"Class_type"`
}

type App_inputstructure struct {
	Appinputstructureid int    `gorm:"column:appinputstructureid;primaryKey"`
	Tenantname          string `gorm:"column:tenantname"`
	Data_Type           string `gorm:"column:"datatype"`
	ActualFlag          string `gorm:"column:actualflag"`
	Year                string `gorm:"column:year"`
	Groupings           string `gorm:"column:groupings"`
	BaseVersion         string `gorm:"column:baseversion"`
	Type1               string `gorm:"column:type1"`
	Type2               string `gorm:"column:type2"`
	Type3               string `gorm:"column:type3"`
	Type4               string `gorm:"column:type4"`
	Type5               string `gorm:"column:type5"`
	Type6               string `gorm:"column:type6"`
	Region              string `gorm:"column:region"`
	Currency            string `gorm:"column:currency"`
	Classtype           string `gorm:"column:classtype"`
}

type Distinctdata struct {
	Groupings   string  `json:"Groupings"`
	BaseVersion string  `json:"BaseVersion"`
	Type1       string  `json:"Type1"`
	Type2       string  `json:"Type2"`
	Type3       string  `json:"Type3"`
	Type4       string  `json:"Type4"`
	Type5       string  `json:"Type5"`
	Type6       string  `json:"Type6"`
	Region      string  `json:"Region"`
	Currency    string  `json:"Currency"`
	Class_Type  string  `json:"Class_Type"`
	Year        string  `json:"Year"`
	Month       string  `json:"Month"`
	Period      string  `json:"Period"`
	Value       float64 `json:"Value"`
}

type App_inputdatas struct {
	Appinputdataid int    `gorm:"primaryKey"`
	Tenantname     string `gorm:"column:tenantname"`
	Groupings      string `gorm:"column:groupings"`
	BaseVersion    string `gorm:"column:baseversion"`
	Type1          string `gorm:"column:type1"`
	Type2          string `gorm:"column:type2"`
	Type3          string `gorm:"column:type3"`
	Type4          string `gorm:"column:type4"`
	Type5          string `gorm:"column:type5"`
	Type6          string `gorm:"column:type6"`
	Region         string `gorm:"column:region"`
	Currency       string `gorm:"column:currency"`
	Class_Type     string `gorm:"column:classtype"`
	Year           string `gorm:"column:year"`
	Month          string `gorm:"column:month"`
	Period         string `gorm:"column:period"`
	Value          string `gorm:"column:value"`
	// Created        string `gorm:"column:created"`
	// Updated        string `gorm:"column:updated"`
}

type EmailMessage struct {
	Message struct {
		Subject string `json:"subject"`
		Body    struct {
			ContentType string `json:"contentType"`
			Content     string `json:"content"`
		} `json:"body"`
		ToRecipients []struct {
			EmailAddress struct {
				Address string `json:"address"`
			} `json:"emailAddress"`
		} `json:"toRecipients"`
	} `json:"message"`
	SaveToSentItems string `json:"saveToSentItems"`
}

// Attributes struct to represent the main attributes of the app
type Attributes struct {
	ID                string                 `json:"id"`
	Name              string                 `json:"name"`
	Description       string                 `json:"description"`
	Thumbnail         string                 `json:"thumbnail"`
	LastReloadTime    string                 `json:"lastReloadTime"`
	CreatedDate       string                 `json:"createdDate"`
	ModifiedDate      string                 `json:"modifiedDate"`
	Owner             string                 `json:"owner"`
	OwnerID           string                 `json:"ownerId"`
	DynamicColor      string                 `json:"dynamicColor"`
	Published         bool                   `json:"published"`
	PublishTime       string                 `json:"publishTime"`
	Custom            map[string]interface{} `json:"custom"`
	HasSectionAccess  bool                   `json:"hasSectionAccess"`
	Encrypted         bool                   `json:"encrypted"`
	OriginAppID       string                 `json:"originAppId"`
	IsDirectQueryMode bool                   `json:"isDirectQueryMode"`
	Usage             string                 `json:"usage"`
	SpaceID           string                 `json:"spaceId"`
	ResourceType      string                 `json:"_resourcetype"`
}

// Response struct to represent the full response
type App_applications struct {
	Attributes Attributes `json:"attributes"`
	Privileges []string   `json:"privileges"`
	Create     []string   `json:"create"`
}

type AppResponse struct {
	Data []App_applications `json:"data"` // Assuming the array is under the "data" field
}

type AppSheets struct {
	SheetsForAllApps []SheetsForAllApps `json:"SheetsForAllApps"`
}

type SheetsForAllApps struct {
	AppID  string  `json:"appId"`
	Sheets []Sheet `json:"sheets"`
}

type Sheet struct {
	QInfo QInfo `json:"qInfo"`
	QMeta QMeta `json:"qMeta"`
	QData QData `json:"qData"`
}

type QInfo struct {
	QId   string `json:"qId"`
	QType string `json:"qType"`
}

type QMeta struct {
	Title        string   `json:"title"`
	Description  string   `json:"description"`
	CreatedDate  string   `json:"createdDate"`
	ModifiedDate string   `json:"modifiedDate"`
	Published    bool     `json:"published"`
	PublishTime  string   `json:"publishTime"`
	Approved     bool     `json:"approved"`
	Owner        string   `json:"owner"`
	OwnerID      string   `json:"ownerId"`
	Privileges   []string `json:"privileges"`
	ResourceType string   `json:"_resourcetype"`
	ObjectType   string   `json:"_objecttype"`
	ID           string   `json:"id"`
}

type QData struct {
	Rank        float64   `json:"rank"`
	Thumbnail   Thumbnail `json:"thumbnail"`
	Columns     int       `json:"columns"`
	Rows        int       `json:"rows"`
	Cells       []Cell    `json:"cells"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
}

type Thumbnail struct {
	QStaticContentUrl StaticContentURL `json:"qStaticContentUrl"`
}

type StaticContentURL struct {
	QUrl string `json:"qUrl"`
}

type Cell struct {
	Name    string `json:"name"`
	Type    string `json:"type"`
	Col     int    `json:"col"`
	Row     int    `json:"row"`
	Colspan int    `json:"colspan"`
	Rowspan int    `json:"rowspan"`
	Bounds  Bounds `json:"bounds"`
}

type Bounds struct {
	X      float64 `json:"x"`
	Y      float64 `json:"y"`
	Width  float64 `json:"width"`
	Height float64 `json:"height"`
}

type AssignedRole struct {
	ID    string `json:"id"`
	Level string `json:"level"`
	Name  string `json:"name"`
	Type  string `json:"type"`
}

type AssignedGroup struct {
	ID            string         `json:"id"`
	Name          string         `json:"name"`
	AssignedRoles []AssignedRole `json:"assignedRoles"`
}

type Links struct {
	Self struct {
		Href string `json:"href"`
	} `json:"self"`
}

type App_users struct {
	AssignedGroups []AssignedGroup `json:"assignedGroups"`
	AssignedRoles  []AssignedRole  `json:"assignedRoles"`
	Created        string          `json:"created"`
	CreatedAt      string          `json:"createdAt"`
	Email          string          `json:"email"`
	ID             string          `json:"id"`
	LastUpdated    string          `json:"lastUpdated"`
	LastUpdatedAt  string          `json:"lastUpdatedAt"`
	Links          Links           `json:"links"`
	Name           string          `json:"name"`
	Picture        string          `json:"picture"`
	Roles          []string        `json:"roles"`
	Status         string          `json:"status"`
	Subject        string          `json:"subject"`
	TenantID       string          `json:"tenantId"`
}

type UsersResponse struct {
	Data []App_users `json:"data"`
}

type App_user_response struct {
	Userid      int                   `gorm:"column:userid;primaryKey"`
	Authname    string                `gorm:"column:authname"`
	Firstname   string                `gorm:"column:firstname"`
	Lastname    string                `gorm:"column:lastname"`
	Password    string                `gorm:"column:password"`
	Otp         string                `gorm:"column:otp"`
	Email       string                `gorm:"column:email"`
	Approleid   int                   `gorm:"column:approleid"`
	Rolename    string                `gorm:"column:rolename"`
	Dialcode    string                `gorm:"column:dialcode"`
	Contactno   string                `gorm:"column:contactno"`
	Configid    int                   `gorm:"column:configid"`
	Address     string                `gorm:"column:address"`
	Suburb      string                `gorm:"column:suburb"`
	City        string                `gorm:"column:city"`
	State       string                `gorm:"column:state"`
	Postcode    string                `gorm:"column:postcode"`
	Tenantid    string                `gorm:"column:tenantid"`
	Tenantname  string                `gorm:"column:tenantname"`
	Subject     string                `gorm:"column:subject"`
	Picture     string                `gorm:"column:picture"`
	Referenceid int                   `gorm:"column:referenceid"`
	Status      int                   `gorm:"column:status"`
	Permissions []App_userpermissions `gorm:"foreignKey:Userid;references:Userid"`
	Tenants     []Tenants             `gorm:"foreignKey:Tenantid;references:Tenantid"`
}

type Tenants struct {
	Apptenantid      int    `json:"apptenantid"`
	Apptenantname    string `json:"apptenantname"`
	Tenantid         string `json:"tenantid"`
	Tenantname       string `json:"tenantname"`
	Tenanturl        string `json:"tenanturl"`
	Alternateurl     string `json:"alternateurl"`
	Keyid            string `json:"keyid"`
	Tenantemail      string `json:"tenantemail"`
	Tenantcontactno  string `json:"tenantcontactno"`
	Tenantaddress    string `json:"tenantaddress"`
	Tenantsuburb     string `json:"tenantsuburb"`
	Tenantcity       string `json:"tenantcity"`
	Tenantstate      string `json:"tenantstate"`
	Tenantpostcode   int    `json:"tenantpostcode"`
	Connectorid      int    `json:"connectorid"`
	Configid         int    `json:"configid"`
	Referenceid      int    `json:"referenceid"`
	Apikey           string `json:"apikey"`
	Webintegrationid string `json:"webintegrationid"`
	Status           *int   `json:"status" `
}

type App_spaces struct {
	Appspaceid int    `gorm:"column:appspaceid;primaryKey;autoIncrement"`
	Spaceid    string `gorm:"column:spaceid"`
	Type       string `gorm:"column:type"`
	Ownerid    string `gorm:"column:ownerid"`
	Name       string `gorm:"column:name"`
	Tenantid   string `gorm:"column:tenantid"`
	Tenantname string `gorm:"column:tenantname"`
	Status     int    `gorm:"column:status"`
	Createdby  string `gorm:"column:createdby"`
	Created    string `gorm:"column:created"`
	Updated    string `gorm:"column:updated"`
}

type App_userpermissions struct {
	Appuserpermissionsid int          `gorm:"primaryKey"`
	Userid               int          `gorm:"column:userid"`
	Username             string       `gorm:"column:username"`
	Apptenantid          int          `gorm:"column:apptenantid"`
	Apptenantname        string       `gorm:"column:apptenantname"`
	Tenantid             string       `gorm:"column:tenantid"`
	Tenantname           string       `gorm:"column:tenantname"`
	Spaceid              string       `gorm:"column:spaceid"`
	Appid                string       `gorm:"column:appid"`
	Appname              string       `gorm:"column:appname"`
	Status               *int         `gorm:"column:status"`
	Created              string       `gorm:"column:created"`
	Update               string       `gorm:"column:updated"`
	Spaces               []App_spaces `gorm:"foreignKey:Spaceid;references:Spaceid"`
	// Applications             []App_application          `gorm:"foreignKey:Appid;references:Appid"`
	App_userpermissionsheets []App_userpermissionsheets `gorm:"foreignKey:Appuserpermissionsid;references:Appuserpermissionsid" json:"App_userpermissionsheets"`
}

type App_userpermissionsheets struct {
	Appuserpermissionsheetsid int    `gorm:"primaryKey"`
	Appuserpermissionsid      int    `gorm:"column:appuserpermissionsid"`
	Appid                     string `gorm:"column:appid"`
	Appname                   string `gorm:"column:appname"`
	Sheetid                   string `gorm:"column:sheetid"`
	Sheetname                 string `gorm:"column:sheetname"`
	Status                    int    `gorm:"column:status;default:0"`
}

type App_application struct {
	Applicationid int          `gorm:"primaryKey"`
	Tenantid      string       `gorm:"column:tenantid"`
	Appid         string       `gorm:"column:appid"`
	Appname       string       `gorm:"column:appname"`
	Category      string       `gorm:"column:category"`
	Spaceid       string       `gorm:"column:spaceid"`
	Owner         string       `gorm:"column:owner"`
	Ownerid       string       `gorm:"column:ownerid"`
	Sensitivity   int          `gorm:"column:sensitivity"`
	Status        int          `gorm:"column:status"`
	Publishtime   string       `gorm:"column:publishtime"`
	App_sheets    []App_sheets `gorm:"foreignKey:Appid;references:Appid"`
}

type App_sheets struct {
	Appsheetid   int    `gorm:"column:appsheetid"`
	Sheetid      string `gorm:"column:sheetid"`
	Sheetname    string `gorm:"column:sheetname"`
	Tenantid     string `gorm:"column:tenantid"`
	Appid        string `gorm:"column:appid"`
	Appname      string `gorm:"column:appname"`
	Category     string `gorm:"column:category"`
	Createddate  string `gorm:"column:createddate"`
	Modifieddate string `gorm:"column:modifieddate"`
	Publishtime  string `gorm:"column:publishtime"`
	Ownerid      string `gorm:"column:ownerid"`
	Ownername    string `gorm:"column:ownername"`
	Status       *int   `gorm:"column:status"`
}

type VehicleData struct {
	ID                       uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	RegistrationNumber       string    `gorm:"column:registration_number;size:20;unique" json:"registration_number"`
	Make                     string    `gorm:"column:make;size:100" json:"make"`
	Model                    string    `gorm:"column:model;size:100" json:"model"`
	FuelType                 string    `gorm:"column:fuel_type;size:50" json:"fuel_type"`
	Colour                   string    `gorm:"column:colour;size:50" json:"colour"`
	MotStatus                string    `gorm:"column:mot_status;size:50" json:"mot_status"`
	MotExpiryDate            string    `gorm:"column:mot_expiry_date;size:50" json:"mot_expiry_date"`
	TaxStatus                string    `gorm:"column:tax_status;size:50" json:"tax_status"`
	TaxDueDate               string    `gorm:"column:tax_due_date;size:50" json:"tax_due_date"`
	YearOfManufacture        int       `gorm:"column:year_of_manufacture" json:"year_of_manufacture"`
	EngineCapacity           int       `gorm:"column:engine_capacity" json:"engine_capacity"`
	CO2Emissions             float64   `gorm:"column:co2_emissions" json:"co2_emissions"`
	Wheelplan                string    `gorm:"column:wheelplan;size:50" json:"wheelplan"`
	Make1                    string    `gorm:"column:make1;size:50" json:"make1"`
	Model1                   string    `gorm:"column:model1;size:50" json:"model1"`
	Description1             string    `gorm:"column:description1;size:50" json:"description1"`
	Colour1                  string    `gorm:"column:colour1;size:50" json:"colour1"`
	Mileage1                 string    `gorm:"column:mileage1;size:50" json:"mileage1"`
	YearOfManufacture1       int       `gorm:"column:year_of_manufacture1" json:"year_of_manufacture1"`
	FuelType1                string    `gorm:"column:fuel_type1;size:50" json:"fuel_type1"`
	EngineCapacity1          string    `gorm:"column:engine_capacity1;size:50" json:"engine_capacity1"`
	PurchasedDate1           string    `gorm:"column:purchased_date1;size:50" json:"purchased_date1"`
	PurchasedPrice1          string    `gorm:"column:purchased_price1;size:50" json:"purchased_price1"`
	PurchasedInvoiceNo1      string    `gorm:"column:purchased_invoice_no1;size:50" json:"purchased_invoice_no1"`
	ResellingPrice1          string    `gorm:"column:reselling_price1;size:50" json:"reselling_price1"`
	PurchaseType1            string    `gorm:"column:purchase_type1;size:50" json:"purchase_type1"`
	PartsExchangeRegNo1      string    `gorm:"column:parts_exchange_reg_no1;size:50" json:"parts_exchange_reg_no1"`
	PartsExchangeStockNo1    string    `gorm:"column:parts_exchange_stock_no1;size:50" json:"parts_exchange_stock_no1"`
	MarkedForExport          bool      `gorm:"column:marked_for_export" json:"marked_for_export"`
	MonthOfFirstRegistration string    `gorm:"column:month_of_first_registration;size:50" json:"month_of_first_registration"`
	FirstUsedDate            string    `gorm:"column:first_used_date;size:50" json:"first_used_date"`
	RegistrationDate         string    `gorm:"column:registration_date;size:50" json:"registration_date"`
	ManufactureDate          string    `gorm:"column:manufacture_date;size:50" json:"manufacture_date"`
	EngineSize               string    `gorm:"column:engine_size;size:50" json:"engine_size"`
	HasOutstandingRecall     string    `gorm:"column:has_outstanding_recall;size:20" json:"has_outstanding_recall"`
	VehicleDataUpdated       int       `gorm:"column:vechile_data_updated" json:"vechile_data_updated"`
	CreatedAt                time.Time `gorm:"column:created_at;autoCreateTime" json:"created_at"`
}
