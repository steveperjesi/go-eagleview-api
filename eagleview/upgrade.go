package eagleview

import (
	"time"
)

type GetReportUpgradeProductsRequest struct {
	ReportID int
}

type PlaceUpgradeOrderRequest struct {
	ReportID int
}

type GetReportUpgradeProductsResponse struct {
	AvailableProducts []AvailableProducts
	CurrentOrder      UpgradeProducts
}

type PlaceUpgradeOrderResponse struct {
	ReportID int
}

type UpgradeProducts struct {
	ReportID                   int
	MeasurementInstructionType int
	ProductID                  int
	ProductDeliveryID          int
	AddOnProductIDs            []int
	AdditionalEmails           string
	ClaimNumber                string
	ClaimInfo                  string
	PONumber                   string
	DateOfLoss                 time.Time
	CatID                      string
}
