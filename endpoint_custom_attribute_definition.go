package solace

import (
	"fmt"
	"time"
)

type CustomAttributeDefinition struct {
	CreatedTime           time.Time `json:"createdTime"`
	UpdatedTime           time.Time `json:"updatedTime"`
	CreatedBy             string    `json:"createdBy"`
	ChangedBy             string    `json:"changedBy"`
	Id                    string    `json:"id"`
	Name                  string    `json:"name"`
	ValueType             string    `json:"valueType"`
	Scope                 string    `json:"scope"`
	AssociatedEntityTypes []string  `json:"associatedEntityTypes"`
	AssociatedEntities    []struct {
		EntityType           string   `json:"entityType"`
		ApplicationDomainIds []string `json:"applicationDomainIds"`
	} `json:"associatedEntities"`
	Type string `json:"type"`
}

type CustomAttributeDefinitionListResponse struct {
	CustomAttributeDefinitions []CustomAttributeDefinition `json:"data"`
	Meta                       Meta                        `json:"meta"`
}

type CustomAttributeDefinitionListPaginator struct {
	firstPage  bool
	pageNumber int
	count      int
	pageSize   int
	nextPage   int
	totalPages int
	config     *RequestConfig
	client     *Client
}

func (lp *CustomAttributeDefinitionListPaginator) HasNoMorePages() bool {
	return lp.nextPage == 0 || lp.pageNumber > lp.totalPages
}

func (lp *CustomAttributeDefinitionListPaginator) NextPage() ([]CustomAttributeDefinition, Meta, error) {
	if !lp.firstPage && lp.HasNoMorePages() {
		return nil, Meta{}, fmt.Errorf("no more pages available")
	}
	var r = &CustomAttributeDefinitionListResponse{}
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

	return r.CustomAttributeDefinitions, r.Meta, nil
}

func (at *Client) NewCustomAttributeDefinitionListPaginator() *CustomAttributeDefinitionListPaginator {
	var config = NewRequestConfig("architecture/customAttributeDefinitions")
	config.Pagination = NewRequestPagination()

	var p = &CustomAttributeDefinitionListPaginator{
		firstPage:  true,
		pageNumber: 1,
		pageSize:   defaultPageSize,
		config:     config,
		client:     at,
	}

	return p
}

type CustomAttributeDefinitionGetResponse struct {
	CustomAttributeDefinition CustomAttributeDefinition `json:"data"`
	Meta                      Meta                      `json:"meta"`
}

func (at *Client) GetCustomAttributeDefinition(customAttributeDefinitionId string) (*CustomAttributeDefinition, error) {
	var config = NewRequestConfig(fmt.Sprintf(`architecture/customAttributeDefinitions/%s`, customAttributeDefinitionId))

	var r = &CustomAttributeDefinitionGetResponse{}
	var _, err = at.Get(config, &r)
	if err != nil {
		return nil, at.handleKnownErrors(err)
	}

	return &r.CustomAttributeDefinition, nil
}
