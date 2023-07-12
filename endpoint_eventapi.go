package solacesdk

import (
	"fmt"
	"time"
)

type EventApi struct {
	CreatedTime         time.Time `json:"createdTime"`
	UpdatedTime         time.Time `json:"updatedTime"`
	CreatedBy           string    `json:"createdBy"`
	ChangedBy           string    `json:"changedBy"`
	Id                  string    `json:"id"`
	Name                string    `json:"name"`
	Shared              bool      `json:"shared"`
	ApplicationDomainId string    `json:"applicationDomainId"`
	NumberOfVersions    int       `json:"numberOfVersions"`
	BrokerType          string    `json:"brokerType"`
	CustomAttributes    []struct {
		CustomAttributeDefinitionId   string `json:"customAttributeDefinitionId"`
		CustomAttributeDefinitionName string `json:"customAttributeDefinitionName"`
		Value                         string `json:"value"`
	} `json:"customAttributes"`
	Type string `json:"type"`
}

type EventApiListResponse struct {
	EventApis []EventApi `json:"data"`
	Meta      Meta       `json:"meta"`
}

type EventApiListPaginator struct {
	firstPage  bool
	pageNumber int
	count      int
	pageSize   int
	nextPage   int
	totalPages int
	config     *RequestConfig
	client     *Client
}

func (lp *EventApiListPaginator) HasNoMorePages() bool {
	return lp.nextPage == 0 || lp.pageNumber > lp.totalPages
}

func (lp *EventApiListPaginator) NextPage() ([]EventApi, Meta, error) {
	if !lp.firstPage && lp.HasNoMorePages() {
		return nil, Meta{}, fmt.Errorf("no more pages available")
	}
	var r = &EventApiListResponse{}
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

	return r.EventApis, r.Meta, nil
}

func (at *Client) NewEventApiListPaginator(args map[string]string) *EventApiListPaginator {
	var config = NewRequestConfig("architecture/eventApis")
	config.Pagination = NewRequestPagination()
	ids, exists := args["ids"]
	if exists && len(ids) > 0 {
		config.Params.Add("ids", ids)
	}
	applicationDomainIds, exists := args["applicationDomainIds"]
	if exists && len(applicationDomainIds) > 0 {
		config.Params.Add("applicationDomainIds", applicationDomainIds)
	}

	var p = &EventApiListPaginator{
		firstPage:  true,
		pageNumber: 1,
		pageSize:   defaultPageSize,
		config:     config,
		client:     at,
	}

	return p
}

type EventApiGetResponse struct {
	EventApi EventApi `json:"data"`
	Meta     Meta     `json:"meta"`
}

func (at *Client) GetEventApi(eventApiId string) (*EventApi, error) {
	var config = NewRequestConfig(fmt.Sprintf(`architecture/eventApis/%s`, eventApiId))

	var r = &EventApiGetResponse{}
	var _, err = at.Get(config, &r)
	if err != nil {
		return nil, at.handleKnownErrors(err)
	}

	return &r.EventApi, nil
}
