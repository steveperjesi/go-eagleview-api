package main

import (
	// "errors"
	"fmt"
	// "strconv"
	"time"

	"github.com/steveperjesi/go-eagleview-api/eagleview"

	"githost.in/peartree/pkg/log"
)

func main() {
	products, err := GetAvailableProducts()
	if err != nil {
		panic(err)
	}

	log.Info().Interface("PRODUCTS", products).Msg("main")

	time.Sleep(3 * time.Second)

	orderResponse, err := PlaceOrder()
	if err != nil {
		panic(err)
	}

	log.Info().Interface("orderResponse", orderResponse).Msg("main")

	// initVars, err := eagleview.GetInitVars()
	// if err != nil {
	// 	panic(err)
	// }

	// // log.Info().Interface("initVars", initVars).Msg("main")

	// client, err := eagleview.NewClient(initVars[0], initVars[1], initVars[2], initVars[3])
	// if err != nil {
	// 	panic(err)
	// }

	// log.Info().Interface("client", client).Msg("main")

	// params := make(map[string]interface{})

	// // Get available products
	// endpoint := "/v2/Product/GetAvailableProducts"

	// response := []eagleview.Product{}
	// err = client.Request("GET", endpoint, params, &response)
	// if err != nil {
	// 	panic(err)
	// }

	// log.Debug().Interface("GetAvailableProducts", response).Msg("main")

	fmt.Println("* DONE *")
}

func PlaceOrder() (eagleview.PlaceOrderResponse, error) {
	initVars, err := eagleview.GetInitVars()
	if err != nil {
		return eagleview.PlaceOrderResponse{}, err
	}

	client, err := eagleview.NewClient(initVars[0], initVars[1], initVars[2], initVars[3])
	if err != nil {
		return eagleview.PlaceOrderResponse{}, err
	}

	testAddresses := []eagleview.Address{}

	addr1 := eagleview.Address{
		Address:     "1313 S Disneyland Dr",
		City:        "Anaheim",
		State:       "CA",
		Zip:         "92802",
		Country:     "US",
		AddressType: 1,
	}

	addr2 := eagleview.Address{
		Address:     "3911 S Figueroa St",
		City:        "Los Angeles",
		State:       "CA",
		Zip:         "90037",
		Country:     "US",
		Latitude:    34.014167,
		Longitude:   -118.287778,
		AddressType: 4,
	}

	testAddresses = append(testAddresses, addr1, addr2)

	attributes := []eagleview.ReportAttribute{}

	attr1 := eagleview.ReportAttribute{
		Attribute: 9,     // Standard Product Is OK To Substitute If Ordered Product Can Not Be Completed
		Value:     "Yes", // Yes, No, Ask
	}

	attr2 := eagleview.ReportAttribute{
		Attribute: 10,    // Premium Product Is OK To Substitute If Ordered Product Can Not Be Completed
		Value:     "Yes", // Yes, No, Ask
	}

	attr3 := eagleview.ReportAttribute{
		Attribute: 24,                        // Additional Emails To Send Report-Related Information To
		Value:     "steve@theprotectall.com", // Semi-colon Seperated List Allowed
	}

	attributes = append(attributes, attr1, attr2, attr3)

	reports := []eagleview.Report{}
	report := eagleview.Report{
		Addresses:  testAddresses,
		BuildingID: "Parking structure only",
		// PrimaryProductID:           63, // EagleView Inform Advanced for Commercial
		PrimaryProductID:           32, // EagleView Premium - Commercial
		DeliveryProductID:          8,  // Regular delivery; Default turnaround time of 3 business days. No additional fees
		AddOnProductIDs:            []int{33},
		MeasurementInstructionType: 4, // CommercialComplex
		ReportAttributes:           attributes,
		ClaimNumber:                "Test-Claim-#PJ8890-1",
		ClaimInfo:                  "Test claim",
		BatchID:                    "F12",
		CatID:                      "CAT5-E87",
		ChangesInLastFourYears:     false,
		PONumber:                   "PO# 9932",
		Comments:                   "Happy happy, joy joy",
		ReferenceID:                "REF123456", // Used with incoming webhook
		InsuredName:                "Luke Skywalker",
		PolicyNumber:               "009866324",
		// UpgradeFromReportID:        0,
	}

	reports = append(reports, report)

	// dateOfLoss, err := time.Parse(time.RFC3339, "2020-02-14T12:00:00Z")
	dateOfLoss, err := time.Parse(time.RFC3339, "2020-02-14T12:00:00Z")
	if err != nil {
		return eagleview.PlaceOrderResponse{}, err
	}

	reports[0].DateOfLoss = dateOfLoss

	// TODO: untested since eagleview does not have a test card to use
	// creditCard := eagleview.CreditCardData{
	// 	FirstName:       "Test",
	// 	LastName:        "Tester",
	// 	ExpirationMonth: 12,
	// 	ExpirationYear:  2024,
	// 	Number:          "0000222244445555",
	// 	Type:            2, // 2 = Visa
	// }

	params := map[string]interface{}{
		"OrderReports": reports,
		"PromoCode":    "",
		// "PlaceOrderUser": "",
		// "CreditCardData": creditCard,
	}

	// Place an order
	endpoint := "/v2/Order/PlaceOrder"

	response := eagleview.PlaceOrderResponse{}

	err = client.Request("POST", endpoint, "", params, &response)
	if err != nil {
		return eagleview.PlaceOrderResponse{}, err
	}

	return response, nil
}

// GetAvailableProducts works !!!
func GetAvailableProducts() ([]eagleview.Product, error) {
	initVars, err := eagleview.GetInitVars()
	if err != nil {
		return nil, err
	}

	client, err := eagleview.NewClient(initVars[0], initVars[1], initVars[2], initVars[3])
	if err != nil {
		return nil, err
	}

	params := make(map[string]interface{})

	// Get available products
	endpoint := "/v2/Product/GetAvailableProducts"

	response := []eagleview.Product{}
	err = client.Request("GET", endpoint, "application/x-www-form-urlencoded", params, &response)
	if err != nil {
		return nil, err
	}

	// log.Debug().Interface("GetAvailableProducts", response).Msg("main")
	return response, nil
}

// results, err := client.Search(yelp.SearchOptions{
// 	Term:      "restaurants",
// 	Latitude:  36.0813328,
// 	Longitude: -115.3161651,
// 	SortBy:    "distance",
// 	// Location: "las vegas",
// 	// Radius:    40000, // 40000 meters is the max allowed value
// 	// Limit:  10,
// 	// Categories: "localservices",
// 	// OpenNow:   true,
// 	// Offset:    10,
// 	// Price:  "1,2,3,4",  // 1 = $, 2 = $$, 3 = $$$, 4 = $$$$
// 	// Attributes: "hot_and_new",
// 	// Locale:     "en_US",
// 	// OpenAt:    1572494399,
// })
// if err != nil {
// 	panic(err)
// }

// fmt.Println("\nTOTAL: ", results.Total)
// fmt.Println("\nREGION: ", results.Region)

// for _, biz := range results.Businesses {
// 	fmt.Println("Name\t\t", biz.Name)
// 	fmt.Println("ID\t\t", biz.Id)
// 	fmt.Println("Alias\t\t", biz.Alias)
// 	fmt.Println("Rating\t\t", biz.Rating)
// 	fmt.Println("Price\t\t", biz.Price)
// 	fmt.Println("IsClosed\t", biz.IsClosed)
// 	fmt.Println("Url\t\t", biz.Url)
// 	fmt.Println("Distance\t", biz.Distance)
// 	fmt.Println("ReviewCount\t", biz.ReviewCount)
// 	fmt.Println("Latitude\t", biz.Coordinates.Latitude)
// 	fmt.Println("Longitude\t", biz.Coordinates.Longitude)

// 	fmt.Println("Phone\t\t", biz.Phone)

// 	fmt.Println("Address1\t", biz.Location.Address1)
// 	fmt.Println("Address2\t", biz.Location.Address2)
// 	fmt.Println("Address3\t", biz.Location.Address3)
// 	fmt.Println("City\t\t", biz.Location.City)
// 	fmt.Println("State\t\t", biz.Location.State)
// 	fmt.Println("ZipCode\t\t", biz.Location.ZipCode)
// 	fmt.Println("Country\t\t", biz.Location.Country)

// 	cats := yelp.CategoriesToString(biz.Categories)
// 	fmt.Println("Categories\t", cats)
// 	fmt.Println("-----")
// }
// }
