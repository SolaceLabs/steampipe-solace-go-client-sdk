package solacesdk

import (
	"net/url"
)

type RequestConfig struct {
	Endpoint   string
	Params     *url.Values
	Pagination *RequestPagination
	Body       interface{}
}

func NewRequestConfig(endpoint string) *RequestConfig {
	return &RequestConfig{
		Endpoint: endpoint,
		Params:   &url.Values{},
	}
}

type RequestPagination struct {
	pageNumber int
	pageSize   int
	paginate   bool
}

func (rp *RequestPagination) NextPage(newPageSize int) {
	rp.pageNumber += 1

	if newPageSize != -1 {
		rp.pageSize = newPageSize
	}
}

func NewRequestPagination() *RequestPagination {
	return &RequestPagination{
		pageNumber: 1,
		pageSize:   defaultPageSize,
		paginate:   true,
	}
}

func NewRequestNoPagination() *RequestPagination {
	return &RequestPagination{
		pageNumber: 1,
		pageSize:   defaultPageSize,
		paginate:   false,
	}
}

func NewRequestSingleElementPagination(paginate bool) *RequestPagination {
	return &RequestPagination{
		pageNumber: 1,
		pageSize:   1,
		paginate:   paginate,
	}
}

type Pagination struct {
	SortBy string `json:"sortBy"`
}

type Meta struct {
	Pagination struct {
		PageNumber int `json:"pageNumber"`
		Count      int `json:"count"`
		PageSize   int `json:"pageSize"`
		NextPage   int `json:"nextPage"`
		TotalPages int `json:"totalPages"`
	} `json:"pagination"`
}
