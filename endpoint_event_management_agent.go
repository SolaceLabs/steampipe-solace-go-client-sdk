package solace

import (
	"fmt"
)

type EventManagementAgent struct {
	CreatedTime                     string   `json:"createdTime"`
	UpdatedTime                     string   `json:"updatedTime"`
	CreatedBy                       string   `json:"createdBy"`
	ChangedBy                       string   `json:"changedBy"`
	Id                              string   `json:"id"`
	Name                            string   `json:"name"`
	Region                          string   `json:"region"`
	ClientUsername                  string   `json:"clientUsername"`
	ClientPassword                  string   `json:"clientPassword"`
	ReferencedByMessagingServiceIds []string `json:"referencedByMessagingServiceIds"`
	OrgId                           string   `json:"orgId"`
	Status                          string   `json:"status"`
	LastConnectedTime               string   `json:"lastConnectedTime"`
	EventManagementAgentRegionId    string   `json:"eventManagementAgentRegionId"`
	Type                            string   `json:"type"`
}

type EventManagementAgentListResponse struct {
	EventManagementAgents []EventManagementAgent `json:"data"`
	Meta                  Meta                   `json:"meta"`
}

type EventManagementAgentListPaginator struct {
	firstPage  bool
	pageNumber int
	count      int
	pageSize   int
	nextPage   int
	totalPages int
	config     *RequestConfig
	client     *Client
}

func (lp *EventManagementAgentListPaginator) HasNoMorePages() bool {
	return lp.nextPage == 0 || lp.pageNumber > lp.totalPages
}

func (lp *EventManagementAgentListPaginator) NextPage() ([]EventManagementAgent, Meta, error) {
	if !lp.firstPage && lp.HasNoMorePages() {
		return nil, Meta{}, fmt.Errorf("no more pages available")
	}
	var r = &EventManagementAgentListResponse{}
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

	return r.EventManagementAgents, r.Meta, nil
}

func (at *Client) NewEventManagementAgentListPaginator() *EventManagementAgentListPaginator {
	var config = NewRequestConfig("architecture/eventManagementAgents")
	config.Pagination = NewRequestPagination()

	var p = &EventManagementAgentListPaginator{
		firstPage:  true,
		pageNumber: 1,
		pageSize:   defaultPageSize,
		config:     config,
		client:     at,
	}

	return p
}

type EventManagementAgentGetResponse struct {
	EventManagementAgent EventManagementAgent `json:"data"`
	Pagination           Pagination           `json:"meta"`
}

func (at *Client) GetEventManagementAgent(eventManagementAgentId string) (*EventManagementAgent, error) {
	var config = NewRequestConfig(fmt.Sprintf(`architecture/eventManagementAgents/%s`, eventManagementAgentId))
	config.Params.Add("expand", "broker")
	config.Params.Add("expand", "serviceConnectionEndpoints")

	var result = &EventManagementAgentGetResponse{}
	var _, err = at.Get(config, &result)
	if err != nil {
		return nil, at.handleKnownErrors(err)
	}

	return &result.EventManagementAgent, nil
}
