package solacesdk

import (
	"fmt"
	"time"
)

type Consumer struct {
	CreatedTime          time.Time `json:"createdTime"`
	UpdatedTime          time.Time `json:"updatedTime"`
	CreatedBy            string    `json:"createdBy"`
	ChangedBy            string    `json:"changedBy"`
	Id                   string    `json:"id"`
	Name                 string    `json:"name"`
	ApplicationVersionId string    `json:"applicationVersionId"`
	BrokerType           string    `json:"brokerType"`
	ConsumerType         string    `json:"consumerType"`
	Subscriptions        []struct {
		Id                       string `json:"id"`
		SubscriptionType         string `json:"subscriptionType"`
		Value                    string `json:"value"`
		AttractedEventVersionIds []struct {
			EventVersionId string   `json:"eventVersionId"`
			EventMeshIds   []string `json:"eventMeshIds"`
		} `json:"attractedEventVersionIds"`
	} `json:"subscriptions"`
	Type string `json:"type"`
}

type ConsumerListResponse struct {
	Consumers []Consumer `json:"data"`
	Meta      Meta       `json:"meta"`
}

type ConsumerListPaginator struct {
	firstPage  bool
	pageNumber int
	count      int
	pageSize   int
	nextPage   int
	totalPages int
	config     *RequestConfig
	client     *Client
}

func (lp *ConsumerListPaginator) HasNoMorePages() bool {
	return lp.nextPage == 0 || lp.pageNumber > lp.totalPages
}

func (lp *ConsumerListPaginator) NextPage() ([]Consumer, Meta, error) {
	if !lp.firstPage && lp.HasNoMorePages() {
		return nil, Meta{}, fmt.Errorf("no more pages available")
	}
	var r = &ConsumerListResponse{}
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

	return r.Consumers, r.Meta, nil
}

func (at *Client) NewConsumerListPaginator(args map[string]string) *ConsumerListPaginator {
	var config = NewRequestConfig("architecture/consumers")
	config.Pagination = NewRequestPagination()
	ids, exists := args["ids"]
	if exists && len(ids) > 0 {
		config.Params.Add("ids", ids)
	}
	applicationVersionIds, exists := args["applicationVersionIds"]
	if exists && len(applicationVersionIds) > 0 {
		config.Params.Add("applicationVersionIds", applicationVersionIds)
	}

	var p = &ConsumerListPaginator{
		firstPage:  true,
		pageNumber: 1,
		pageSize:   defaultPageSize,
		config:     config,
		client:     at,
	}

	return p
}

type ConsumerGetResponse struct {
	Consumer Consumer `json:"data"`
	Meta     Meta     `json:"meta"`
}

func (at *Client) GetConsumer(consumerId string) (*Consumer, error) {
	var config = NewRequestConfig(fmt.Sprintf(`architecture/consumers/%s`, consumerId))

	var r = &ConsumerGetResponse{}
	var _, err = at.Get(config, &r)
	if err != nil {
		return nil, at.handleKnownErrors(err)
	}

	return &r.Consumer, nil
}
