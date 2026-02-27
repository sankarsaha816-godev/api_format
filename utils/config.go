package utils

import (
	"fmt"
	"os"

	"github.com/spf13/viper"
	"golang.org/x/crypto/bcrypt"
)

// OutlookConfig holds Outlook/Microsoft email configuration
type OutlookConfig struct {
	ClientID     string
	ClientSecret string
	TenantID     string
	Mailbox      string
}

// GetOutlookConfig loads Outlook configuration from config.yaml
func GetOutlookConfig(env string) OutlookConfig {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")

	err := viper.ReadInConfig()
	if err != nil {
		fmt.Println("fatal error config file: default \n", err)
		os.Exit(1)
	}

	envPrefix := env + "."
	return OutlookConfig{
		ClientID:     viper.GetString(envPrefix + "OUTLOOK_CLIENT_ID"),
		ClientSecret: viper.GetString(envPrefix + "OUTLOOK_CLIENT_SECRET"),
		TenantID:     viper.GetString(envPrefix + "OUTLOOK_TENANT_ID"),
		Mailbox:      viper.GetString(envPrefix + "OUTLOOK_MAILBOX"),
	}
}


func DevConfig() (Port, DBname, Password, Username, Host, key, secret, clientid, clientsecret, scopes, tokenurl string) {

	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")

	err := viper.ReadInConfig()
	if err != nil {
		fmt.Println("fatal error config file: default \n", err)
		os.Exit(1)
	}

	dbname := viper.GetString("DEV.A1")
	dbpassword := viper.GetString("DEV.A2")
	dbusername := viper.GetString("DEV.A3")
	dbport := viper.GetString("DEV.A4")
	dbhost := viper.GetString("DEV.A5")
	contextkey := viper.GetString("DEV.A6")
	jwtkey := viper.GetString("DEV.A7")
	clientid = viper.GetString("DEV.A8")
	clientsecret = viper.GetString("DEV.A9")
	scopes = viper.GetString("DEV.A10")
	tokenurl = viper.GetString("DEV.A11")
	return dbport, dbname, dbpassword, dbusername, dbhost, contextkey, jwtkey, clientid, clientsecret, scopes, tokenurl
}

func LiveConfig() (Port, DBname, Password, Username, Host, key, secret, clientid, clientsecret, scopes, tokenurl string) {

	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")

	err := viper.ReadInConfig()
	if err != nil {
		fmt.Println("fatal error config file: default \n", err)
		os.Exit(1)
	}
	dbname := viper.GetString("LIVE.A1")
	dbpassword := viper.GetString("LIVE.A2")
	dbusername := viper.GetString("LIVE.A3")
	dbport := viper.GetString("LIVE.A4")
	dbhost := viper.GetString("LIVE.A5")
	contextkey := viper.GetString("LIVE.A6")
	jwtkey := viper.GetString("LIVE.A7")
	clientid = viper.GetString("LIVE.A8")
	clientsecret = viper.GetString("LIVE.A9")
	scopes = viper.GetString("LIVE.A10")
	tokenurl = viper.GetString("LIVE.A11")
	return dbport, dbname, dbpassword, dbusername, dbhost, contextkey, jwtkey, clientid, clientsecret, scopes, tokenurl
}

func ProdConfig() (Port, DBname, Password, Username, Host, key, secret, clientid, clientsecret, scopes, tokenurl string) {

	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")

	err := viper.ReadInConfig()
	if err != nil {
		fmt.Println("fatal error config file: default \n", err)
		os.Exit(1)
	}
	dbname := viper.GetString("PROD.A1")
	dbpassword := viper.GetString("PROD.A2")
	dbusername := viper.GetString("PROD.A3")
	dbport := viper.GetString("PROD.A4")
	dbhost := viper.GetString("PROD.A5")
	contextkey := viper.GetString("PROD.A6")
	jwtkey := viper.GetString("PROD.A7")
	clientid = viper.GetString("PROD.A8")
	clientsecret = viper.GetString("PROD.A9")
	scopes = viper.GetString("PROD.A10")
	tokenurl = viper.GetString("PROD.A11")
	return dbport, dbname, dbpassword, dbusername, dbhost, contextkey, jwtkey, clientid, clientsecret, scopes, tokenurl
}


func AutoConfig() (Port, DBname, Password, Username, Host, key, secret, clientid, clientsecret, scopes, tokenurl string) {

	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")

	err := viper.ReadInConfig()
	if err != nil {
		fmt.Println("fatal error config file: default \n", err)
		os.Exit(1)
	}
	dbname := viper.GetString("AUTO_PROD.A1")
	dbpassword := viper.GetString("AUTO_PROD.A2")
	dbusername := viper.GetString("AUTO_PROD.A3")
	dbport := viper.GetString("AUTO_PROD.A4")
	dbhost := viper.GetString("AUTO_PROD.A5")
	contextkey := viper.GetString("AUTO_PROD.A6")
	jwtkey := viper.GetString("AUTO_PROD.A7")
	clientid = viper.GetString("AUTO_PROD.A8")
	clientsecret = viper.GetString("AUTO_PROD.A9")
	scopes = viper.GetString("AUTO_PROD.A10")
	tokenurl = viper.GetString("AUTO_PROD.A11")
	return dbport, dbname, dbpassword, dbusername, dbhost, contextkey, jwtkey, clientid, clientsecret, scopes, tokenurl
}


func AutoConfigDEV() (Port, DBname, Password, Username, Host, key, secret, clientid, clientsecret, scopes, tokenurl string) {

	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")

	err := viper.ReadInConfig()
	if err != nil {
		fmt.Println("fatal error config file: default \n", err)
		os.Exit(1)
	}
	dbname := viper.GetString("AUTO_DEV.A1")
	dbpassword := viper.GetString("AUTO_DEV.A2")
	dbusername := viper.GetString("AUTO_DEV.A3")
	dbport := viper.GetString("AUTO_DEV.A4")
	dbhost := viper.GetString("AUTO_DEV.A5")
	contextkey := viper.GetString("AUTO_DEV.A6")
	jwtkey := viper.GetString("AUTO_DEV.A7")
	clientid = viper.GetString("AUTO_DEV.A8")
	clientsecret = viper.GetString("AUTO_DEV.A9")
	scopes = viper.GetString("AUTO_DEV.A10")
	tokenurl = viper.GetString("AUTO_DEV.A11")
	return dbport, dbname, dbpassword, dbusername, dbhost, contextkey, jwtkey, clientid, clientsecret, scopes, tokenurl
}


func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

// CheckPassword hash compares raw password with it's hashed values
func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	if err != nil {
		return false
	}
	return true
}

func Urls() (Port, CertFilePath, Keyfilepath, Getticket string) {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()
	if err != nil {
		fmt.Println("fatal error config file: default \n", err)
		os.Exit(1)
	}
	// Declare var
	port := viper.GetString("DEV.B1")
	certFilePath := viper.GetString("DEV.B2")
	keyfilepath := viper.GetString("DEV.B3")
	getticket := viper.GetString("DEV.B4")

	return port, certFilePath, keyfilepath, getticket
}

func Qlikdata() (Applicationsurl, Sheeturl, Xrfkey, Username, Password string) {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")

	err := viper.ReadInConfig()
	if err != nil {
		fmt.Println("fatal error config file: default \n", err)
		os.Exit(1)
	}
	applicationsurl := viper.GetString("DEV.B5")
	sheeturl := viper.GetString("DEV.B9")
	xrfkey := viper.GetString("DEV.B6")
	username := viper.GetString("DEV.B7")
	password := viper.GetString("DEV.B8")
	return applicationsurl, sheeturl, xrfkey, username, password
}

func Folders() (Landing, Projects string) {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")

	err := viper.ReadInConfig()
	if err != nil {
		fmt.Println("fatal error config file: default \n", err)
		os.Exit(1)
	}
	landing := viper.GetString("DEV.C1")
	projects := viper.GetString("DEV.C2")

	return landing, projects
}



func Urls3() (Clientid, Clientsecret, Scopes, Tokenurl, BC1, BC2, Per1, Per2, Bcloc, BC3 string) {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()
	if err != nil {
		fmt.Println("fatal error config file: default \n", err)
		os.Exit(1)
	}
	// Declare var
	clientid := viper.GetString("DEV.A8")
	clientsecret := viper.GetString("DEV.A9")
	scopes := viper.GetString("DEV.A10")
	tokenurl := viper.GetString("DEV.A11")
	bc1 := viper.GetString("DEV.BC1")
	bc2 := viper.GetString("DEV.BC2")
	per1 := viper.GetString("DEV.C1")
	per2 := viper.GetString("DEV.BC4")
	bcloc := viper.GetString("DEV.C2")
	bc3 := viper.GetString("DEV.BC3")

	return clientid, clientsecret, scopes, tokenurl, bc1, bc2, per1, per2, bcloc, bc3
}

func Dril() (Budgetloc1, Budgetstruct string) {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()
	if err != nil {
		fmt.Println("fatal error config file: default \n", err)
		os.Exit(1)
	}
	budgetloc1 := viper.GetString("DEV.C3")
	budgetstruct := viper.GetString("DEV.C4")

	return budgetloc1, budgetstruct
}

func Distinct() (Budget string) {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()
	if err != nil {
		fmt.Println("fatal error config file: default \n", err)
		os.Exit(1)
	}
	distdata := viper.GetString("DEV.C5")

	return distdata
}

func Otp() (a, b, c, d, e, f string) {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()
	if err != nil {
		fmt.Println("fatal error config file: default \n", err)
		os.Exit(1)
	}
	A := viper.GetString("DEV.A")
	B := viper.GetString("DEV.B")
	C := viper.GetString("DEV.C")
	D := viper.GetString("DEV.D")
	E := viper.GetString("DEV.E")
	F := viper.GetString("DEV.F")

	return A, B, C, D, E, F
}

func Kafka() (a, b, c, d, e, f string) {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()
	if err != nil {
		fmt.Println("fatal error config file: default \n", err)
		os.Exit(1)
	}
	A := viper.GetString("DEV.A")
	B := viper.GetString("DEV.B")
	C := viper.GetString("DEV.C")
	D := viper.GetString("DEV.D")
	E := viper.GetString("DEV.E")
	F := viper.GetString("DEV.F")

	return A, B, C, D, E, F
}

func Hubspot() (a, b string) {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()
	if err != nil {
		fmt.Println("fatal error config file: default \n", err)
		os.Exit(1)
	}
	G := viper.GetString("DEV.G")
	H := viper.GetString("DEV.H")
	return G, H
}


func Qlikcloud() (a, b, c, d, e, f string) {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()
	if err != nil {
		fmt.Println("fatal error config file: default \n", err)
		os.Exit(1)
	}
	I := viper.GetString("DEV.I")
	J := viper.GetString("DEV.J")
	K := viper.GetString("DEV.K")
	L := viper.GetString("DEV.L")
	M := viper.GetString("DEV.M")
	N := viper.GetString("DEV.N")

	return I, J, K, L, M, N
}

func DVSA() (a, b, c, d, e, f, g, h string) {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()
	if err != nil {
		fmt.Println("fatal error config file: default \n", err)
		os.Exit(1)
	}
	O := viper.GetString("DEV.O")
	P := viper.GetString("DEV.P")
	Q := viper.GetString("DEV.Q")
	R := viper.GetString("DEV.R")
	S := viper.GetString("DEV.S")
	T := viper.GetString("DEV.T")
	U := viper.GetString("DEV.U")
	V := viper.GetString("DEV.V")
	return O, P, Q, R, S, T, U, V
}

func Qlikembed() (a, b, c, d string) {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()
	if err != nil {
		fmt.Println("fatal error config file: default \n", err)
		os.Exit(1)
	}
	W := viper.GetString("AUTO_DEV.W")
	X := viper.GetString("AUTO_DEV.X")
	Y := viper.GetString("AUTO_DEV.Y")
	Z := viper.GetString("AUTO_DEV.Z")
	return W, X, Y, Z
}


func Insightspace1() (spacesEndpoint, spaceName , accessKey , secretKey , region, uploads string) {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")

	err := viper.ReadInConfig()
	if err != nil {
		fmt.Println("fatal error config file: default \n", err)
		os.Exit(1)
	}
	aa := viper.GetString("AUTO_DEV.AA")
	ab := viper.GetString("AUTO_DEV.AB")
	ac := viper.GetString("AUTO_DEV.AC")
	ad := viper.GetString("AUTO_DEV.AD")
	ae := viper.GetString("AUTO_DEV.AE")
	af := viper.GetString("AUTO_DEV.AF")

	return aa, ab, ac, ad, ae, af
}

func Insightspace2() (spacesEndpoint, spaceName , accessKey , secretKey , region, uploads string) {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")

	err := viper.ReadInConfig()
	if err != nil {
		fmt.Println("fatal error config file: default \n", err)
		os.Exit(1)
	}
	aa := viper.GetString("AUTO_PROD.AA")
	ab := viper.GetString("AUTO_PROD.AB")
	ac := viper.GetString("AUTO_PROD.AC")
	ad := viper.GetString("AUTO_PROD.AD")
	ae := viper.GetString("AUTO_PROD.AE")
	af := viper.GetString("AUTO_PROD.AF")

	return aa, ab, ac, ad, ae, af
}

// const (
// 	ClientID     = "9444ceab-d3e7-4055-a58b-646c7bb97b11"
// 	ClientSecret = "YgR8Q~f4u9vJxFju4TNTt1y43pmlBMVm4nhgqcqG"
// 	TenantID     = "88c82b29-b553-4794-980c-010e970904c7"
// 	FromEmail    = "Noreply@insightdelivered.com"
// 	TokenURL     = "https://login.microsoftonline.com/" + TenantID + "/oauth2/v2.0/token"
// 	GraphSendURL = "https://graph.microsoft.com/v1.0/users/" + FromEmail + "/sendMail"
// )

// // ----------------------------------------------------------
// // STEP 1: Get Access Token from Azure AD
// // ----------------------------------------------------------
// func getGraphAccessToken() (string, error) {
// 	data := "client_id=" + ClientID +
// 		"&scope=https%3A%2F%2Fgraph.microsoft.com%2F.default" +
// 		"&client_secret=" + ClientSecret +
// 		"&grant_type=client_credentials"

// 	req, _ := http.NewRequest("POST", TokenURL, bytes.NewBufferString(data))
// 	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

// 	client := &http.Client{}
// 	resp, err := client.Do(req)
// 	if err != nil {
// 		return "", err
// 	}
// 	defer resp.Body.Close()

// 	body, _ := io.ReadAll(resp.Body)

// 	if resp.StatusCode != 200 {
// 		return "", fmt.Errorf("token error %d: %s", resp.StatusCode, string(body))
// 	}

// 	var res struct {
// 		AccessToken string `json:"access_token"`
// 	}
// 	json.Unmarshal(body, &res)

// 	return res.AccessToken, nil
// }

// // ----------------------------------------------------------
// // STEP 2: Send Email with Attachment (PDF Invoice etc.)
// // ----------------------------------------------------------
// func SendInvoiceEmail(to string, subject string, message string, fileName string, fileBytes []byte) error {

// 	// 1. Get Microsoft Graph Access Token
// 	token, err := getGraphAccessToken()
// 	if err != nil {
// 		return err
// 	}

// 	// 2. Encode attachment to Base64
// 	encoded := base64.StdEncoding.EncodeToString(fileBytes)

// 	emailData := map[string]interface{}{
// 		"message": map[string]interface{}{
// 			"subject": subject,
// 			"body": map[string]string{
// 				"contentType": "HTML",
// 				"content":     message,
// 			},
// 			"toRecipients": []map[string]map[string]string{
// 				{"emailAddress": {"address": to}},
// 			},
// 			"attachments": []map[string]string{
// 				{
// 					"@odata.type":   "#microsoft.graph.fileAttachment",
// 					"name":          fileName,
// 					"contentBytes":  encoded,
// 					"contentType":   "application/pdf",
// 				},
// 			},
// 		},
// 		"saveToSentItems": "true",
// 	}

// 	jsonBody, _ := json.Marshal(emailData)

// 	req, _ := http.NewRequest("POST", GraphSendURL, bytes.NewBuffer(jsonBody))
// 	req.Header.Add("Authorization", "Bearer "+token)
// 	req.Header.Add("Content-Type", "application/json")

// 	client := &http.Client{}
// 	resp, err := client.Do(req)
// 	if err != nil {
// 		return err
// 	}
// 	defer resp.Body.Close()

// 	respBody, _ := io.ReadAll(resp.Body)

// 	if resp.StatusCode >= 300 {
// 		return fmt.Errorf("graph sendMail error %d: %s", resp.StatusCode, string(respBody))
// 	}

// 	return nil
// }
// LoadR2Config loads R2 (Cloudflare) credentials from config.yaml and sets environment variables.
// It checks AUTO_DEV, AUTO_PROD, DEV, LIVE sections in that order for R2_ACCESS_KEY_ID, etc.
func LoadR2Config() {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")

	if err := viper.ReadInConfig(); err != nil {
		// config file not found or error reading, skip silently
		fmt.Printf("Warning: Failed to read config.yaml for R2 config: %v\n", err)
		return
	}

	setIfNotEmpty := func(key, envName string) {
		if viper.IsSet(key) {
			val := viper.GetString(key)
			if val != "" && os.Getenv(envName) == "" {
				os.Setenv(envName, val)
			}
		}
	}

	// Try to load from environment-specific sections in order: AUTO_DEV, AUTO_PROD, DEV, LIVE
	envSections := []string{"AUTO_DEV", "AUTO_PROD", "DEV", "LIVE"}
	found := false
	for _, section := range envSections {
		if viper.IsSet(section + ".R2_ACCESS_KEY_ID") {
			setIfNotEmpty(section+".R2_ACCESS_KEY_ID", "R2_ACCESS_KEY_ID")
			setIfNotEmpty(section+".R2_SECRET_ACCESS_KEY", "R2_SECRET_ACCESS_KEY")
			setIfNotEmpty(section+".R2_ENDPOINT", "R2_ENDPOINT")
			setIfNotEmpty(section+".R2_BUCKET", "R2_BUCKET")
			setIfNotEmpty(section+".R2_REGION", "R2_REGION")
			found = true
			break
		}
	}

	if found {
		fmt.Println("✓ R2 configuration loaded from config.yaml")
	}
}