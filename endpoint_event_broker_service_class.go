package solacesdk

import (
	"fmt"
)

type EventBrokerServiceClass struct {
	Id                      string `json:"id"`
	Type                    string `json:"type"`
	Name                    string `json:"name"`
	VpnConnections          int    `json:"vpnConnections"`
	BrokerScalingTier       string `json:"brokerScalingTier"`
	VpnMaxSpoolSize         int    `json:"vpnMaxSpoolSize"`
	MaxNumberVpns           int    `json:"maxNumberVpns"`
	HighAvailabilityCapable bool   `json:"highAvailabilityCapable"`
}

type EventBrokerServiceClassListResponse struct {
	EventBrokerServiceClasses []EventBrokerServiceClass `json:"data"`
	Meta                      Meta                      `json:"meta"`
}

type EventBrokerServiceClassListPaginator struct {
	firstPage  bool
	pageNumber int
	count      int
	pageSize   int
	nextPage   int
	totalPages int
	paginate   bool
	config     *RequestConfig
	client     *Client
}

func (lp *EventBrokerServiceClassListPaginator) HasMorePages() bool {
	return lp.nextPage > 0 || lp.firstPage
}

func (lp *EventBrokerServiceClassListPaginator) NextPage() ([]EventBrokerServiceClass, Meta, error) {
	if !lp.HasMorePages() {
		return nil, Meta{}, fmt.Errorf("no more pages available")
	}

	var r = &EventBrokerServiceClassListResponse{}
	var _, err = lp.client.Get(lp.config, r)
	if err != nil {
		return nil, Meta{}, lp.client.handleKnownErrors(err)
	}

	lp.firstPage = false
	lp.pageNumber += 1
	lp.config.Pagination.NextPage(lp.pageSize)

	return r.EventBrokerServiceClasses, r.Meta, nil
}

func (at *Client) NewEventBrokerServiceClassListPaginator() *EventBrokerServiceClassListPaginator {
	var config = NewRequestConfig("missionControl/serviceClasses")

	config.Pagination = NewRequestNoPagination()

	var p = &EventBrokerServiceClassListPaginator{
		firstPage:  true,
		pageNumber: 1,
		pageSize:   defaultPageSize,
		config:     config,
		client:     at,
	}

	return p
}

type EventBrokerServiceClassGetResponse struct {
	EventBrokerServiceClass EventBrokerServiceClass `json:"data"`
}

func (at *Client) GetEventBrokerServiceClass(eventBrokerServiceClassId string) (*EventBrokerServiceClass, error) {
	var config = NewRequestConfig(fmt.Sprintf(`missionControl/serviceClasses/%s`, eventBrokerServiceClassId))

	var result = &EventBrokerServiceClassGetResponse{}
	var _, err = at.Get(config, &result)
	if err != nil {
		return nil, at.handleKnownErrors(err)
	}

	return &result.EventBrokerServiceClass, nil
}
