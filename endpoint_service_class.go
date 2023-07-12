package solacesdk

import (
	"fmt"
)

type ServiceClass struct {
	Id                      string `json:"id"`
	Type                    string `json:"type"`
	Name                    string `json:"name"`
	VpnConnections          int    `json:"vpnConnections"`
	BrokerScalingTier       string `json:"brokerScalingTier"`
	VpnMaxSpoolSize         int    `json:"vpnMaxSpoolSize"`
	MaxNumberVpns           int    `json:"maxNumberVpns"`
	HighAvailabilityCapable bool   `json:"highAvailabilityCapable"`
}
type ServiceClassListResponse struct {
	ServiceClasses []ServiceClass `json:"data"`
	Meta           Meta           `json:"meta"`
}

type ServiceClassListPaginator struct {
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

func (lp *ServiceClassListPaginator) HasNoMorePages() bool {
	return lp.nextPage == 0 || lp.pageNumber > lp.totalPages
}

func (lp *ServiceClassListPaginator) NextPage() ([]ServiceClass, Meta, error) {
	if !lp.firstPage && lp.HasNoMorePages() {
		return nil, Meta{}, fmt.Errorf("no more pages available")
	}
	var r = &ServiceClassListResponse{}
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

	return r.ServiceClasses, r.Meta, nil
}

func (at *Client) NewServiceClassListPaginator() *ServiceClassListPaginator {
	var config = NewRequestConfig("missionControl/serviceClasses")
	config.Pagination = NewRequestNoPagination()

	var p = &ServiceClassListPaginator{
		firstPage:  true,
		pageNumber: 1,
		pageSize:   defaultPageSize,
		config:     config,
		client:     at,
	}

	return p
}

type ServiceClassGetResponse struct {
	ServiceClass ServiceClass `json:"data"`
	Meta         Meta         `json:"meta"`
}

func (at *Client) GetServiceClass(serviceClassId string) (*ServiceClass, error) {
	var config = NewRequestConfig(fmt.Sprintf(`missionControl/serviceClasses/%s`, serviceClassId))

	var r = &ServiceClassGetResponse{}
	var _, err = at.Get(config, &r)
	if err != nil {
		return nil, at.handleKnownErrors(err)
	}

	return &r.ServiceClass, nil
}
