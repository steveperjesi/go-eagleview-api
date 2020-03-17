package eagleview

// type PlaceOrderRequest struct {
// 	OrderReports   []Report       `json:"OrderReports"`
// 	PromoCode      string         `json:"PromoCode"`
// 	PlaceOrderUser string         `json:"PlaceOrderUser,omitempty"`
// 	CreditCardData CreditCardData `json:"CreditCardData"`
// }

// {
//     "OrderId": 26998516,
//     "ReportIds": [
//         41001465
//     ]
// }

type PlaceOrderResponse struct {
	OrderID   int   `json:"OrderId"`
	ReportIDs []int `json:"ReportIds"`
}
