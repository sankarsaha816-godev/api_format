package controllers

import (
	"bytes"
	"crypto/tls"
	"encoding/base64"
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
)

const enquiryTemplate = `<?xml version="1.0" encoding="UTF-8"?>
<soapenv:Envelope xmlns:soapenv="http://schemas.xmlsoap.org/soap/envelope/"
  xmlns:sup="http://webservices.hpi.co.uk/SupplementaryEnquiry%s">
  <soapenv:Header/>
  <soapenv:Body>
    <sup:EnquiryRequest>
      <sup:Authentication>
        <sup:SubscriberDetails>
          <sup:CustomerCode>%s</sup:CustomerCode>
          <sup:Initials>%s</sup:Initials>
          <sup:Password>%s</sup:Password>
        </sup:SubscriberDetails>
      </sup:Authentication>
      <sup:Request>
        <sup:Asset>
          <sup:Vrm>%s</sup:Vrm>
          <sup:Vin>%s</sup:Vin>
          <sup:Mileage>%s</sup:Mileage>
          <sup:Reference>%s</sup:Reference>
        </sup:Asset>
        <sup:PrimaryProduct>
          <sup:Code>%s</sup:Code>
        </sup:PrimaryProduct>
        %s
      </sup:Request>
    </sup:EnquiryRequest>
  </soapenv:Body>
</soapenv:Envelope>`

const parameterizedEnquiryTemplate = `<?xml version="1.0" encoding="UTF-8"?>
<soapenv:Envelope xmlns:soapenv="http://schemas.xmlsoap.org/soap/envelope/"
  xmlns:sup="http://webservices.hpi.co.uk/SupplementaryEnquiry%s">
  <soapenv:Header/>
  <soapenv:Body>
    <sup:ParameterizedEnquiryRequest>
      <sup:Authentication>
        <sup:SubscriberDetails>
          <sup:CustomerCode>%s</sup:CustomerCode>
          <sup:Initials>%s</sup:Initials>
          <sup:Password>%s</sup:Password>
        </sup:SubscriberDetails>
      </sup:Authentication>
      <sup:ParameterizedRequest>
        <sup:Asset>
          <sup:Vrm>%s</sup:Vrm>
          <sup:Vin>%s</sup:Vin>
          <sup:Mileage>%s</sup:Mileage>
          <sup:Reference>%s</sup:Reference>
        </sup:Asset>
        <sup:ParameterList>
          %s
        </sup:ParameterList>
        <sup:PrimaryProduct>
          <sup:Code>%s</sup:Code>
        </sup:PrimaryProduct>
        %s
      </sup:ParameterizedRequest>
    </sup:ParameterizedEnquiryRequest>
  </soapenv:Body>
</soapenv:Envelope>`

const controlledEnquiryTemplate = `<?xml version="1.0" encoding="UTF-8"?>
<soapenv:Envelope xmlns:soapenv="http://schemas.xmlsoap.org/soap/envelope/"
  xmlns:sup="http://webservices.hpi.co.uk/SupplementaryEnquiry%s">
  <soapenv:Header/>
  <soapenv:Body>
    <sup:ControlledEnquiryRequest>
      <sup:Authentication>
        <sup:SubscriberDetails>
          <sup:CustomerCode>%s</sup:CustomerCode>
          <sup:Initials>%s</sup:Initials>
          <sup:Password>%s</sup:Password>
        </sup:SubscriberDetails>
      </sup:Authentication>
      <sup:ControlledRequest>
        <sup:Asset>
          <sup:Vrm>%s</sup:Vrm>
          <sup:Vin>%s</sup:Vin>
          <sup:Mileage>%s</sup:Mileage>
          <sup:Reference>%s</sup:Reference>
        </sup:Asset>
        <sup:ParameterList>
          %s
        </sup:ParameterList>
        <sup:PrimaryProduct>
          <sup:Code>%s</sup:Code>
        </sup:PrimaryProduct>
        %s
      </sup:ControlledRequest>
    </sup:ControlledEnquiryRequest>
  </soapenv:Body>
</soapenv:Envelope>`

const documentEnquiryTemplate = `<?xml version="1.0" encoding="UTF-8"?>
<soapenv:Envelope xmlns:soapenv="http://schemas.xmlsoap.org/soap/envelope/"
  xmlns:sup="http://webservices.hpi.co.uk/SupplementaryEnquiry%s">
  <soapenv:Header/>
  <soapenv:Body>
    <sup:DocumentEnquiryRequest>
      %s
    </sup:DocumentEnquiryRequest>
  </soapenv:Body>
</soapenv:Envelope>`

type Credentials struct {
	CustomerCode string `json:"customerCode"`
	Initials     string `json:"initials"`
	Password     string `json:"password"`
}

type EnquiryRequest struct {
	Environment           string            `json:"environment"`
	Version               string            `json:"version"`
	EnquiryType           string            `json:"enquiryType"`
	Credentials           Credentials       `json:"credentials"`
	Vrm                   string            `json:"vrm"`
	Vin                   string            `json:"vin"`
	Mileage               string            `json:"mileage"`
	Reference             string            `json:"reference"`
	PrimaryProduct        string            `json:"primaryProduct"`
	SupplementaryProducts []string          `json:"supplementaryProducts"`
	Parameters            map[string]string `json:"parameters"`
	DocumentParams        map[string]string `json:"documentParams"`
}

type EnquiryResponse struct {
	Success bool `json:"success"`
	// RawRequest    string                 `json:"rawRequest"`
	// RawResponse   string                 `json:"rawResponse"`
	Endpoint      string                 `json:"endpoint"`
	SOAPAction    string                 `json:"soapAction"`
	StatusCode    int                    `json:"statusCode"`
	Duration      string                 `json:"duration"`
	Error         string                 `json:"error,omitempty"`
	ParsedData    map[string]interface{} `json:"parsedData,omitempty"`
	DashboardData map[string]interface{} `json:"dashboardData,omitempty"`
	Summary       map[string]interface{} `json:"summary,omitempty"`
}

type legacyHpiCheckRequest struct {
	RegNo       string `json:"regno"`
	Initials    string `json:"initials"`
	Mileage     int    `json:"mileage,omitempty"`
	Reference   string `json:"reference,omitempty"`
	VIN         string `json:"vin,omitempty"`
	Product     string `json:"product,omitempty"`
	Customer    string `json:"customerCode,omitempty"`
	Password    string `json:"password,omitempty"`
	Version     string `json:"version,omitempty"`
	Environment string `json:"environment,omitempty"`
}

func getEndpoint(env, version string) string {
	base := "https://wss.hpi.co.uk"
	if strings.EqualFold(env, "test") {
		base = "https://pat-wss.hpi.co.uk"
	}
	return fmt.Sprintf("%s/TradeSoap/services/SupplementaryEnquiry%s/", base, version)
}

func getSOAPAction(version, enquiryType string) string {
	actionMap := map[string]string{
		"standard":      "enquire",
		"parameterized": "parameterizedEnquiry",
		"controlled":    "controlledEnquiry",
		"document":      "documentEnquiry",
		"version":       "getServiceVersion",
	}
	action := actionMap[strings.ToLower(strings.TrimSpace(enquiryType))]
	if action == "" {
		action = "enquire"
	}
	return fmt.Sprintf("http://webservices.hpi.co.uk/SupplementaryEnquiry%s/%s", version, action)
}

func buildSupplementaryXML(products []string) string {
	var sb strings.Builder
	for _, code := range products {
		if strings.TrimSpace(code) != "" {
			sb.WriteString(fmt.Sprintf("        <sup:SupplementaryProduct>\n          <sup:Code>%s</sup:Code>\n        </sup:SupplementaryProduct>\n", code))
		}
	}
	return sb.String()
}

func buildParametersXML(params map[string]string) string {
	var sb strings.Builder
	for name, value := range params {
		sb.WriteString(fmt.Sprintf("          <sup:Parameter>\n            <sup:Name>%s</sup:Name>\n            <sup:Value>%s</sup:Value>\n          </sup:Parameter>\n", name, value))
	}
	return sb.String()
}

func buildDocumentParamsXML(params map[string]string) string {
	var sb strings.Builder
	for name, value := range params {
		sb.WriteString(fmt.Sprintf("      <sup:Parameter>\n        <sup:Name>%s</sup:Name>\n        <sup:Value>%s</sup:Value>\n      </sup:Parameter>\n", name, value))
	}
	return sb.String()
}

func buildSOAPEnvelope(req EnquiryRequest) string {
	version := req.Version
	suppXML := buildSupplementaryXML(req.SupplementaryProducts)
	vrm := strings.TrimSpace(strings.ReplaceAll(req.Vrm, " ", ""))

	switch strings.ToLower(strings.TrimSpace(req.EnquiryType)) {
	case "parameterized":
		paramsXML := buildParametersXML(req.Parameters)
		return fmt.Sprintf(parameterizedEnquiryTemplate,
			version,
			req.Credentials.CustomerCode,
			req.Credentials.Initials,
			req.Credentials.Password,
			vrm,
			req.Vin,
			req.Mileage,
			req.Reference,
			paramsXML,
			req.PrimaryProduct,
			suppXML,
		)
	case "controlled":
		paramsXML := buildParametersXML(req.Parameters)
		return fmt.Sprintf(controlledEnquiryTemplate,
			version,
			req.Credentials.CustomerCode,
			req.Credentials.Initials,
			req.Credentials.Password,
			vrm,
			req.Vin,
			req.Mileage,
			req.Reference,
			paramsXML,
			req.PrimaryProduct,
			suppXML,
		)
	case "document":
		docParamsXML := buildDocumentParamsXML(req.DocumentParams)
		return fmt.Sprintf(documentEnquiryTemplate, version, docParamsXML)
	default:
		return fmt.Sprintf(enquiryTemplate,
			version,
			req.Credentials.CustomerCode,
			req.Credentials.Initials,
			req.Credentials.Password,
			vrm,
			req.Vin,
			req.Mileage,
			req.Reference,
			req.PrimaryProduct,
			suppXML,
		)
	}
}

func sendSOAPRequest(endpoint, soapAction, soapBody string) ([]byte, int, time.Duration, error) {
	client := &http.Client{
		Timeout: 30 * time.Second,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: false,
				MinVersion:         tls.VersionTLS12,
			},
		},
	}

	httpReq, err := http.NewRequest(http.MethodPost, endpoint, bytes.NewBufferString(soapBody))
	if err != nil {
		return nil, 0, 0, fmt.Errorf("failed to create request: %w", err)
	}

	httpReq.Header.Set("Content-Type", "text/xml; charset=utf-8")
	httpReq.Header.Set("SOAPAction", soapAction)

	start := time.Now()
	resp, err := client.Do(httpReq)
	duration := time.Since(start)
	if err != nil {
		return nil, 0, duration, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, resp.StatusCode, duration, fmt.Errorf("failed to read response: %w", err)
	}

	return body, resp.StatusCode, duration, nil
}

func parseXMLToStructured(data []byte) map[string]interface{} {
	decoder := xml.NewDecoder(bytes.NewReader(data))
	root := make(map[string]interface{})
	stack := []map[string]interface{}{root}
	nameStack := []string{}

	for {
		token, err := decoder.Token()
		if err != nil {
			break
		}

		switch t := token.(type) {
		case xml.StartElement:
			localName := t.Name.Local
			nameStack = append(nameStack, localName)
			child := make(map[string]interface{})

			for _, attr := range t.Attr {
				aName := attr.Name.Local
				if aName == "nil" && attr.Value == "1" {
					child["_nil"] = true
				} else if aName == "type" {
					child["_type"] = attr.Value
				} else if aName != "" {
					child["@"+aName] = attr.Value
				}
			}

			parent := stack[len(stack)-1]
			if existing, ok := parent[localName]; ok {
				switch v := existing.(type) {
				case []interface{}:
					v = append(v, child)
					parent[localName] = v
				default:
					parent[localName] = []interface{}{v, child}
				}
			} else {
				parent[localName] = child
			}

			stack = append(stack, child)

		case xml.CharData:
			text := strings.TrimSpace(string(t))
			if text != "" && len(stack) > 1 {
				current := stack[len(stack)-1]
				current["_text"] = text
			}

		case xml.EndElement:
			if len(stack) > 1 {
				current := stack[len(stack)-1]

				if len(current) == 1 {
					if text, ok := current["_text"]; ok {
						parent := stack[len(stack)-2]
						name := nameStack[len(nameStack)-1]
						if arr, isArr := parent[name].([]interface{}); isArr {
							arr[len(arr)-1] = text
							parent[name] = arr
						} else {
							parent[name] = text
						}
					}
				}

				if _, isNil := current["_nil"]; isNil && len(current) <= 2 {
					parent := stack[len(stack)-2]
					name := nameStack[len(nameStack)-1]
					if arr, isArr := parent[name].([]interface{}); isArr {
						arr[len(arr)-1] = nil
						parent[name] = arr
					} else {
						parent[name] = nil
					}
				}

				stack = stack[:len(stack)-1]
			}
			if len(nameStack) > 0 {
				nameStack = nameStack[:len(nameStack)-1]
			}
		}
	}

	return root
}

func extractDashboardData(parsed map[string]interface{}) map[string]interface{} {
	dashboard := make(map[string]interface{})
	envelope := getMap(parsed, "Envelope")
	body := getMap(envelope, "Body")

	var results map[string]interface{}
	for _, key := range []string{"EnquiryResponse", "ParameterizedEnquiryResponse", "ControlledEnquiryResponse"} {
		resp := getMap(body, key)
		if resp != nil {
			results = getMap(resp, "RequestResults")
			break
		}
	}

	fault := getMap(body, "Fault")
	if fault != nil {
		dashboard["_fault"] = true
		dashboard["faultCode"] = getString(fault, "faultcode")
		dashboard["faultString"] = getString(fault, "faultstring")
		detail := getMap(fault, "detail")
		if detail != nil {
			hpiFault := getMap(detail, "HpiSoapFault")
			if hpiFault != nil {
				dashboard["faultErrors"] = hpiFault["Error"]
			}
		}
		return dashboard
	}

	if results == nil {
		resp := getMap(body, "DocumentEnquiryResponse")
		if resp != nil {
			dashboard["documentResponse"] = resp
		}
		return dashboard
	}

	asset := getMap(results, "Asset")
	if asset == nil {
		return dashboard
	}

	assetId := getMap(asset, "AssetIdentification")
	if assetId != nil {
		dashboard["identification"] = map[string]interface{}{
			"vrm":       getString(assetId, "Vrm"),
			"vin":       getString(assetId, "Vin"),
			"mileage":   getString(assetId, "Mileage"),
			"reference": getString(assetId, "Reference"),
		}
	}

	primary := getMap(asset, "PrimaryAssetData")
	if primary != nil {
		dvla := getMap(primary, "DVLA")
		if dvla != nil {
			dashboard["dvla"] = extractDVLA(dvla)
		}
		smmt := getMap(primary, "SMMT")
		if smmt != nil {
			dashboard["smmt"] = extractSMMT(smmt)
		}
		cr := primary["CrossReference"]
		if crStr, ok := cr.(string); ok {
			dashboard["crossReference"] = crStr
		}
		addInfo := getMap(primary, "AdditionalInformation")
		if addInfo != nil {
			dashboard["additionalInfo"] = map[string]interface{}{
				"co2Rating":              getString(addInfo, "CO2Rating"),
				"plateTransferIndicator": getString(addInfo, "PlateTransferIndicator"),
			}
		}
		fullCheck := getMap(primary, "FullCheck")
		if fullCheck != nil {
			dashboard["fullCheck"] = extractFullCheck(fullCheck)
		}
		flagsCheck := getMap(primary, "FlagsCheck")
		if flagsCheck != nil {
			dashboard["flagsCheck"] = extractFlags(flagsCheck)
		}
		instep := getMap(primary, "Instep")
		if instep != nil {
			dashboard["instep"] = instep
		}
		tp := primary["TranslatePlus"]
		if tp != nil {
			dashboard["translatePlus"] = tp
		}
		td := getMap(primary, "TecDocData")
		if td != nil {
			dashboard["tecDoc"] = td
		}
		cw := getMap(primary, "CrushWatchResults")
		if cw != nil {
			dashboard["crushWatch"] = cw
		}
	}

	supp := getMap(asset, "SupplementaryAssetData")
	if supp != nil {
		dashboard["supplementary"] = extractSupplementary(supp)
	}

	if warnings := results["Warning"]; warnings != nil {
		dashboard["warnings"] = warnings
	}

	return dashboard
}

func extractDVLA(dvla map[string]interface{}) map[string]interface{} {
	result := make(map[string]interface{})

	makeMap := getMap(dvla, "Make")
	if makeMap != nil {
		result["makeCode"] = getString(makeMap, "Code")
		result["makeDescription"] = getString(makeMap, "Description")
	}
	model := getMap(dvla, "Model")
	if model != nil {
		result["modelCode"] = getString(model, "Code")
		result["modelDescription"] = getString(model, "Description")
	}
	body := getMap(dvla, "Body")
	if body != nil {
		result["bodyCode"] = getString(body, "Code")
		result["bodyDescription"] = getString(body, "Description")
		result["doors"] = getString(body, "Doors")
		result["weight"] = getString(body, "Weight")
		colour := getMap(body, "Colour")
		if colour != nil {
			result["previousColours"] = getString(colour, "NumberPreviousColours")
			current := getMap(colour, "Current")
			if current != nil {
				result["currentColourCode"] = getString(current, "Code")
				result["currentColour"] = getString(current, "Description")
			}
			prev := getMap(colour, "Previous")
			if prev != nil {
				result["previousColourCode"] = getString(prev, "Code")
				result["previousColour"] = getString(prev, "Description")
				result["previousColourDate"] = getString(prev, "DateChanged")
			}
			orig := getMap(colour, "Original")
			if orig != nil {
				result["originalColourCode"] = getString(orig, "Code")
				result["originalColour"] = getString(orig, "Description")
			}
		}
		wp := getMap(body, "WheelPlan")
		if wp != nil {
			result["wheelPlan"] = getString(wp, "Description")
		}
	}
	engine := getMap(dvla, "Engine")
	if engine != nil {
		result["engineSize"] = getString(engine, "Size")
		result["engineNumber"] = getString(engine, "Number")
		fuel := getMap(engine, "Fuel")
		if fuel != nil {
			result["fuelType"] = getString(fuel, "Description")
			result["fuelCode"] = getString(fuel, "Code")
		}
	}
	keepers := getMap(dvla, "Keepers")
	if keepers != nil {
		result["lastKeeperChange"] = getString(keepers, "LastChangeOfKeeperDate")
		result["totalPreviousKeepers"] = getString(keepers, "TotalPreviousKeepers")
		pk := getMap(keepers, "PreviousKeeper")
		if pk != nil {
			result["prevKeeperAcquired"] = getString(pk, "DateAcquiredVehicle")
			result["prevKeeperDisposed"] = getString(pk, "DateDisposedOfVehicle")
		}
	}
	keyDates := getMap(dvla, "KeyDates")
	if keyDates != nil {
		fr := getMap(keyDates, "FirstRegistered")
		if fr != nil {
			result["firstRegistered"] = getString(fr, "Date")
			vi := getMap(fr, "VinVrmIndicator")
			if vi != nil {
				result["vinVrmIndicator"] = getString(vi, "Description")
			}
		}
		mfg := getMap(keyDates, "Manufactured")
		if mfg != nil {
			result["yearOfManufacture"] = getString(mfg, "Year")
		}
		result["scrapped"] = getString(keyDates, "Scrapped")
		result["exported"] = getString(keyDates, "Exported")
		result["isImported"] = getString(keyDates, "IsImported")
	}
	result["isFromNorthernIreland"] = getString(dvla, "IsFromNorthernIreland")

	return result
}

func extractSMMT(smmt map[string]interface{}) map[string]interface{} {
	result := make(map[string]interface{})
	makeMap := getMap(smmt, "Make")
	if makeMap != nil {
		result["makeCode"] = getString(makeMap, "Code")
		result["makeDescription"] = getString(makeMap, "Description")
	}
	model := getMap(smmt, "Model")
	if model != nil {
		result["modelCode"] = getString(model, "Code")
		result["modelDescription"] = getString(model, "Description")
	}
	body := getMap(smmt, "Body")
	if body != nil {
		result["bodyCode"] = getString(body, "Code")
		result["bodyDescription"] = getString(body, "Description")
		result["doors"] = getString(body, "Doors")
	}
	tx := getMap(smmt, "Transmission")
	if tx != nil {
		result["transmissionCode"] = getString(tx, "Code")
		result["transmissionDescription"] = getString(tx, "Description")
	}
	ms := getMap(smmt, "MarketSector")
	if ms != nil {
		result["marketSector"] = getString(ms, "Code")
	}
	return result
}

func extractFullCheck(fc map[string]interface{}) map[string]interface{} {
	result := make(map[string]interface{})
	pt := getMap(fc, "PlateTransfers")
	if pt != nil {
		result["plateTransferIndicator"] = getString(pt, "PlateTransferIndicator")
		result["plateTransferDate"] = getString(pt, "Date")
		result["plates"] = pt["Plate"]
	}
	sw := getMap(fc, "SecurityWatch")
	if sw != nil {
		result["securityWatch"] = map[string]interface{}{
			"totalVrmReadings": getString(sw, "TotalVRMReadings"),
			"totalVinReadings": getString(sw, "TotalVINReadings"),
			"telephone":        getString(sw, "Telephone"),
		}
	}
	fa := getMap(fc, "FinanceAgreements")
	if fa != nil {
		result["financeAgreements"] = fa["FinanceDetails"]
	}
	vcar := getMap(fc, "VCAR")
	if vcar != nil {
		result["vcar"] = map[string]interface{}{
			"damages":     vcar["Damages"],
			"thefts":      vcar["Thefts"],
			"inspections": vcar["Inspections"],
		}
	}
	si := getMap(fc, "StolenIncidents")
	if si != nil {
		result["stolenIncidents"] = si["StolenIncident"]
	}
	ih := getMap(fc, "IncidentHistory")
	if ih != nil {
		result["incidentHistory"] = ih
	}
	kh := getMap(fc, "KeeperHistory")
	if kh != nil {
		result["keeperHistory"] = kh["Keeper"]
	}
	return result
}

func extractFlags(fc map[string]interface{}) []map[string]interface{} {
	var flags []map[string]interface{}
	flagData := fc["Flag"]
	if flagData == nil {
		return flags
	}
	switch v := flagData.(type) {
	case []interface{}:
		for _, item := range v {
			if m, ok := item.(map[string]interface{}); ok {
				flags = append(flags, map[string]interface{}{
					"code":  getString(m, "Code"),
					"name":  getString(m, "Name"),
					"value": getString(m, "Value"),
				})
			}
		}
	case map[string]interface{}:
		flags = append(flags, map[string]interface{}{
			"code":  getString(v, "Code"),
			"name":  getString(v, "Name"),
			"value": getString(v, "Value"),
		})
	}
	return flags
}

func extractSupplementary(supp map[string]interface{}) map[string]interface{} {
	result := make(map[string]interface{})
	keyMap := map[string]string{
		"CAPBlackBookPlus":        "CAPBlackBook",
		"CAPBlackBookPlusData":    "CAPBlackBook",
		"BlackBookPlus":           "CAPBlackBook",
		"CAPValuationData":        "CAPValuation",
		"CAPLiveValuationData":    "CAPLiveValuation",
		"CAPLiveValuation":        "CAPLiveValuation",
		"MarketValuationData":     "MarketValuation",
		"AutotraderValuation":     "MarketValuation",
		"NMRCheckData":            "NMRCheck",
		"NationalMileageRegister": "NMRCheck",
		"NMRDetailsData":          "NMRDetails",
		"NmrDetailsList":          "NMRDetails",
		"DvlaMOTData":             "DvlaMOT",
		"DvlaMot":                 "DvlaMOT",
		"DvlaMOT":                 "DvlaMOT",
		"DVLAMOT":                 "DvlaMOT",
		"DVLAMot":                 "DvlaMOT",
		"FullDvlaMOTData":         "FullDvlaMOT",
		"FullDvlaMot":             "FullDvlaMOT",
		"FullDVLAMOT":             "FullDvlaMOT",
		"FullDVLAMot":             "FullDvlaMOT",
		"FullMOT":                 "FullDvlaMOT",
		"FullMot":                 "FullDvlaMOT",
		"MOTHistoryData":          "MOTHistory",
		"MOTHistoryList":          "MOTHistory",
		"TrackerData":             "Tracker",
		"CAPCodeData":             "CAP",
		"CAPMultiMatchData":       "CAPMultiMatch",
		"SpecCheckResult":         "SpecCheck",
		"SpecCheckResults":        "SpecCheck",
		"SpecCheckData":           "SpecCheck",
		"OptionListResponse":      "SpecCheck",
		"OptionListResult":        "SpecCheck",
		"OptionList":              "SpecCheck",
		"ExtraSMMTData":           "ExtraSMMT",
		"AnnotationsData":         "Annotations",
		"V5CheckResult":           "V5Check",
		"EnvironmentalSheet":      "EnvironmentalSheet",
	}

	for key, val := range supp {
		if key == "_nil" || key == "_type" || val == nil {
			continue
		}
		normalizedKey := key
		if mapped, ok := keyMap[key]; ok {
			normalizedKey = mapped
		}
		result[normalizedKey] = val
	}

	return result
}

func getMap(m map[string]interface{}, key string) map[string]interface{} {
	if m == nil {
		return nil
	}
	val, ok := m[key]
	if !ok || val == nil {
		return nil
	}
	result, ok := val.(map[string]interface{})
	if !ok {
		return nil
	}
	return result
}

func getString(m map[string]interface{}, key string) string {
	if m == nil {
		return ""
	}
	val, ok := m[key]
	if !ok || val == nil {
		return ""
	}
	switch v := val.(type) {
	case string:
		return v
	case map[string]interface{}:
		if t, ok := v["_text"]; ok {
			return fmt.Sprintf("%v", t)
		}
		return ""
	default:
		return fmt.Sprintf("%v", v)
	}
}

func xmlToMap(data []byte) map[string]interface{} {
	result := make(map[string]interface{})
	decoder := xml.NewDecoder(bytes.NewReader(data))
	var path []string
	for {
		token, err := decoder.Token()
		if err != nil {
			break
		}
		switch t := token.(type) {
		case xml.StartElement:
			path = append(path, t.Name.Local)
		case xml.CharData:
			text := strings.TrimSpace(string(t))
			if text != "" && len(path) > 0 {
				key := strings.Join(path, ".")
				result[key] = text
			}
		case xml.EndElement:
			if len(path) > 0 {
				path = path[:len(path)-1]
			}
		}
	}
	return result
}

func buildSummary(dashboard map[string]interface{}, req EnquiryRequest) map[string]interface{} {
	summary := map[string]interface{}{
		"enquiryType": strings.ToLower(strings.TrimSpace(req.EnquiryType)),
		"reference":   strings.TrimSpace(req.Reference),
		"primary":     strings.TrimSpace(req.PrimaryProduct),
	}

	if id, ok := dashboard["identification"].(map[string]interface{}); ok {
		if vrm := getString(id, "vrm"); vrm != "" {
			summary["vrm"] = vrm
		}
		if vin := getString(id, "vin"); vin != "" {
			summary["vin"] = vin
		}
		if mileage := getString(id, "mileage"); mileage != "" {
			summary["mileage"] = mileage
		}
		if reference := getString(id, "reference"); reference != "" {
			summary["reference"] = reference
		}
	}

	if dvla, ok := dashboard["dvla"].(map[string]interface{}); ok {
		makeDescription := getString(dvla, "makeDescription")
		modelDescription := getString(dvla, "modelDescription")
		if makeDescription != "" {
			summary["make"] = makeDescription
		}
		if modelDescription != "" {
			summary["model"] = modelDescription
		}
	}

	if fullCheck, ok := dashboard["fullCheck"].(map[string]interface{}); ok {
		if financeAgreements, ok := fullCheck["financeAgreements"]; ok {
			summary["financeCount"] = toItemCount(financeAgreements)
		}
	}

	if supplementary, ok := dashboard["supplementary"].(map[string]interface{}); ok {
		if nmrCheck, ok := supplementary["NMRCheck"].(map[string]interface{}); ok {
			if discrepancy := getString(nmrCheck, "Discrepancy"); discrepancy != "" {
				summary["nmrDiscrepancy"] = strings.EqualFold(discrepancy, "Y")
				summary["nmrDiscrepancyCode"] = discrepancy
			}
			if mileageStatus, ok := nmrCheck["MileageStatus"].(map[string]interface{}); ok {
				if code := getString(mileageStatus, "Code"); code != "" {
					summary["nmrStatusCode"] = code
				}
				if desc := getString(mileageStatus, "Description"); desc != "" {
					summary["nmrStatusDescription"] = desc
				}
			}
		}
	}

	if fault, ok := dashboard["_fault"].(bool); ok && fault {
		summary["fault"] = true
		if faultCode := getString(dashboard, "faultCode"); faultCode != "" {
			summary["faultCode"] = faultCode
		}
		if faultString := getString(dashboard, "faultString"); faultString != "" {
			summary["faultString"] = faultString
		}
	}

	return summary
}

func toItemCount(value interface{}) int {
	switch items := value.(type) {
	case []interface{}:
		return len(items)
	case nil:
		return 0
	default:
		return 1
	}
}

func formatXML(data []byte) string {
	var buf bytes.Buffer
	decoder := xml.NewDecoder(bytes.NewReader(data))
	encoder := xml.NewEncoder(&buf)
	encoder.Indent("", "  ")

	for {
		token, err := decoder.Token()
		if err != nil {
			break
		}
		_ = encoder.EncodeToken(token)
	}
	_ = encoder.Flush()
	return buf.String()
}

func sanitizeReference(value string) string {
	trimmed := strings.TrimSpace(value)
	if trimmed == "" {
		return ""
	}

	var builder strings.Builder
	for _, r := range strings.ToUpper(trimmed) {
		if (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9') {
			builder.WriteRune(r)
		}
	}

	cleaned := builder.String()
	if len(cleaned) > 15 {
		cleaned = cleaned[:15]
	}
	return cleaned
}

func hasFaultCode(dashboard map[string]interface{}, code string) bool {
	faultErrors, ok := dashboard["faultErrors"]
	if !ok || faultErrors == nil {
		return false
	}

	switch errors := faultErrors.(type) {
	case []interface{}:
		for _, item := range errors {
			if errMap, ok := item.(map[string]interface{}); ok {
				if strings.EqualFold(getString(errMap, "Code"), code) {
					return true
				}
			}
		}
	case []map[string]interface{}:
		for _, errMap := range errors {
			if strings.EqualFold(getString(errMap, "Code"), code) {
				return true
			}
		}
	case map[string]interface{}:
		return strings.EqualFold(getString(errors, "Code"), code)
	}

	return false
}

func runEnquiry(req EnquiryRequest) EnquiryResponse {
	endpoint := getEndpoint(req.Environment, req.Version)
	soapAction := getSOAPAction(req.Version, req.EnquiryType)
	soapBody := buildSOAPEnvelope(req)

	respBody, statusCode, duration, err := sendSOAPRequest(endpoint, soapAction, soapBody)
	response := EnquiryResponse{
		// RawRequest: soapBody,
		Endpoint:   endpoint,
		SOAPAction: soapAction,
		StatusCode: statusCode,
		Duration:   duration.String(),
	}

	if err != nil {
		response.Success = false
		response.Error = err.Error()
		return response
	}

	formattedResponse := formatXML(respBody)
	if formattedResponse == "" {
		formattedResponse = string(respBody)
	}
	// response.RawResponse = formattedResponse
	response.Success = statusCode == http.StatusOK
	response.ParsedData = xmlToMap(respBody)

	structured := parseXMLToStructured(respBody)
	response.DashboardData = extractDashboardData(structured)
	response.Summary = buildSummary(response.DashboardData, req)

	if statusCode == http.StatusInternalServerError {
		response.Success = false
		response.Error = "SOAP Fault returned (HTTP 500) - check raw response for details"
	}

	return response
}

func executeEnquiry(req EnquiryRequest) EnquiryResponse {
	response := runEnquiry(req)

	if response.Success || !strings.EqualFold(req.Environment, "test") {
		return response
	}

	if response.StatusCode != http.StatusInternalServerError {
		return response
	}

	hasInvalidReference := hasFaultCode(response.DashboardData, "2021")
	hasInvalidMileage := hasFaultCode(response.DashboardData, "014")
	forceStandardRetry := strings.EqualFold(req.EnquiryType, "standard")
	if !forceStandardRetry && !hasInvalidReference && !hasInvalidMileage {
		return response
	}

	timestampRef := "REF" + time.Now().Format("060102150405")
	sanitizedRef := sanitizeReference(req.Reference)
	if sanitizedRef == "" {
		sanitizedRef = timestampRef
	}

	retryCandidates := []struct {
		reference string
		mileage   string
	}{
		{reference: sanitizedRef, mileage: "0"},
		{reference: "POSTMAN001", mileage: "0"},
		{reference: timestampRef, mileage: "0"},
		{reference: "POSTMAN001", mileage: "1"},
	}

	var lastRetryResp EnquiryResponse
	var lastRetryRef, lastRetryMileage string
	for idx, candidate := range retryCandidates {
		retryReq := req
		retryReq.Reference = candidate.reference
		retryReq.Mileage = candidate.mileage

		retryResp := runEnquiry(retryReq)
		if retryResp.Summary == nil {
			retryResp.Summary = map[string]interface{}{}
		}
		retryResp.Summary["retryApplied"] = true
		retryResp.Summary["retryReason"] = "HPI rejected reference/mileage on first attempt"
		retryResp.Summary["retryAttempt"] = idx + 1
		retryResp.Summary["retryOriginalReference"] = req.Reference
		retryResp.Summary["retryOriginalMileage"] = req.Mileage
		retryResp.Summary["retryFinalReference"] = retryReq.Reference
		retryResp.Summary["retryFinalMileage"] = retryReq.Mileage

		if retryResp.Success {
			return retryResp
		}

		lastRetryResp = retryResp
		lastRetryRef = retryReq.Reference
		lastRetryMileage = retryReq.Mileage
	}

	if lastRetryResp.Summary != nil {
		lastRetryResp.Summary["retrySucceeded"] = false
	}

	if lastRetryResp.StatusCode != 0 {
		return lastRetryResp
	}

	if response.Summary == nil {
		response.Summary = map[string]interface{}{}
	}
	response.Summary["retryApplied"] = true
	response.Summary["retrySucceeded"] = false
	response.Summary["retryOriginalReference"] = req.Reference
	response.Summary["retryOriginalMileage"] = req.Mileage
	response.Summary["retryFinalReference"] = lastRetryRef
	response.Summary["retryFinalMileage"] = lastRetryMileage
	return response
}

func HPIEnquiry(c *fiber.Ctx) error {
	var req EnquiryRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(http.StatusBadRequest).JSON(EnquiryResponse{
			Success: false,
			Error:   "Invalid request body: " + err.Error(),
		})
	}

	if req.Credentials.CustomerCode == "" {
		req.Credentials.CustomerCode = os.Getenv("HPI_CUSTOMER_CODE")
	}
	if req.Credentials.Initials == "" {
		req.Credentials.Initials = os.Getenv("HPI_INITIALS")
	}
	if req.Credentials.Password == "" {
		req.Credentials.Password = os.Getenv("HPI_PASSWORD")
	}

	if req.Credentials.CustomerCode == "" || req.Credentials.Password == "" {
		return c.Status(http.StatusBadRequest).JSON(EnquiryResponse{
			Success: false,
			Error:   "Customer code and password are required",
		})
	}
	if req.Vrm == "" && req.Vin == "" && !strings.EqualFold(req.EnquiryType, "document") {
		return c.Status(http.StatusBadRequest).JSON(EnquiryResponse{
			Success: false,
			Error:   "Either VRM or VIN must be provided",
		})
	}

	if req.Version == "" {
		req.Version = "V1"
	}
	if req.Environment == "" {
		req.Environment = "test"
	}
	if req.EnquiryType == "" {
		req.EnquiryType = "standard"
	}
	if req.Mileage == "" {
		req.Mileage = "0"
	}
	if req.PrimaryProduct == "" && !strings.EqualFold(req.EnquiryType, "document") {
		req.PrimaryProduct = "HPI64"
	}

	response := executeEnquiry(req)
	if !response.Success && response.Error != "" && response.StatusCode == 0 {
		return c.Status(http.StatusBadGateway).JSON(response)
	}

	return c.Status(http.StatusOK).JSON(response)
}

func HPICheckVehicle(c *fiber.Ctx) error {
	var legacy legacyHpiCheckRequest
	if err := c.BodyParser(&legacy); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"code":    http.StatusBadRequest,
			"status":  false,
			"message": "invalid request body",
			"error":   err.Error(),
		})
	}

	vrm := strings.ToUpper(strings.ReplaceAll(strings.TrimSpace(legacy.RegNo), " ", ""))
	if vrm == "" {
		vrm = strings.ToUpper(strings.ReplaceAll(strings.TrimSpace(c.Query("regno")), " ", ""))
	}
	if vrm == "" && strings.TrimSpace(legacy.VIN) == "" {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"code":    http.StatusBadRequest,
			"status":  false,
			"message": "regno or vin is required",
		})
	}

	initials := strings.ToUpper(strings.TrimSpace(legacy.Initials))
	if initials == "" {
		initials = strings.ToUpper(strings.TrimSpace(c.Query("initials")))
	}
	if initials == "" {
		initials = os.Getenv("HPI_INITIALS")
	}

	mileage := "0"
	if legacy.Mileage > 0 {
		mileage = fmt.Sprintf("%d", legacy.Mileage)
	}

	request := EnquiryRequest{
		Environment:    firstNonEmpty(legacy.Environment, "test"),
		Version:        firstNonEmpty(legacy.Version, "V1"),
		EnquiryType:    "standard",
		Vrm:            vrm,
		Vin:            strings.TrimSpace(legacy.VIN),
		Mileage:        mileage,
		Reference:      strings.TrimSpace(legacy.Reference),
		PrimaryProduct: firstNonEmpty(strings.TrimSpace(legacy.Product), "HPI64"),
		Credentials: Credentials{
			CustomerCode: firstNonEmpty(strings.TrimSpace(legacy.Customer), os.Getenv("HPI_CUSTOMER_CODE")),
			Initials:     initials,
			Password:     firstNonEmpty(strings.TrimSpace(legacy.Password), os.Getenv("HPI_PASSWORD")),
		},
	}

	if request.Credentials.CustomerCode == "" || request.Credentials.Password == "" {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"code":    http.StatusBadRequest,
			"status":  false,
			"message": "missing HPI credentials",
			"hint":    "set HPI_CUSTOMER_CODE and HPI_PASSWORD env vars or pass in payload",
		})
	}

	response := executeEnquiry(request)
	ok := response.Success && response.StatusCode == http.StatusOK
	statusCode := http.StatusOK
	if !ok {
		statusCode = http.StatusBadGateway
	}

	return c.Status(statusCode).JSON(fiber.Map{
		"code":    statusCode,
		"status":  ok,
		"message": map[bool]string{true: "HPI check completed", false: "HPI check failed"}[ok],
		"details": response,
	})
}

func HPIHealth(c *fiber.Ctx) error {
	return c.Status(http.StatusOK).JSON(fiber.Map{
		"code":      http.StatusOK,
		"status":    true,
		"message":   "ok",
		"timestamp": time.Now().Format(time.RFC3339),
		"service":   "HPI Trade Webservices Proxy",
		"details": fiber.Map{
			"provider":         "hpi trade soap",
			"test_endpoint":    getEndpoint("test", "V1"),
			"live_endpoint":    getEndpoint("live", "V1"),
			"enquiry_endpoint": "/hpi/enquiry",
		},
	})
}

func HPILoginTest(c *fiber.Ctx) error {
	customer := os.Getenv("HPI_CUSTOMER_CODE")
	password := os.Getenv("HPI_PASSWORD")
	initials := os.Getenv("HPI_INITIALS")
	if customer == "" || password == "" {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"code":    http.StatusBadRequest,
			"status":  false,
			"message": "missing HPI credentials in environment",
		})
	}

	req := EnquiryRequest{
		Environment: "test",
		Version:     "V1",
		EnquiryType: "version",
		Credentials: Credentials{
			CustomerCode: customer,
			Initials:     initials,
			Password:     password,
		},
		Vrm:            "",
		Vin:            "",
		Mileage:        "0",
		Reference:      "login-test",
		PrimaryProduct: "HPI64",
	}

	resp := executeEnquiry(req)
	if resp.Error != "" && resp.StatusCode == 0 {
		return c.Status(http.StatusBadGateway).JSON(fiber.Map{
			"code":    http.StatusBadGateway,
			"status":  false,
			"message": "connectivity test failed",
			"error":   resp.Error,
		})
	}

	return c.Status(http.StatusOK).JSON(fiber.Map{
		"code":    http.StatusOK,
		"status":  true,
		"message": "connectivity test completed",
		"details": resp,
	})
}

func HPIProducts(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"primaryProducts": []fiber.Map{
			{"code": "HPI64", "name": "ID Check", "description": "Basic HPI ID check"},
			{"code": "HPI63", "name": "Full HPI Check", "description": "Full HPI check including finance, stolen, VCAR, plate transfers"},
			{"code": "HPI31", "name": "Translate", "description": "Translate-type check with Instep elements (BrokerNet)"},
			{"code": "HPI75", "name": "Translate Plus", "description": "TranslatePlus-style check with insurance/security details"},
			{"code": "HPF12", "name": "Flags Check", "description": "Flags Check (includes basic HPI ID check)"},
			{"code": "KT001", "name": "TecDoc Check", "description": "TecDoc Check (includes basic HPI ID check)"},
			{"code": "EA001", "name": "Crush Watch", "description": "Crush Watch check (codes EA001 to EA005)"},
			{"code": "HPIMM", "name": "HPI Multi-Match", "description": "Combination of HPI64, ADSMT and CAP multi-match"},
			{"code": "HPIT1", "name": "Translate Plus + Tracker", "description": "TranslatePlus with Tracker data"},
		},
		"supplementaryProducts": []fiber.Map{
			{"code": "CAPCD", "name": "CAP Code", "description": "CAP code lookup", "requiresPrimary": "any"},
			{"code": "HPI25", "name": "CAP Valuation", "description": "CAP vehicle valuation with mileage and condition", "requiresPrimary": "any"},
			{"code": "CAPBB", "name": "CAP Black Book Plus", "description": "Historic and future valuations (requires HPI25)", "requiresPrimary": "any"},
			{"code": "CAPL1", "name": "CAP Live Valuation", "description": "Extended CAP valuation with daily live values", "requiresPrimary": "any"},
			{"code": "ADSMT", "name": "Extra SMMT", "description": "Extended SMMT technical data", "requiresPrimary": "any"},
			{"code": "HP001", "name": "NMR Check", "description": "National Mileage Register check (requires HPI63)", "requiresPrimary": "HPI63"},
			{"code": "NMRDT", "name": "NMR Details", "description": "NMR detailed mileage history (requires HPI63 + HP001)", "requiresPrimary": "HPI63"},
			{"code": "HP003", "name": "NMR Investigation", "description": "NMR investigation request (requires HPI63 + HP001)", "requiresPrimary": "HPI63"},
			{"code": "PDF01", "name": "Trade Certificates", "description": "PDF certificate via email (requires HPI63)", "requiresPrimary": "HPI63"},
			{"code": "ENVSH", "name": "Environmental Sheet", "description": "Environmental data PDF", "requiresPrimary": "any"},
			{"code": "HPI58", "name": "V5 Check", "description": "V5 registration document verification", "requiresPrimary": "any"},
			{"code": "HPI70", "name": "Annotations", "description": "Previous HPI check annotations history", "requiresPrimary": "any"},
			{"code": "AVX01", "name": "Screen Check (OEM)", "description": "OEM codes only", "requiresPrimary": "any"},
			{"code": "AVX04", "name": "Screen Check (OEM+ARGIC)", "description": "OEM and ARGIC codes", "requiresPrimary": "any"},
			{"code": "AVX05", "name": "Screen Check (ARGIC)", "description": "ARGIC codes only", "requiresPrimary": "any"},
			{"code": "AVX02", "name": "Spec Check", "description": "Standard and optional equipment list", "requiresPrimary": "any"},
			{"code": "AVX03", "name": "Spec Check (Franchise)", "description": "Spec check excluding own-make vehicles", "requiresPrimary": "any"},
			{"code": "PDF06", "name": "Spec Check PDF", "description": "PDF of spec check data (requires AVX02/AVX03)", "requiresPrimary": "any"},
			{"code": "VALDP", "name": "Autotrader Market Values", "description": "Autotrader retail/trade/private/part-exchange values", "requiresPrimary": "any"},
			{"code": "HPIM1", "name": "MOT Enquiry", "description": "MOT expiry date", "requiresPrimary": "any"},
			{"code": "HPIMT", "name": "Full DVLA MOT", "description": "MOT expiry date and vehicle tax data", "requiresPrimary": "any"},
			{"code": "CAPMM", "name": "CAP Multi-Match", "description": "Multiple CAP code matches", "requiresPrimary": "any"},
			{"code": "HPITR", "name": "Tracker", "description": "Tracker fitted status (requires HPI63)", "requiresPrimary": "HPI63"},
			{"code": "PDF07", "name": "PDF Response", "description": "PDF of results in response (parameterized only)", "requiresPrimary": "HPI63"},
			{"code": "HPIMH", "name": "MOT History", "description": "Full MOT history with advisories and failures", "requiresPrimary": "any"},
		},
		"enquiryTypes": []fiber.Map{
			{"type": "standard", "name": "Standard Enquiry", "description": "Basic HPI check with optional supplementary products"},
			{"type": "parameterized", "name": "Parameterized Enquiry", "description": "For products requiring extra parameters (NMR, PDF, Spec Check)"},
			{"type": "controlled", "name": "Controlled Enquiry", "description": "Access-controlled enquiry based on account rules"},
			{"type": "document", "name": "Document Enquiry", "description": "Retrieve a previously generated PDF certificate"},
		},
		"parameterTemplates": fiber.Map{
			"nmr": []fiber.Map{
				{"name": "NMR_ODOMETER_READING", "description": "Current odometer reading", "required": true},
				{"name": "NMR_ODOMETER_UNITS", "description": "M=miles, K=kilometres", "required": true, "default": "M"},
				{"name": "NMR_ODOMETER_READING_DATE", "description": "Date of reading (dd/mm/yyyy)", "required": false},
				{"name": "NMR_ODOMETER_CHANGED", "description": "Has odometer been changed? Y/N", "required": true, "default": "N"},
				{"name": "NMR_ODOMETER_TOTAL", "description": "Total from all odometers (if changed)", "required": false},
				{"name": "NMR_ODOMETER_READING_1", "description": "Previous odometer reading (if changed)", "required": false},
				{"name": "NMR_ODOMETER_READING_1_DATE", "description": "Previous reading date (if changed)", "required": false},
			},
			"nmrInvestigation": []fiber.Map{
				{"name": "PDF_EMAIL_ADDRESS", "description": "Email for PDF certificate", "required": true},
				{"name": "NMR_INV_STOCK_REFERENCE", "description": "Dealer stock reference (max 15 chars)", "required": true},
			},
			"tradeCertificate": []fiber.Map{
				{"name": "PDF_EMAIL_ADDRESS", "description": "Email for PDF certificate", "required": true},
				{"name": "NMR_ODOMETER_READING", "description": "Odometer reading (nearest mile)", "required": true},
			},
			"environmentalSheet": []fiber.Map{
				{"name": "PDF_EMAIL_ADDRESS", "description": "Email for PDF (or use ENV_SHEET_PDF_IN_RESPONSE)", "required": false},
				{"name": "ENV_SHEET_PDF_IN_RESPONSE", "description": "Set TRUE to embed PDF in response", "required": false},
			},
			"specCheckPdf": []fiber.Map{
				{"name": "SPEC_CHECK_PDF_IN_RESPONSE", "description": "Set TRUE for PDF in response", "required": false, "default": "TRUE"},
				{"name": "PDF_EMAIL_ADDRESS", "description": "Email for PDF delivery", "required": false},
				{"name": "SHOW_STANDARD_SPEC_CHECK_OPTIONS", "description": "Show standard options", "required": false, "default": "false"},
			},
			"v5Check": []fiber.Map{
				{"name": "V5_SERIAL_NUMBER", "description": "V5 document serial number", "required": true},
				{"name": "V5_ISSUE_DATE", "description": "V5 issue date (dd/mm/yyyy)", "required": true},
			},
			"documentEnquiry": []fiber.Map{
				{"name": "product", "description": "Product code (e.g. HPI63)", "required": true},
				{"name": "customer", "description": "Customer code", "required": true},
				{"name": "password", "description": "Password", "required": true},
				{"name": "vrm", "description": "Vehicle registration mark", "required": true},
				{"name": "vin", "description": "Vehicle identification number (optional)", "required": false},
				{"name": "select", "description": "'most recent' or 'last allowed'", "required": false, "default": "most recent"},
			},
		},
		"v2Differences": fiber.Map{
			"HPI63": "Full keeper history + Incident History (IncidentRecord & SalvageRecord)",
			"HPF12": "Additional Flag code 95 for Incident History",
			"VALDP": "Removed deprecated fields, schema updated for current fields",
		},
	})
}

func HPIDecodePDF(c *fiber.Ctx) error {
	var body struct {
		Base64Data string `json:"base64Data"`
		Filename   string `json:"filename"`
	}
	if err := c.BodyParser(&body); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request"})
	}

	decoded, err := base64.StdEncoding.DecodeString(body.Base64Data)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Invalid base64 data"})
	}

	filename := body.Filename
	if filename == "" {
		filename = "document.pdf"
	}

	c.Set("Content-Type", "application/pdf")
	c.Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filename))
	return c.Send(decoded)
}
