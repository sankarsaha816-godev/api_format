package domain

import (
	"bitbucket/api_format/db"
	"bitbucket/api_format/models"
	"bitbucket/api_format/utils"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
)

func Createtenant(u models.Tenants) (int, error) {
	tx := db.DB.Begin()
	fmt.Println("--success--")

	// Copy required fields from User to Tenant
	t1 := tx.Table("tenants").Create(&u)
	if t1.Error != nil {
		fmt.Println("t1.err", tx.Error)
		tx.Rollback()
	}

	err1 := tx.Commit().Error

	if err1 != nil {
		print("err1", err1)
		panic(err1)
	}

	var lastInsertedTenant models.Tenants
	if err := db.DB.Last(&lastInsertedTenant).Error; err != nil {
		fmt.Println("Error fetching last inserted record:", err)
		return 0, err
	}
	tenantID := lastInsertedTenant.Apptenantid
	fmt.Println("Last inserted Tenant ID:", tenantID)
	return tenantID, nil
}

func GetTenantId(contactno string) int {

	var tid int

	q1 := `select apptenantid from tenants WHERE tenantcontactno='` + contactno + `'`

	db.DB.Raw(q1).Find(&tid)

	return tid

}

// func Getticket(windowsid string) string {
// 	client := &http.Client{}

// 	str := windowsid
// 	fmt.Println("----------str---------", str)
// 	parts := strings.Split(str, "\\")

// 	firstPart := parts[0]
// 	secondPart := parts[len(parts)-1]
// 	fmt.Println("First part:", firstPart)
// 	fmt.Println("Second part:", secondPart)

// 	// _, _, _, getticket := utils.Urls()

// 	// create a new GET request to the third-party API
// 	url := "https://qea-web.insightdelivered.com:23002/api/ticket/id/demouser4"
// 	// + firstPart + "/" + secondPart
// 	fmt.Println("===== url =====", url)
// 	req, err := http.NewRequest("GET", url, nil)
// 	if err != nil {
// 		fmt.Println("Error creating request:", err)
// 	}

// 	// send the request using the client
// 	resp, err := client.Do(req)
// 	if err != nil {
// 		fmt.Println("Error sending request:", err)
// 	}
// 	defer resp.Body.Close()

// 	// read the response body
// 	responseBody, err := ioutil.ReadAll(resp.Body)
// 	if err != nil {
// 		fmt.Println("Error reading response body:", err)
// 	}

// 	// print the response body
// 	fmt.Println(string(responseBody))

// 	output := string(responseBody)

// 	jsonString := output

// 	// Parse the JSON string into a map[string]interface{}
// 	var data map[string]interface{}

// 	err1 := json.Unmarshal([]byte(jsonString), &data)
// 	if err1 != nil {
// 		fmt.Println("Error unmarshaling JSON:", err)
// 	}

// 	// Access each key-value pair using the corresponding key
// 	message := data["message"]
// 	ticket := data["ticket"]

// 	// Print the values of each key
// 	fmt.Println("message:", message)
// 	fmt.Println("ticket:", ticket)
// 	ticketString := ticket.(string)

// 	return ticketString

// }

func Getticket(windowsid string) (string, error) {
	_, certUrl, _, getticket := utils.Urls()

	// Create an HTTP client to download the certificate
	httpClient := &http.Client{}

	// Download the CA certificate from the provided URL
	resp, err := httpClient.Get(certUrl)
	if err != nil {
		return "", fmt.Errorf("failed to download CA certificate: %v", err)
	}
	defer resp.Body.Close()

	// Read the downloaded certificate into a byte slice
	caCert, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read CA certificate: %v", err)
	}

	// Create a CA certificate pool and add the downloaded QlikSense CA
	caCertPool := x509.NewCertPool()
	if !caCertPool.AppendCertsFromPEM(caCert) {
		return "", fmt.Errorf("failed to append CA certificate")
	}

	// Create a custom TLS configuration using the CA pool
	tlsConfig := &tls.Config{
		RootCAs: caCertPool,
	}

	// Create an HTTP client with the custom TLS configuration
	transport := &http.Transport{
		TLSClientConfig: tlsConfig,
	}
	client := &http.Client{
		Transport: transport,
	}

	// Parse the Windows ID
	parts := strings.Split(windowsid, "\\")
	if len(parts) < 2 {
		return "", fmt.Errorf("invalid Windows ID format")
	}
	firstPart := parts[0]
	secondPart := parts[len(parts)-1]

	// Create the URL for the GET request to the third-party API
	url := getticket + firstPart + "/" + secondPart
	fmt.Println("===== url =====", url)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", fmt.Errorf("error creating request: %v", err)
	}

	// Send the request using the custom HTTP client
	resp, err = client.Do(req)
	if err != nil {
		return "", fmt.Errorf("error sending request: %v", err)
	}
	defer resp.Body.Close()

	// Read the response body
	responseBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("error reading response body: %v", err)
	}

	// Parse the JSON response
	var data map[string]interface{}
	if err := json.Unmarshal(responseBody, &data); err != nil {
		return "", fmt.Errorf("error unmarshaling JSON: %v", err)
	}

	// Extract and return the ticket
	ticket, ok := data["ticket"].(string)
	if !ok || ticket == "" {
		return "", fmt.Errorf("no ticket available for the user")
	}

	return ticket, nil
}

func Gettenantbyname(contactno string) int {

	var tid int
	// var tname string

	q1 := `select apptenantid from tenants WHERE tenantname='` + contactno + `'`

	db.DB.Raw(q1).Find(&tid)

	fmt.Println("=======tid======", tid)

	return tid

}

func Gettenantidbytid(contactno int) string {

	var tid string

	q1 := `select tenantid from tenants WHERE apptenantid = ` + strconv.Itoa(contactno)

	db.DB.Raw(q1).Find(&tid)

	fmt.Println("=======tid======", tid)

	return tid

}

func Gettenantappsbytid(contactno string) []utils.App_application {

	var tid []utils.App_application

	q1 := `SELECT * FROM app_applications WHERE STATUS = 1 AND tenantid ='` + contactno + `' `

	db.DB.Raw(q1).Find(&tid)

	// fmt.Println("=======tid======", tid)

	return tid

}

// Assuming `db` is your GORM database instance
func Gettenanturlbytid(contactno string) (string, string, error) {
	var url, key string

	// Use parameterized query to prevent SQL injection
	query := "SELECT alternateurl, apikey FROM tenants WHERE tenantid = ?"
	result := db.DB.Raw(query, contactno).Row()

	// Scan results into variables
	if err := result.Scan(&url, &key); err != nil {
		return "", "", fmt.Errorf("failed to retrieve tenant by tenant ID: %v", err)
	}

	fmt.Println("=======Tenant URL======", url)
	fmt.Println("=======API Key======", key)

	return url, key, nil
}

func Gettenantbyid(tid int) models.Tenants {

	var tenants models.Tenants

	q1 := `SELECT * from tenants WHERE apptenantid = ` + strconv.Itoa(tid)

	db.DB.Raw(q1).Find(&tenants)

	return tenants

}

func Deletetenantbyid(uid int) (bool, error) {

	q1 := `DELETE FROM tenants WHERE apptenantid = ?`

	result := db.DB.Exec(q1, uid)

	if result.Error != nil {
		return false, result.Error
	}

	return true, nil
}

func Getalltenants() []models.Tenants {

	var tenants []models.Tenants

	q1 := `SELECT * from tenants `

	db.DB.Raw(q1).Find(&tenants)

	return tenants

}

func Checktenantpermissions(tenants []models.Tenantpermissions) int {
	var tenant models.Tenantpermissions

	// Loop through each tenant permission to check its existence
	for _, t := range tenants {
		result := db.DB.Where("tenantpermissionsid = ?", t.Tenantpermissionsid).First(&tenant)
		if result.Error == nil && result.RowsAffected > 0 {
			return tenant.Tenantpermissionsid
		}
	}
	return 0
}

func Createtenantpermissions(tenants []models.Tenantpermissions) (int, error) {
	// Begin a transaction
	tx := db.DB.Begin()

	if err := tx.Error; err != nil {
		return 0, err
	}

	// Insert the tenant permissions into the database
	if err := tx.Create(&tenants).Error; err != nil {
		tx.Rollback() // Rollback the transaction in case of error
		return 0, err
	}

	// Commit the transaction
	if err := tx.Commit().Error; err != nil {
		return 0, err
	}

	// we can return the Tenantid of the first entry
	if len(tenants) > 0 {
		return tenants[0].Apptenantid, nil
	}

	return 0, nil
}

func Getpermissionsbyid(tid int) []models.Tenantpermissions {

	var tenants []models.Tenantpermissions

	q1 := `SELECT * from tenantpermissions WHERE apptenantid = ` + strconv.Itoa(tid)

	db.DB.Raw(q1).Find(&tenants)

	return tenants

}

// func Gettenantpermissionmodules(tid string) []models.App_modules {

// 	var tenants []models.App_modules

// 	q1 := `SELECT DISTINCT b.* FROM app_permissions_features a INNER JOIN app_modules b ON a.appmoduleid = b.appmoduleid
// WHERE a.appid = '` + tid + `'`

// 	db.DB.Raw(q1).Find(&tenants)

// 	return tenants

// }

// func Getpermissionsbyid(tid int) ([]models.Tenantpermissions, error) {
// 	var tenants []models.Tenantpermissions

// 	// Build the query using GORM with proper joins
// 	err := db.DB.Table("tenantpermissions as a").
// 		Select("a.*, b.*").
// 		Joins("INNER JOIN app_permissions_features b ON a.appid = b.appid").
// 		// Joins("INNER JOIN app_modules c ON c.appmoduleid = b.appmoduleid").
// 		Where("a.tenantid = ?", tid).
// 		Preload("app_permissions_features").
// 		Find(&tenants).Error

// 	if err != nil {
// 		return nil, err
// 	}

// 	return tenants, nil
// }

func Getavailablepermissionsbyid(tid int) []utils.App_application {

	var tenants []utils.App_application

	q1 := `SELECT * FROM app_applications WHERE STATUS = 1 AND appid NOT IN (SELECT appid FROM tenantpermissions WHERE apptenantid =` + strconv.Itoa(tid) + `) `

	db.DB.Raw(q1).Find(&tenants)

	return tenants

}

func Gettenantuser(userid int) models.Tenantuserinfo {

	var user models.Tenantuserinfo

	// Retrieve the user with related permissions from the database
	if err := db.DB.Preload("Tenantpermissions").First(&user, userid).Error; err != nil {
		panic(err)
	}

	// Return the user information as a JSON response
	return user
}
