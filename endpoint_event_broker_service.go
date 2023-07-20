package solace

import (
	"fmt"
	"time"
)

type EventBrokerService struct {
	CreatedTime         time.Time `json:"createdTime"`
	UpdatedTime         time.Time `json:"updatedTime"`
	CreatedBy           string    `json:"createdBy"`
	ChangedBy           string    `json:"changedBy"`
	Id                  string    `json:"id"`
	Type                string    `json:"type"`
	Name                string    `json:"name"`
	OwnedBy             string    `json:"ownedBy"`
	InfrastructureId    string    `json:"infrastructureId"`
	DatacenterId        string    `json:"datacenterId"`
	ServiceClassId      string    `json:"serviceClassId"`
	EventMeshId         string    `json:"eventMeshId"`
	OngoingOperationIds []string  `json:"ongoingOperationIds"`
	AdminState          string    `json:"adminState"`
	CreationState       string    `json:"creationState"`
	Locked              bool      `json:"locked"`
}

type EventBrokerServiceListResponse struct {
	EventBrokerServices []EventBrokerService `json:"data"`
	Meta                Meta                 `json:"meta"`
}

type EventBrokerServiceListPaginator struct {
	firstPage  bool
	pageNumber int
	count      int
	pageSize   int
	nextPage   int
	totalPages int
	config     *RequestConfig
	client     *Client
}

func (lp *EventBrokerServiceListPaginator) HasNoMorePages() bool {
	return lp.nextPage == 0 || lp.pageNumber > lp.totalPages
}

func (lp *EventBrokerServiceListPaginator) NextPage() ([]EventBrokerService, Meta, error) {
	if !lp.firstPage && lp.HasNoMorePages() {
		return nil, Meta{}, fmt.Errorf("no more pages available")
	}
	var r = &EventBrokerServiceListResponse{}
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

	return r.EventBrokerServices, r.Meta, nil
}

func (at *Client) NewEventBrokerServiceListPaginator() *EventBrokerServiceListPaginator {
	var config = NewRequestConfig("missionControl/eventBrokerServices")
	config.Pagination = NewRequestPagination()

	var p = &EventBrokerServiceListPaginator{
		firstPage:  true,
		pageNumber: 1,
		pageSize:   defaultPageSize,
		config:     config,
		client:     at,
	}

	return p
}

type EventBrokerServiceGetResponse struct {
	EventBrokerService EventBrokerService `json:"data"`
	Pagination         Pagination         `json:"meta"`
}

func (at *Client) GetEventBrokerService(eventBrokerServiceId string) (*EventBrokerService, error) {
	var config = NewRequestConfig(fmt.Sprintf(`missionControl/eventBrokerServices/%s`, eventBrokerServiceId))
	config.Params.Add("expand", "broker")
	config.Params.Add("expand", "serviceConnectionEndpoints")

	var result = &EventBrokerServiceGetResponse{}
	var _, err = at.Get(config, &result)
	if err != nil {
		return nil, at.handleKnownErrors(err)
	}

	return &result.EventBrokerService, nil
}
