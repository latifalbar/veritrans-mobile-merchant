package merchant

//APIResponse response of the API call
type APIResponse struct {
	StatusCode    string `json:"status_code"`
	StatusMessage string `json:"status_message"`
	TokenID       string `json:"token_id"`
	Bank          string `json:"bank"`
	RedirectURL   string `json:"redirect_url"`
}

// Promotion struct
type Promotion struct {
	InstallmentList []Installment `json:"installment"`
	DiscountList    []Discount    `json:"discount"`
}

// Discount is a model for promotion
type Discount struct {
	Title       string   `json:"title"`
	Description string   `json:"description"`
	Percentage  int64    `json:"discount_percentage"`
	Bins        []string `json:"bins"`
}

// Installment is a model for installment promotion
type Installment struct {
	Title            string   `json:"title"`
	Description      string   `json:"description"`
	Percentage       int64    `json:"discount_percentage"`
	Bins             []string `json:"bins"`
	InstallmentTerms []string `json:"installment_terms"`
}

// AuthenticatedModel is a model for generated authentication
type AuthenticatedModel struct {
	Token string `json:"token"`
	Cards []Card `json:"cards"`
}

// CardRequest is a model for card request details
type CardRequest struct {
	StatusCode   string `json:"status_code"`
	SavedTokenID string `json:"saved_token_id"`
	MaskedCard   string `json:"masked_card"`
}

// Card is a model for card details
type Card struct {
	SavedTokenID string `json:"saved_token_id"`
	MaskedCard   string `json:"masked_card"`
}
