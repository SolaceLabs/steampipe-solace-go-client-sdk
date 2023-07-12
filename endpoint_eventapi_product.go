package solacesdk

import (
	"fmt"
	"time"
)

type EventApiProduct struct {
	CreatedTime         time.Time `json:"createdTime"`
	UpdatedTime         time.Time `json:"updatedTime"`
	CreatedBy           string    `json:"createdBy"`
	ChangedBy           string    `json:"changedBy"`
	Id                  string    `json:"id"`
	Name                string    `json:"name"`
	ApplicationDomainId string    `json:"applicationDomainId"`
	Shared              bool      `json:"shared"`
	NumberOfVersions    int       `json:"numberOfVersions"`
	BrokerType          string    `json:"brokerType"`
	Type                string    `json:"type"`
	CustomAttributes    []struct {
		CustomAttributeDefinitionId   string `json:"customAttributeDefinitionId"`
		CustomAttributeDefinitionName string `json:"customAttributeDefinitionName"`
		Value                         string `json:"value"`
	} `json:"customAttributes"`
}

type EventApiProductListResponse struct {
	EventApiProducts []EventApiProduct `json:"data"`
	Meta             Meta              `json:"meta"`
}

type EventApiProductListPaginator struct {
	firstPage  bool
	pageNumber int
	count      int
	pageSize   int
	nextPage   int
	totalPages int
	config     *RequestConfig
	client     *Client
}

func (lp *EventApiProductListPaginator) HasNoMorePages() bool {
	return lp.nextPage == 0 || lp.pageNumber > lp.totalPages
}

func (lp *EventApiProductListPaginator) NextPage() ([]EventApiProduct, Meta, error) {
	if !lp.firstPage && lp.HasNoMorePages() {
		return nil, Meta{}, fmt.Errorf("no more pages available")
	}
	var r = &EventApiProductListResponse{}
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

	return r.EventApiProducts, r.Meta, nil
}

func (at *Client) NewEventApiProductListPaginator(args map[string]string) *EventApiProductListPaginator {
	var config = NewRequestConfig("architecture/eventApiProducts")
	config.Pagination = NewRequestPagination()
	ids, exists := args["ids"]
	if exists && len(ids) > 0 {
		config.Params.Add("ids", ids)
	}
	applicationDomainIds, exists := args["applicationDomainIds"]
	if exists && len(applicationDomainIds) > 0 {
		config.Params.Add("applicationDomainIds", applicationDomainIds)
	}

	var p = &EventApiProductListPaginator{
		firstPage:  true,
		pageNumber: 1,
		pageSize:   defaultPageSize,
		config:     config,
		client:     at,
	}

	return p
}

type EventApiProductGetResponse struct {
	EventApiProduct EventApiProduct `json:"data"`
	Meta            Meta            `json:"meta"`
}

func (at *Client) GetEventApiProduct(eventApiProductId string) (*EventApiProduct, error) {
	var config = NewRequestConfig(fmt.Sprintf(`architecture/eventApiProducts/%s`, eventApiProductId))

	var r = &EventApiProductGetResponse{}
	var _, err = at.Get(config, &r)
	if err != nil {
		return nil, at.handleKnownErrors(err)
	}

	return &r.EventApiProduct, nil
}
