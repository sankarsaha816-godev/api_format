package models

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

type Tenantpermissions struct {
	Tenantpermissionsid int    `json:"tenantpermissionsid"`
	Apptenantid         int    `json:"apptenantid"`
	Apptenantname       string `json:"apptenantname"`
	Tenantid            string `json:"tenantid"`
	Tenantname          string `json:"tenantname"`
	Appid               string `json:"appid"`
	Appname             string `json:"appname"`
	Category            string `json:"category"`
	Publish             string `json:"publish"`
	Status              *int   `json:"status"`
	// App_modules         []App_modules `gorm:"foreignKey:Appid;references:Appid" json:"App_modules"`
}

// type Tenantpermissions struct {
// 	Tenantpermissionsid int           `json:"tenantpermissionsid"`
// 	Tenantid            int           `json:"tenantid"`
// 	Tenantname          string        `json:"tenantname"`
// 	Appid               string        `json:"appid"`
// 	Appname             string        `json:"appname"`
// 	Category            string        `json:"category"`
// 	Publish             string        `json:"publish"`
// 	Status              *int          `json:"status"`
// 	App_modules         []App_modules `gorm:"many2many:app_permissions_features;foreignKey:Appid;joinForeignKey:appid;References:Appmoduleid;joinReferences:appmoduleid" json:"App_modules"`
// }

type Tenantuserinfo struct {
	Userid      int    `json:"userid" gorm:"Primary_Key"`
	Authname    string `json:"authname"`
	Firstname   string `json:"firstname"`
	Lastname    string `json:"lastname"`
	Password    string `json:"password"`
	Email       string `json:"email"`
	Dialcode    string `json:"dialcode"`
	Contactno   string `json:"contactno"`
	Configid    int    `json:"configid"`
	Address     string `json:"address"`
	Suburb      string `json:"suburb"`
	City        string `json:"city"`
	State       string `json:"state"`
	Postcode    string `json:"postcode"`
	Tenantid    int    `json:"tenantid"`
	Windowsid   string `json:"windowsid"`
	Partnerid   int    `json:"partnerid"`
	Referenceid int    `json:"referenceid"`
	Roleid      int    `json:"roleid"`
	Status      *int   `json:"status"`
	// Tenantpermissions []App_permissions `json:"tenantpermissions"`
}
