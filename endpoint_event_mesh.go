package solacesdk

import (
	"fmt"
	"time"
)

type EventMesh struct {
	CreatedTime   time.Time `json:"createdTime"`
	UpdatedTime   time.Time `json:"updatedTime"`
	CreatedBy     string    `json:"createdBy"`
	ChangedBy     string    `json:"changedBy"`
	Id            string    `json:"id"`
	Name          string    `json:"name"`
	EnvironmentId string    `json:"environmentId"`
	Description   string    `json:"description"`
	BrokerType    string    `json:"brokerType"`
	Type          string    `json:"type"`
}

type EventMeshListResponse struct {
	EventMeshs []EventMesh `json:"data"`
	Meta       Meta        `json:"meta"`
}

type EventMeshListPaginator struct {
	firstPage  bool
	pageNumber int
	count      int
	pageSize   int
	nextPage   int
	totalPages int
	config     *RequestConfig
	client     *Client
}

func (lp *EventMeshListPaginator) HasNoMorePages() bool {
	return lp.nextPage == 0 || lp.pageNumber > lp.totalPages
}

func (lp *EventMeshListPaginator) NextPage() ([]EventMesh, Meta, error) {
	if !lp.firstPage && lp.HasNoMorePages() {
		return nil, Meta{}, fmt.Errorf("no more pages available")
	}
	var r = &EventMeshListResponse{}
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

	return r.EventMeshs, r.Meta, nil
}

func (at *Client) NewEventMeshListPaginator() *EventMeshListPaginator {
	var config = NewRequestConfig("architecture/eventMeshes")
	config.Pagination = NewRequestPagination()

	var p = &EventMeshListPaginator{
		firstPage:  true,
		pageNumber: 1,
		pageSize:   defaultPageSize,
		config:     config,
		client:     at,
	}

	return p
}

type EventMeshGetResponse struct {
	EventMesh EventMesh `json:"data"`
	Meta      Meta      `json:"meta"`
}

func (at *Client) GetEventMesh(eventMeshId string) (*EventMesh, error) {
	var config = NewRequestConfig(fmt.Sprintf(`architecture/eventMeshes/%s`, eventMeshId))

	var r = &EventMeshGetResponse{}
	var _, err = at.Get(config, &r)
	if err != nil {
		return nil, at.handleKnownErrors(err)
	}

	return &r.EventMesh, nil
}
