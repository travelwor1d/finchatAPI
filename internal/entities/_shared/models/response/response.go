package response

type Paging struct {
	Page         int `json:"page"`
	Size         int `json:"size"`
	TotalCount   int `json:"totalCount"`
	CurrentCount int `json:"currentCount"`
}

type PageItem struct {
	Label string `json:"label"`
	Value int    `json:"value"`
}
