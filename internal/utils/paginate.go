package utils

import (
	"fmt"
	"net/url"
	"strconv"
)

type PaginationParams struct {
	Skip    int64
	Limit   int64
	PerPage int64
	Page    int64
}

type PageInfo struct {
	TotalItems  int64 `json:"total_items"`
	TotalPages  int64 `json:"total_pages"`
	CurrentPage int64 `json:"current_page"`
	HasNextPage bool  `json:"has_next_page"`
}

// Parses the pagination params from the request query params
// Returns the pagination params and an error
func ParsePaginationQueryParams(queryParams url.Values) (PaginationParams, error) {
	perPage := queryParams["perPage"]
	page := queryParams["page"]

	if len(perPage) == 0 {
		perPage = []string{"20"}
	}
	if len(page) == 0 {
		page = []string{"1"}
	}

	var skip int64 = 0
	var limit int64 = 0
	pageInt, err := strconv.Atoi(page[0])
	if err != nil {
		return PaginationParams{}, fmt.Errorf("error parsing 'page' query with value %v", page[0])
	}
	perPageInt, err := strconv.Atoi(perPage[0])
	if err != nil {
		return PaginationParams{}, fmt.Errorf("error parsing 'perPage' query with value %v", perPage[0])
	}
	skip = int64((pageInt - 1) * perPageInt)
	limit = int64(perPageInt)
	return PaginationParams{
		Skip:    skip,
		Limit:   limit,
		PerPage: int64(perPageInt),
		Page:    int64(pageInt),
	}, nil
}
