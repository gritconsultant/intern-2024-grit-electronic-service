package response

type PaymentResponses struct {
	ID            int                   `json:"id"`
	UpdatedBy     int                   `json:"updated_by"`
	Admin         AdminPaymentResp      `bun:"admin"`
	SystemBank    SystemBankRespPayment `bun:"systembank"`
	Price         float64               `json:"price"`
	Status        string                `json:"status"`
	Image         string                `json:"image"`
	BankName      string                `json:"bank_name"`
	AccountName   string                `json:"account_name"`
	AccountNumber string                `json:"account_number"`
	Created_at    int64                 `json:"created_at"`
	Updated_at    int64                 `json:"updated_at"`
}

type PaymentUserResp struct {
}

type PaymentRespOrderDetail struct {
	ID     int     `json:"id"`
	Price  float64 `json:"price"`
	Status string  `json:"status"`
}
