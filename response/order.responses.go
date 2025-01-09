package response

type OrderResponses struct {
	ID           int                    `json:"id"`
	TotalPrice  float64                `json:"total_price"`
	TotalAmount int                  `json:"total_amount"`
	Status       string                 `json:"status"`
	Created_at   int64                  `json:"created_at"`
	Updated_at   int64                  `json:"updated_at"`
}
