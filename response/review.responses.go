package response

// type ReviewResponses struct {
// 	ID          int               `json:"id"`
// 	User        UserRespReview    `json:"user"`
// 	Product     ProductRespReview `json:"product"`
// 	Rating      int               `json:"rating"`
// 	TextReview  string            `json:"text_review"`
// 	ImageReview []string          `json:"image_review"`
// 	Created_at  int64             `json:"created_at"`
// 	Updated_at  int64             `json:"updated_at"`
// }

type ReviewProductResp struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Rating   int    `json:"rating"`
	Text     string `json:"text_review"`
}

// type ReviewResponses struct {
// 	ID          int64    `json:"id"`
// 	User        string   `json:"user"`    // จาก username
// 	Product     string   `json:"product"` // จาก product_name
// 	Rating      int      `json:"rating"`
// 	TextReview  string   `json:"text_review"` // จาก description
// 	Description string   `json:"description"`
// 	ImageReview []string `json:"image_review"`
// 	CreatedAt   string   `json:"created_at"`
// 	UpdatedAt   string   `json:"updated_at"`
// }

type ReviewResponses struct {
	ID          int64  `bun:"id" json:"id"`
	User        string `bun:"user" json:"user"`
	Product     string `bun:"product" json:"product"`
	Rating      int    `bun:"rating" json:"rating"`
	Description string `bun:"description" json:"description"`
	ImageReview []string `json:"image_review"`
	CreatedAt   string `bun:"created_at" json:"created_at"`
	UpdatedAt   string `bun:"updated_at" json:"updated_at"`
}
