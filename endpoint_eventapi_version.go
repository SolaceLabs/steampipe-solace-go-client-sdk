package solacesdk

import (
	"fmt"
	"time"
)

type EventApiVersion struct {
	CreatedTime                       time.Time `json:"createdTime"`
	UpdatedTime                       time.Time `json:"updatedTime"`
	CreatedBy                         string    `json:"createdBy"`
	ChangedBy                         string    `json:"changedBy"`
	Id                                string    `json:"id"`
	EventApiId                        string    `json:"eventApiId"`
	Description                       string    `json:"description"`
	Version                           string    `json:"version"`
	DisplayName                       string    `json:"displayName"`
	ProducedEventVersionIds           []string  `json:"producedEventVersionIds"`
	ConsumedEventVersionIds           []string  `json:"consumedEventVersionIds"`
	DeclaredEventApiProductVersionIds []string  `json:"declaredEventApiProductVersionIds"`
	CustomAttributes                  []struct {
		CustomAttributeDefinitionId   string `json:"customAttributeDefinitionId"`
		CustomAttributeDefinitionName string `json:"customAttributeDefinitionName"`
		Value                         string `json:"value"`
	} `json:"customAttributes"`
	StateId string `json:"stateId"`
	Type    string `json:"type"`
}

type EventApiVersionListResponse struct {
	EventApiVersions []EventApiVersion `json:"data"`
	Meta             Meta              `json:"meta"`
}

type EventApiVersionListPaginator struct {
	firstPage  bool
	pageNumber int
	count      int
	pageSize   int
	nextPage   int
	totalPages int
	config     *RequestConfig
	client     *Client
}

func (lp *EventApiVersionListPaginator) HasNoMorePages() bool {
	return lp.nextPage == 0 || lp.pageNumber > lp.totalPages
}

func (lp *EventApiVersionListPaginator) NextPage() ([]EventApiVersion, Meta, error) {
	if !lp.firstPage && lp.HasNoMorePages() {
		return nil, Meta{}, fmt.Errorf("no more pages available")
	}
	var r = &EventApiVersionListResponse{}
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

	return r.EventApiVersions, r.Meta, nil
}

func (at *Client) NewEventApiVersionListPaginator(args map[string]string) *EventApiVersionListPaginator {
	var config = NewRequestConfig("architecture/eventApiVersions")
	config.Pagination = NewRequestPagination()
	ids, exists := args["ids"]
	if exists && len(ids) > 0 {
		config.Params.Add("ids", ids)
	}
	eventApiIds, exists := args["eventApiIds"]
	if exists && len(eventApiIds) > 0 {
		config.Params.Add("eventApiIds", eventApiIds)
	}

	var p = &EventApiVersionListPaginator{
		firstPage:  true,
		pageNumber: 1,
		pageSize:   defaultPageSize,
		config:     config,
		client:     at,
	}

	return p
}

type EventApiVersionGetResponse struct {
	EventApiVersion EventApiVersion `json:"data"`
	Meta            Meta            `json:"meta"`
}

func (at *Client) GetEventApiVersion(eventApiVersionId string) (*EventApiVersion, error) {
	var config = NewRequestConfig(fmt.Sprintf(`architecture/eventApiVersions/%s`, eventApiVersionId))

	var r = &EventApiVersionGetResponse{}
	var _, err = at.Get(config, &r)
	if err != nil {
		return nil, at.handleKnownErrors(err)
	}

	return &r.EventApiVersion, nil
}
