package solacesdk

import (
	"fmt"
)

type LifecycleState struct {
	Id          string `json:"id"`
	Description string `json:"description"`
	Name        string `json:"name"`
	StateOrder  int    `json:"stateOrder"`
	Type        string `json:"type"`
}

type LifecycleStateListResponse struct {
	LifecycleStates []LifecycleState `json:"data"`
	Meta            Meta             `json:"meta"`
}

type LifecycleStateListPaginator struct {
	firstPage  bool
	pageNumber int
	count      int
	pageSize   int
	nextPage   int
	totalPages int
	config     *RequestConfig
	client     *Client
}

func (lp *LifecycleStateListPaginator) HasNoMorePages() bool {
	return lp.nextPage == 0 || lp.pageNumber > lp.totalPages
}

func (lp *LifecycleStateListPaginator) NextPage() ([]LifecycleState, Meta, error) {
	if !lp.firstPage && lp.HasNoMorePages() {
		return nil, Meta{}, fmt.Errorf("no more pages available")
	}
	var r = &LifecycleStateListResponse{}
	var _, err = lp.client.Get(lp.config, r)
	if err != nil {
		return nil, Meta{}, lp.client.handleKnownErrors(err)
	}

	if lp.firstPage == true {
		lp.firstPage = false
		lp.count = r.Meta.Pagination.Count
		lp.totalPages = r.Meta.Pagination.TotalPages
	}
	lp.nextPage = r.Meta.Pagination.NextPage
	lp.config.Pagination.NextPage(lp.pageSize)

	return r.LifecycleStates, r.Meta, nil
}

func (at *Client) NewLifecycleStateListPaginator() *LifecycleStateListPaginator {
	var config = NewRequestConfig("architecture/states")
	config.Pagination = NewRequestNoPagination()

	var p = &LifecycleStateListPaginator{
		firstPage:  true,
		pageNumber: 1,
		pageSize:   defaultPageSize,
		config:     config,
		client:     at,
	}

	return p
}
