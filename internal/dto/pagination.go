package dto
type PaginatedResponse struct {
	Page  int         `json:"page"`
	Limit int         `json:"limit"`
	Total int64       `json:"total"`
	Items interface{} `json:"items"`
}
