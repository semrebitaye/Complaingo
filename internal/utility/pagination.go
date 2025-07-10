package utility

import (
	"encoding/json"
	"fmt"
	"strconv"
)

type PaginationParam struct {
	Page    string `json:"page"`
	PerPage string `json:"per_page"`
	Sort    string `json:"sort"`
	Search  string `json:"search"`
	Filter  string `json:"filter"`
}

type Filter struct {
	ColumnName string `json:"column_name"`
	Operator   string `json:"operator"`
	Value      any    `json:"value"`
}

type Sort struct {
	ColumnName string `json:"column_name"`
	Value      string `json:"value"`
}

type FilterParam struct {
	Page    int      `json:"page"`
	PerPage int      `json:"per_page"`
	Sort    Sort     `json:"sort"`
	Search  string   `json:"search"`
	Filters []Filter `json:"filters"`
}

func ExtractPagination(param PaginationParam) (FilterParam, error) {
	page, err := strconv.Atoi(param.Page)
	if err != nil && page <= 0 {
		page = 1 //default page number
	}

	per_page, err := strconv.Atoi(param.PerPage)
	if err != nil && per_page <= 0 {
		per_page = 10 //default limit
	}

	var sort Sort
	if param.Sort == "" {
		sort.ColumnName = "id"
		sort.Value = "asc"
	} else {
		err := json.Unmarshal([]byte(param.Sort), &sort)
		if err != nil {
			return FilterParam{}, fmt.Errorf("invalid sort format: %w", err)
		}

		// validate sort value
		if sort.Value != "asc" && sort.Value != "desc" {
			return FilterParam{}, fmt.Errorf("invalid sort direction: must be 'asc' or 'desc'")
		}
	}

	var filter []Filter
	if param.Filter != "" {
		err := json.Unmarshal([]byte(param.Filter), &filter)
		if err != nil {
			return FilterParam{}, fmt.Errorf("invalid filter format: %w", err)
		}
	}

	return FilterParam{
		Page:    page,
		PerPage: per_page,
		Sort:    sort,
		Search:  param.Search,
		Filters: filter,
	}, nil
}
