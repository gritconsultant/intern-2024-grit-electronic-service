package requests

type CategoryRequest struct {
	Page   int64  `form:"page"`
	Size   int64  `form:"size"`
	Search string `form:"search"`
}

type CategoryIdRequest struct {
	ID int `uri:"id"`
}

type CategoryCreateRequest struct {
	Name     string `json:"name"`
	IsActive bool   `json:"is_active"`
}

type CategoryUpdateRequest struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	IsActive bool   `json:"is_active"`
}
