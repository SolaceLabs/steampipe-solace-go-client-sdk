package solacesdk

import (
	"fmt"
	"time"
)

type ConfigurationType struct {
	CreatedTime           time.Time `json:"createdTime"`
	UpdatedTime           time.Time `json:"updatedTime"`
	CreatedBy             string    `json:"createdBy"`
	ChangedBy             string    `json:"changedBy"`
	Id                    string    `json:"id"`
	Name                  string    `json:"name"`
	BrokerType            string    `json:"brokerType"`
	AssociatedEntityTypes []string  `json:"associatedEntityTypes"`
	ValueSchema           struct {
	} `json:"valueSchema"`
	Type string `json:"type"`
}

type ConfigurationTypeListResponse struct {
	ConfigurationTypes []ConfigurationType `json:"data"`
	Meta               Meta                `json:"meta"`
}

type ConfigurationTypeListPaginator struct {
	firstPage  bool
	pageNumber int
	count      int
	pageSize   int
	nextPage   int
	totalPages int
	config     *RequestConfig
	client     *Client
}

func (lp *ConfigurationTypeListPaginator) HasNoMorePages() bool {
	return lp.nextPage == 0 || lp.pageNumber > lp.totalPages
}

func (lp *ConfigurationTypeListPaginator) NextPage() ([]ConfigurationType, Meta, error) {
	if !lp.firstPage && lp.HasNoMorePages() {
		return nil, Meta{}, fmt.Errorf("no more pages available")
	}
	var r = &ConfigurationTypeListResponse{}
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

	return r.ConfigurationTypes, r.Meta, nil
}

func (at *Client) NewConfigurationTypeListPaginator() *ConfigurationTypeListPaginator {
	var config = NewRequestConfig("architecture/configurationTypes")
	config.Pagination = NewRequestPagination()

	var p = &ConfigurationTypeListPaginator{
		firstPage:  true,
		pageNumber: 1,
		pageSize:   defaultPageSize,
		config:     config,
		client:     at,
	}

	return p
}

type ConfigurationTypeGetResponse struct {
	ConfigurationType ConfigurationType `json:"data"`
	Meta              Meta              `json:"meta"`
}

func (at *Client) GetConfigurationType(configurationTypeId string) (*ConfigurationType, error) {
	var config = NewRequestConfig(fmt.Sprintf(`architecture/configurationTypes/%s`, configurationTypeId))

	var r = &ConfigurationTypeGetResponse{}
	var _, err = at.Get(config, &r)
	if err != nil {
		return nil, at.handleKnownErrors(err)
	}

	return &r.ConfigurationType, nil
}
