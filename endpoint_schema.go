package solace

import (
	"fmt"
	"time"
)

type Schema struct {
	CreatedTime          time.Time `json:"createdTime"`
	UpdatedTime          time.Time `json:"updatedTime"`
	CreatedBy            string    `json:"createdBy"`
	ChangedBy            string    `json:"changedBy"`
	Id                   string    `json:"id"`
	ApplicationDomainId  string    `json:"applicationDomainId"`
	Name                 string    `json:"name"`
	Shared               bool      `json:"shared"`
	SchemaType           string    `json:"schemaType"`
	NumberOfVersions     int       `json:"numberOfVersions"`
	EventVersionRefCount int       `json:"eventVersionRefCount"`
	CustomAttributes     []struct {
		CustomAttributeDefinitionId   string `json:"customAttributeDefinitionId"`
		CustomAttributeDefinitionName string `json:"customAttributeDefinitionName"`
		Value                         string `json:"value"`
	} `json:"customAttributes"`
	Type string `json:"type"`
}

type SchemaListResponse struct {
	Schemas []Schema `json:"data"`
	Meta    Meta     `json:"meta"`
}

type SchemaListPaginator struct {
	firstPage  bool
	pageNumber int
	count      int
	pageSize   int
	nextPage   int
	totalPages int
	config     *RequestConfig
	client     *Client
}

func (lp *SchemaListPaginator) HasNoMorePages() bool {
	return lp.nextPage == 0 || lp.pageNumber > lp.totalPages
}

func (lp *SchemaListPaginator) NextPage() ([]Schema, Meta, error) {
	if !lp.firstPage && lp.HasNoMorePages() {
		return nil, Meta{}, fmt.Errorf("no more pages available")
	}
	var r = &SchemaListResponse{}
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

	return r.Schemas, r.Meta, nil
}

func (at *Client) NewSchemaListPaginator(args map[string]string) *SchemaListPaginator {
	var config = NewRequestConfig("architecture/schemas")
	config.Pagination = NewRequestPagination()
	ids, exists := args["ids"]
	if exists && len(ids) > 0 {
		config.Params.Add("ids", ids)
	}
	applicationDomainIds, exists := args["applicationDomainIds"]
	if exists && len(applicationDomainIds) > 0 {
		config.Params.Add("applicationDomainIds", applicationDomainIds)
	}

	var p = &SchemaListPaginator{
		firstPage:  true,
		pageNumber: 1,
		pageSize:   defaultPageSize,
		config:     config,
		client:     at,
	}

	return p
}

type SchemaGetResponse struct {
	Schema Schema `json:"data"`
	Meta   Meta   `json:"meta"`
}

func (at *Client) GetSchema(schemaId string) (*Schema, error) {
	var config = NewRequestConfig(fmt.Sprintf(`architecture/schemas/%s`, schemaId))

	var r = &SchemaGetResponse{}
	var _, err = at.Get(config, &r)
	if err != nil {
		return nil, at.handleKnownErrors(err)
	}

	return &r.Schema, nil
}
