package eagleview

// TODO: Remove the "AvailableProducts" struct as []Product works just fine

type AvailableProducts struct {
	Products []Product
}

type Product struct {
	ProductID                  int               `json:"productID"`
	Name                       string            `json:"name"`
	Description                string            `json:"description"`
	ProductGroup               string            `json:"productGroup,omitempty"` // deprecated
	TemporarilyUnavailable     bool              `json:"isTemporarilyUnavailable"`
	PriceMin                   float64           `json:"priceMin"`
	PriceMax                   float64           `json:"priceMax"`
	DeliveryProducts           []DeliveryProduct `json:"deliveryProducts"`
	AddOnProducts              []DeliveryProduct `json:"addOnProducts"`
	MeasurementInstructionType []int             `json:"measurementInstructionTypes"`
	TypeOfStructure            int               `json:"TypeOfStructure"`
	RoofProduct                bool              `json:"IsRoofProduct"`
	SortOrder                  int               `json:"SortOrder"`
	AllowsUserSubmittedPhotos  bool              `json:"AllowsUserSubmittedPhotos"`
	DetailedDescription        string            `json:"DetailedDescription"`
}

type DeliveryProduct struct {
	ProductID              int     `json:"productID"`
	Name                   string  `json:"name"`
	Description            string  `json:"description"`
	TemporarilyUnavailable bool    `json:"isTemporarilyUnavailable"`
	PriceMin               float64 `json:"priceMin"`
	PriceMax               float64 `json:"priceMax"`
	DeliveryProductIDs     []int   `json:"deliveryProductIds"`
	DetailedDescription    string  `json:"DetailedDescription"`
}
