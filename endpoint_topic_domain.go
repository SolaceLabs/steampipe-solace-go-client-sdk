package solacesdk

import (
	"fmt"
	"time"
)

type TopicDomain struct {
	CreatedTime         time.Time `json:"createdTime"`
	UpdatedTime         time.Time `json:"updatedTime"`
	CreatedBy           string    `json:"createdBy"`
	ChangedBy           string    `json:"changedBy"`
	Id                  string    `json:"id"`
	ApplicationDomainId string    `json:"applicationDomainId"`
	BrokerType          string    `json:"brokerType"`
	AddressLevels       []struct {
		Name             string `json:"name"`
		AddressLevelType string `json:"addressLevelType"`
		EnumVersionId    string `json:"enumVersionId"`
	} `json:"addressLevels"`
	Type string `json:"type"`
}
type TopicDomainListResponse struct {
	TopicDomains []TopicDomain `json:"data"`
	Meta         Meta          `json:"meta"`
}

type TopicDomainListPaginator struct {
	firstPage  bool
	pageNumber int
	count      int
	pageSize   int
	nextPage   int
	totalPages int
	config     *RequestConfig
	client     *Client
}

func (lp *TopicDomainListPaginator) HasNoMorePages() bool {
	return lp.nextPage == 0 || lp.pageNumber > lp.totalPages
}

func (lp *TopicDomainListPaginator) NextPage() ([]TopicDomain, Meta, error) {
	if !lp.firstPage && lp.HasNoMorePages() {
		return nil, Meta{}, fmt.Errorf("no more pages available")
	}
	var r = &TopicDomainListResponse{}
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

	return r.TopicDomains, r.Meta, nil
}

func (at *Client) NewTopicDomainListPaginator() *TopicDomainListPaginator {
	var config = NewRequestConfig("architecture/topicDomains")
	config.Pagination = NewRequestPagination()

	var p = &TopicDomainListPaginator{
		firstPage:  true,
		pageNumber: 1,
		pageSize:   defaultPageSize,
		config:     config,
		client:     at,
	}

	return p
}

type TopicDomainGetResponse struct {
	TopicDomain TopicDomain `json:"data"`
	Meta        Meta        `json:"meta"`
}

func (at *Client) GetTopicDomain(topicDomainId string) (*TopicDomain, error) {
	var config = NewRequestConfig(fmt.Sprintf(`architecture/topicDomains/%s`, topicDomainId))

	var r = &TopicDomainGetResponse{}
	var _, err = at.Get(config, &r)
	if err != nil {
		return nil, at.handleKnownErrors(err)
	}

	return &r.TopicDomain, nil
}
