package solacesdk

import (
	"fmt"
	"time"
)

type Application struct {
	CreatedTime         time.Time `json:"createdTime"`
	UpdatedTime         time.Time `json:"updatedTime"`
	CreatedBy           string    `json:"createdBy"`
	ChangedBy           string    `json:"changedBy"`
	Id                  string    `json:"id"`
	Name                string    `json:"name"`
	ApplicationType     string    `json:"applicationType"`
	BrokerType          string    `json:"brokerType"`
	ApplicationDomainId string    `json:"applicationDomainId"`
	NumberOfVersions    int       `json:"numberOfVersions"`
	Type                string    `json:"type"`
	CustomAttributes    []struct {
		CustomAttributeDefinitionId   string `json:"customAttributeDefinitionId"`
		CustomAttributeDefinitionName string `json:"customAttributeDefinitionName"`
		Value                         string `json:"value"`
	} `json:"customAttributes"`
}
type ApplicationListResponse struct {
	Applications []Application `json:"data"`
	Meta         Meta          `json:"meta"`
}

type ApplicationListPaginator struct {
	firstPage  bool
	pageNumber int
	count      int
	pageSize   int
	nextPage   int
	totalPages int
	config     *RequestConfig
	client     *Client
}

func (lp *ApplicationListPaginator) HasNoMorePages() bool {
	return lp.nextPage == 0 || lp.pageNumber > lp.totalPages
}

func (lp *ApplicationListPaginator) NextPage() ([]Application, Meta, error) {
	if !lp.firstPage && lp.HasNoMorePages() {
		return nil, Meta{}, fmt.Errorf("no more pages available")
	}
	var r = &ApplicationListResponse{}
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

	return r.Applications, r.Meta, nil
}

func (at *Client) NewApplicationListPaginator(args map[string]string) *ApplicationListPaginator {
	var config = NewRequestConfig("architecture/applications")
	ids, exists := args["ids"]
	if exists && len(ids) > 0 {
		config.Params.Add("ids", ids)
	}
	applicationDomainId, exists := args["applicationDomainId"]
	if exists && len(applicationDomainId) > 0 {
		config.Params.Add("applicationDomainId", applicationDomainId)
	}

	config.Pagination = NewRequestPagination()

	var p = &ApplicationListPaginator{
		firstPage:  true,
		pageNumber: 1,
		pageSize:   defaultPageSize,
		config:     config,
		client:     at,
	}

	return p
}

type ApplicationGetResponse struct {
	Application Application `json:"data"`
	Meta        Meta        `json:"meta"`
}

func (at *Client) GetApplication(applicationId string) (*Application, error) {
	var config = NewRequestConfig(fmt.Sprintf(`architecture/applications/%s`, applicationId))

	var r = &ApplicationGetResponse{}
	var _, err = at.Get(config, &r)
	if err != nil {
		return nil, at.handleKnownErrors(err)
	}

	return &r.Application, nil
}
