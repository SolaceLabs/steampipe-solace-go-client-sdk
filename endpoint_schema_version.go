package solace

import (
	"fmt"
	"time"
)

type SchemaVersion struct {
	CreatedTime                  time.Time `json:"createdTime"`
	UpdatedTime                  time.Time `json:"updatedTime"`
	CreatedBy                    string    `json:"createdBy"`
	ChangedBy                    string    `json:"changedBy"`
	Id                           string    `json:"id"`
	SchemaId                     string    `json:"schemaId"`
	Description                  string    `json:"description"`
	Version                      string    `json:"version"`
	DisplayName                  string    `json:"displayName"`
	Content                      string    `json:"content"`
	ReferencedByEventVersionIds  []string  `json:"referencedByEventVersionIds"`
	ReferencedBySchemaVersionIds []string  `json:"referencedBySchemaVersionIds"`
	SchemaVersionReferences      []struct {
		SchemaVersionId string `json:"schemaVersionId"`
	} `json:"schemaVersionReferences"`
	CustomAttributes []struct {
		CustomAttributeDefinitionId   string `json:"customAttributeDefinitionId"`
		CustomAttributeDefinitionName string `json:"customAttributeDefinitionName"`
		Value                         string `json:"value"`
	} `json:"customAttributes"`
	StateId string `json:"stateId"`
	Type    string `json:"type"`
}

type SchemaVersionListResponse struct {
	SchemaVersions []SchemaVersion `json:"data"`
	Meta           Meta            `json:"meta"`
}

type SchemaVersionListPaginator struct {
	firstPage  bool
	pageNumber int
	count      int
	pageSize   int
	nextPage   int
	totalPages int
	config     *RequestConfig
	client     *Client
}

func (lp *SchemaVersionListPaginator) HasNoMorePages() bool {
	return lp.nextPage == 0 || lp.pageNumber > lp.totalPages
}

func (lp *SchemaVersionListPaginator) NextPage() ([]SchemaVersion, Meta, error) {
	if !lp.firstPage && lp.HasNoMorePages() {
		return nil, Meta{}, fmt.Errorf("no more pages available")
	}
	var r = &SchemaVersionListResponse{}
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

	return r.SchemaVersions, r.Meta, nil
}

func (at *Client) NewSchemaVersionListPaginator(args map[string]string) *SchemaVersionListPaginator {
	var config = NewRequestConfig("architecture/schemaVersions")
	config.Pagination = NewRequestPagination()
	ids, exists := args["ids"]
	if exists && len(ids) > 0 {
		config.Params.Add("ids", ids)
	}
	schemaIds, exists := args["schemaIds"]
	if exists && len(schemaIds) > 0 {
		config.Params.Add("schemaIds", schemaIds)
	}

	var p = &SchemaVersionListPaginator{
		firstPage:  true,
		pageNumber: 1,
		pageSize:   defaultPageSize,
		config:     config,
		client:     at,
	}

	return p
}

type SchemaVersionGetResponse struct {
	SchemaVersion SchemaVersion `json:"data"`
	Meta          Meta          `json:"meta"`
}

func (at *Client) GetSchemaVersion(schemaVersionId string) (*SchemaVersion, error) {
	var config = NewRequestConfig(fmt.Sprintf(`architecture/schemaVersions/%s`, schemaVersionId))

	var r = &SchemaVersionGetResponse{}
	var _, err = at.Get(config, &r)
	if err != nil {
		return nil, at.handleKnownErrors(err)
	}

	return &r.SchemaVersion, nil
}
