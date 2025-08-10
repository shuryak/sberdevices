package sbertypes

type Pagination struct {
	Limit  string `json:"limit"`
	Offset string `json:"offset"`
	Total  string `json:"total"`
}
