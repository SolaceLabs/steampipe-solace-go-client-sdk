package solace

import (
	"fmt"
	"time"
)

type Environment struct {
	CreatedTime         time.Time `json:"createdTime"`
	UpdatedTime         time.Time `json:"updatedTime"`
	CreatedBy           string    `json:"createdBy"`
	ChangedBy           string    `json:"changedBy"`
	Id                  string    `json:"id"`
	Name                string    `json:"name"`
	Description         string    `json:"description"`
	Revision            int       `json:"revision"`
	NumberOfEventMeshes int       `json:"numberOfEventMeshes"`
	Type                string    `json:"type"`
}

type EnvironmentListResponse struct {
	Environments []Environment `json:"data"`
	Meta         Meta          `json:"meta"`
}

type EnvironmentListPaginator struct {
	firstPage  bool
	pageNumber int
	count      int
	pageSize   int
	nextPage   int
	totalPages int
	config     *RequestConfig
	client     *Client
}

func (lp *EnvironmentListPaginator) HasNoMorePages() bool {
	return lp.nextPage == 0 || lp.pageNumber > lp.totalPages
}

func (lp *EnvironmentListPaginator) NextPage() ([]Environment, Meta, error) {
	if !lp.firstPage && lp.HasNoMorePages() {
		return nil, Meta{}, fmt.Errorf("no more pages available")
	}
	var r = &EnvironmentListResponse{}
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

	return r.Environments, r.Meta, nil
}

func (at *Client) NewEnvironmentListPaginator() *EnvironmentListPaginator {
	var config = NewRequestConfig("architecture/environments")
	config.Pagination = NewRequestPagination()

	var p = &EnvironmentListPaginator{
		firstPage:  true,
		pageNumber: 1,
		pageSize:   defaultPageSize,
		config:     config,
		client:     at,
	}

	return p
}

type EnvironmentGetResponse struct {
	Environment Environment `json:"data"`
	Meta        Meta        `json:"meta"`
}

func (at *Client) GetEnvironment(environmentId string) (*Environment, error) {
	var config = NewRequestConfig(fmt.Sprintf(`architecture/environments/%s`, environmentId))

	var r = &EnvironmentGetResponse{}
	var _, err = at.Get(config, &r)
	if err != nil {
		return nil, at.handleKnownErrors(err)
	}

	return &r.Environment, nil
}
