package solace

import (
	"fmt"
	"time"
)

type EnumVersion struct {
	CreatedTime time.Time `json:"createdTime"`
	UpdatedTime time.Time `json:"updatedTime"`
	CreatedBy   string    `json:"createdBy"`
	ChangedBy   string    `json:"changedBy"`
	Id          string    `json:"id"`
	EnumId      string    `json:"enumId"`
	Description string    `json:"description"`
	Version     string    `json:"version"`
	DisplayName string    `json:"displayName"`
	Values      []struct {
		CreatedTime   time.Time `json:"createdTime"`
		UpdatedTime   time.Time `json:"updatedTime"`
		CreatedBy     string    `json:"createdBy"`
		ChangedBy     string    `json:"changedBy"`
		Id            string    `json:"id"`
		EnumVersionId string    `json:"enumVersionId"`
		Value         string    `json:"value"`
		Label         string    `json:"label"`
		Type          string    `json:"type"`
	} `json:"values"`
	ReferencedByEventVersionIds []string `json:"referencedByEventVersionIds"`
	ReferencedByTopicDomainIds  []string `json:"referencedByTopicDomainIds"`
	StateId                     string   `json:"stateId"`
	CustomAttributes            []struct {
		CustomAttributeDefinitionId   string `json:"customAttributeDefinitionId"`
		CustomAttributeDefinitionName string `json:"customAttributeDefinitionName"`
		Value                         string `json:"value"`
	} `json:"customAttributes"`
	Type string `json:"type"`
}

type EnumVersionListResponse struct {
	EnumVersions []EnumVersion `json:"data"`
	Meta         Meta          `json:"meta"`
}

type EnumVersionListPaginator struct {
	firstPage  bool
	pageNumber int
	count      int
	pageSize   int
	nextPage   int
	totalPages int
	config     *RequestConfig
	client     *Client
}

func (lp *EnumVersionListPaginator) HasNoMorePages() bool {
	return lp.nextPage == 0 || lp.pageNumber > lp.totalPages
}

func (lp *EnumVersionListPaginator) NextPage() ([]EnumVersion, Meta, error) {
	if !lp.firstPage && lp.HasNoMorePages() {
		return nil, Meta{}, fmt.Errorf("no more pages available")
	}
	var r = &EnumVersionListResponse{}
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

	return r.EnumVersions, r.Meta, nil
}

func (at *Client) NewEnumVersionListPaginator(args map[string]string) *EnumVersionListPaginator {
	var config = NewRequestConfig("architecture/enumVersions")
	config.Pagination = NewRequestPagination()
	ids, exists := args["ids"]
	if exists && len(ids) > 0 {
		config.Params.Add("ids", ids)
	}
	enumIds, exists := args["enumIds"]
	if exists && len(enumIds) > 0 {
		config.Params.Add("enumIds", enumIds)
	}

	var p = &EnumVersionListPaginator{
		firstPage:  true,
		pageNumber: 1,
		pageSize:   defaultPageSize,
		config:     config,
		client:     at,
	}

	return p
}

type EnumVersionGetResponse struct {
	EnumVersion EnumVersion `json:"data"`
	Meta        Meta        `json:"meta"`
}

func (at *Client) GetEnumVersion(enumVersionId string) (*EnumVersion, error) {
	var config = NewRequestConfig(fmt.Sprintf(`architecture/enumVersions/%s`, enumVersionId))

	var r = &EnumVersionGetResponse{}
	var _, err = at.Get(config, &r)
	if err != nil {
		return nil, at.handleKnownErrors(err)
	}

	return &r.EnumVersion, nil
}
