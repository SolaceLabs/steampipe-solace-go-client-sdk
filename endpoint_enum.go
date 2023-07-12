package solace

import (
	"fmt"
	"time"
)

type Enum struct {
	CreatedTime          time.Time `json:"createdTime"`
	UpdatedTime          time.Time `json:"updatedTime"`
	CreatedBy            string    `json:"createdBy"`
	ChangedBy            string    `json:"changedBy"`
	Id                   string    `json:"id"`
	ApplicationDomainId  string    `json:"applicationDomainId"`
	Name                 string    `json:"name"`
	Shared               bool      `json:"shared"`
	NumberOfVersions     int       `json:"numberOfVersions"`
	EventVersionRefCount int       `json:"eventVersionRefCount"`
	CustomAttributes     []struct {
		CustomAttributeDefinitionId   string `json:"customAttributeDefinitionId"`
		CustomAttributeDefinitionName string `json:"customAttributeDefinitionName"`
		Value                         string `json:"value"`
	} `json:"customAttributes"`
	Type string `json:"type"`
}

type EnumListResponse struct {
	Enums []Enum `json:"data"`
	Meta  Meta   `json:"meta"`
}

type EnumListPaginator struct {
	firstPage  bool
	pageNumber int
	count      int
	pageSize   int
	nextPage   int
	totalPages int
	config     *RequestConfig
	client     *Client
}

func (lp *EnumListPaginator) HasNoMorePages() bool {
	return lp.nextPage == 0 || lp.pageNumber > lp.totalPages
}

func (lp *EnumListPaginator) NextPage() ([]Enum, Meta, error) {
	if !lp.firstPage && lp.HasNoMorePages() {
		return nil, Meta{}, fmt.Errorf("no more pages available")
	}
	var r = &EnumListResponse{}
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

	return r.Enums, r.Meta, nil
}

func (at *Client) NewEnumListPaginator(args map[string]string) *EnumListPaginator {
	var config = NewRequestConfig("architecture/enums")
	config.Pagination = NewRequestPagination()
	ids, exists := args["ids"]
	if exists && len(ids) > 0 {
		config.Params.Add("ids", ids)
	}
	applicationDomainIds, exists := args["applicationDomainIds"]
	if exists && len(applicationDomainIds) > 0 {
		config.Params.Add("applicationDomainIds", applicationDomainIds)
	}

	var p = &EnumListPaginator{
		firstPage:  true,
		pageNumber: 1,
		pageSize:   defaultPageSize,
		config:     config,
		client:     at,
	}

	return p
}

type EnumGetResponse struct {
	Enum Enum `json:"data"`
	Meta Meta `json:"meta"`
}

func (at *Client) GetEnum(enumId string) (*Enum, error) {
	var config = NewRequestConfig(fmt.Sprintf(`architecture/enums/%s`, enumId))

	var r = &EnumGetResponse{}
	var _, err = at.Get(config, &r)
	if err != nil {
		return nil, at.handleKnownErrors(err)
	}

	return &r.Enum, nil
}
