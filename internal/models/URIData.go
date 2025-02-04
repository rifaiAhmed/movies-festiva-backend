package models

type UriId struct {
	ID int `uri:"id" binding:"required"`
}

type ComponentServerSide struct {
	Limit     int    `json:"limit"`
	Skip      int    `json:"skip"`
	SortType  string `json:"sort_type"`
	SortBy    string `json:"sort_by"`
	Search    string `json:"search"`
	Offset    int    `json:"offset"`
	Condition string `json:"condition"`
	From      string `json:"from"`
	To        string `json:"to"`
}
