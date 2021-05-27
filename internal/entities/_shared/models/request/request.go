package request

type SortItem struct {
	Field string `json:"field"`
	Dir   string `json:"dir"`
}

type PagingItem struct {
	PageSize   int `json:"pageSize"`
	PageNumber int `json:"pageNumber"`
}

type CustomFilterItem struct {
	Name   string   `json:"name"`
	Values []string `json:"values"`
}

type NewGridList struct {
	PageSize      int                `json:"pageSize"`
	PageNumber    int                `json:"pageNumber"`
	HiddenFilters []string           `json:"hiddenFilters"`
	CustomFilters []CustomFilterItem `json:"customFilters"`
	Sorts         []SortItem         `json:"sorts"`
}
