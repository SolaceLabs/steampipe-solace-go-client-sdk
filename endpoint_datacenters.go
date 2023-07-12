package solace

import (
	"fmt"
)

type Datacenter struct {
	Id                           string   `json:"id"`
	Type                         string   `json:"type"`
	Name                         string   `json:"name"`
	DatacenterType               string   `json:"datacenterType"`
	Provider                     string   `json:"provider"`
	OperState                    string   `json:"operState"`
	CreatedBy                    string   `json:"createdBy"`
	CreatedTime                  string   `json:"createdTime"`
	UpdatedBy                    string   `json:"updatedBy"`
	UpdatedTime                  string   `json:"updatedTime"`
	Available                    bool     `json:"available"`
	SupportedServiceClasses      []string `json:"supportedServiceClasses"`
	CloudAgentVersion            string   `json:"cloudAgentVersion"`
	K8sServiceType               string   `json:"k8sServiceType"`
	NumSupportedPrivateEndpoints int      `json:"numSupportedPrivateEndpoints"`
	NumSupportedPublicEndpoints  int      `json:"numSupportedPublicEndpoints"`
	OrganizationId               string   `json:"organizationId"`
	Location                     struct {
		Longitude string `json:"longitude"`
		Latitude  string `json:"latitude"`
		Continent string `json:"continent"`
	} `json:"location"`
	SpoolScaleUpCapabilityInfo struct {
		SpoolScaleUpCapabilityState string `json:"spoolScaleUpCapabilityState"`
		SpoolScaleUpTestTimestamp   string `json:"spoolScaleUpTestTimestamp"`
		SpoolScaleUpTestMessage     string `json:"spoolScaleUpTestMessage"`
	} `json:"spoolScaleUpCapabilityInfo"`
}

type DatacenterListResponse struct {
	Datacenters []Datacenter `json:"data"`
	Meta        Meta         `json:"meta"`
}

type DatacenterListPaginator struct {
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

func (lp *DatacenterListPaginator) HasMorePages() bool {
	return lp.nextPage > 0 || lp.firstPage
}

func (lp *DatacenterListPaginator) NextPage() ([]Datacenter, Meta, error) {
	if !lp.HasMorePages() {
		return nil, Meta{}, fmt.Errorf("no more pages available")
	}

	var r = &DatacenterListResponse{}
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

	return r.Datacenters, r.Meta, nil
}

func (at *Client) NewDatacenterListPaginator() *DatacenterListPaginator {
	var config = NewRequestConfig("missionControl/datacenters")

	config.Pagination = NewRequestNoPagination()

	var p = &DatacenterListPaginator{
		firstPage:  true,
		pageNumber: 1,
		pageSize:   defaultPageSize,
		config:     config,
		client:     at,
	}

	return p
}

type DatacenterGetResponse struct {
	Datacenter Datacenter `json:"data"`
}

func (at *Client) GetDatacenter(datacenterId string) (*Datacenter, error) {
	var config = NewRequestConfig(fmt.Sprintf(`missionControl/datacenters/%s`, datacenterId))

	var result = &DatacenterGetResponse{}
	var _, err = at.Get(config, &result)
	if err != nil {
		return nil, at.handleKnownErrors(err)
	}

	return &result.Datacenter, nil
}
