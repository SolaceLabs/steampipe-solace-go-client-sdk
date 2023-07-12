package solace

import (
	"fmt"
	"time"
)

type Configuration struct {
	CreatedTime         time.Time `json:"createdTime"`
	UpdatedTime         time.Time `json:"updatedTime"`
	CreatedBy           string    `json:"createdBy"`
	ChangedBy           string    `json:"changedBy"`
	Id                  string    `json:"id"`
	ContextType         string    `json:"contextType"`
	ContextId           string    `json:"contextId"`
	ConfigurationTypeId string    `json:"configurationTypeId"`
	EntityType          string    `json:"entityType"`
	EntityId            string    `json:"entityId"`
	Value               struct{}  `json:"value"`
	Type                string    `json:"type"`
}

type ConfigurationListResponse struct {
	Configurations []Configuration `json:"data"`
	Meta           Meta            `json:"meta"`
}

type ConfigurationListPaginator struct {
	firstPage  bool
	pageNumber int
	count      int
	pageSize   int
	nextPage   int
	totalPages int
	config     *RequestConfig
	client     *Client
}

func (lp *ConfigurationListPaginator) HasNoMorePages() bool {
	return lp.nextPage == 0 || lp.pageNumber > lp.totalPages
}

func (lp *ConfigurationListPaginator) NextPage() ([]Configuration, Meta, error) {
	if !lp.firstPage && lp.HasNoMorePages() {
		return nil, Meta{}, fmt.Errorf("no more pages available")
	}
	var r = &ConfigurationListResponse{}
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

	return r.Configurations, r.Meta, nil
}

func (at *Client) NewConfigurationListPaginator() *ConfigurationListPaginator {
	var config = NewRequestConfig("architecture/configurations")
	config.Pagination = NewRequestPagination()

	var p = &ConfigurationListPaginator{
		firstPage:  true,
		pageNumber: 1,
		pageSize:   defaultPageSize,
		config:     config,
		client:     at,
	}

	return p
}

type ConfigurationGetResponse struct {
	Configuration Configuration `json:"data"`
	Meta          Meta          `json:"meta"`
}

func (at *Client) GetConfiguration(configurationId string) (*Configuration, error) {
	var config = NewRequestConfig(fmt.Sprintf(`architecture/configurations/%s`, configurationId))

	var r = &ConfigurationGetResponse{}
	var _, err = at.Get(config, &r)
	if err != nil {
		return nil, at.handleKnownErrors(err)
	}

	return &r.Configuration, nil
}
