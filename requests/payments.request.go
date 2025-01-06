package requests

type PaymentRequest struct {
	Page   int64  `form:"page"`
	Size   int64  `form:"size"`
	Search string `form:"search"`
}

type PaymentIdRequest struct {
	ID int64 `uri:"id"`
}

type PaymentCreateRequest struct {
	Price         int64  `json:"price"`
	Amount        int64  `json:"amount"`
	Status        string `json:"status"`
	Slip          string `json:"slip"`
	BankName      string `json:"bank_name"`
	AccountName   string `json:"account_name"`
	AccountNumber string `json:"account_number"`
}

type PaymentUpdateRequest struct {
	Id            int64  `json:"id"`
	Price         int64  `json:"price"`
	Amount        int64  `json:"amount"`
	Status        string `json:"status"`
	Slip          string `json:"slip"`
	BankName      string `json:"bank_name"`
	AccountName   string `json:"account_name"`
	AccountNumber string `json:"account_number"`
}
