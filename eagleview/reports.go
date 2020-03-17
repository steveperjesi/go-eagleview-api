package eagleview

import (
	"time"
)

// type OrderReportsRequest struct {
// 	OrderReports   []OrderReport
// 	PromoCode      string
// 	PlaceOrderUser string
// 	CreditCardData CreditCardData
// }

// type OrderReportResponse struct {
// 	MinPrice                  float64
// 	MaxPrice                  float64
// 	Discounts                 []Discount
// 	PrimaryProductPriceQuotes ReportPriceQuote
// 	DeliveryProductPriceQuote ReportPriceQuote
// 	AddOnProductPriceQuotes   ReportPriceQuote
// }

type ReportPriceQuote struct {
	ProductID int
	MinPrice  float64
	MaxPrice  float64
}

type Discount struct {
	Description      string
	MinAmount        float64
	MaxAmount        float64
	RejectionMessage string
}

type Report struct {
	Addresses                  []Address         `json:"ReportAddresses"`
	BuildingID                 string            `json:"BuildingId"`
	PrimaryProductID           int               `json:"PrimaryProductId"`
	DeliveryProductID          int               `json:"DeliveryProductId"`
	AddOnProductIDs            []int             `json:"AddOnProductIds"`
	MeasurementInstructionType int               `json:"MeasurementInstructionType"`
	ReportAttributes           []ReportAttribute `json:"ReportAttibutes"`
	ClaimNumber                string            `json:"ClaimNumber"`
	ClaimInfo                  string            `json:"ClaimInfo"`
	BatchID                    string            `json:"BatchId"`
	CatID                      string            `json:"CatId"`
	ChangesInLastFourYears     bool              `json:"ChangesInLast4Years"`
	PONumber                   string            `json:"PONumber"`
	Comments                   string            `json:"Comments"`
	ReferenceID                string            `json:"ReferenceId"`
	InsuredName                string            `json:"InsuredName"`
	UpgradeFromReportID        int               `json:"UpgradeFromReportId"`
	PolicyNumber               string            `json:"PolicyNumber"`
	DateOfLoss                 time.Time         `json:"DateOfLoss"` // "2018-07-19T15:48:45.4836116-07:00"
}

type Address struct {
	Address     string  `json:"Address"`
	City        string  `json:"City"`
	State       string  `json:"State"`
	Zip         string  `json:"Zip"`
	Country     string  `json:"Country"`
	Latitude    float64 `json:"Latitude"`
	Longitude   float64 `json:"Longitude"`
	AddressType int     `json:"AddressType"` // "1" = User address, "4" = User lat/long
}

type ReportAttribute struct {
	Attribute int    `json:"Attribute"`
	Value     string `json:"Value"`
}

type CreditCardData struct {
	FirstName       string `json:"CardFirstName,omitempty"`
	LastName        string `json:"CardLastName,omitempty"`
	ExpirationMonth int    `json:"ExpirationMonth,omitempty"`
	ExpirationYear  int    `json:"ExpirationYear,omitempty"`
	Number          string `json:"CreditCardNumber,omitempty"`
	Type            int    `json:"CreditCardType,omitempty"` // "1" = Amex, "2" = Visa, "3" = MC, "4" = Disc
}
