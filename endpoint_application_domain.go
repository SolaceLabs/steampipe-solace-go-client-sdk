package solace

import (
	"fmt"
	"time"
)

type ApplicationDomain struct {
	Id                                   string    `json:"id"`
	Name                                 string    `json:"name"`
	Description                          string    `json:"description"`
	UniqueTopicAddressEnforcementEnabled bool      `json:"uniqueTopicAddressEnforcementEnabled"`
	TopicDomainEnforcementEnabled        bool      `json:"topicDomainEnforcementEnabled"`
	CreatedBy                            string    `json:"createdBy"`
	CreatedTime                          time.Time `json:"createdTime"`
	ChangedBy                            string    `json:"changedBy"`
	UpdatedTime                          time.Time `json:"updatedTime"`
	Stats                                struct {
		SchemaCount          int `json:"schemaCount"`
		EventCount           int `json:"eventCount"`
		ApplicationCount     int `json:"applicationCount"`
		EnumCount            int `json:"enumCount"`
		EventApiCount        int `json:"eventApiCount"`
		EventApiProductCount int `json:"eventApiProductCount"`
	} `json:"stats"`
	CustomAttributes []struct {
		CustomAttributeDefinitionId   string `json:"customAttributeDefinitionId"`
		CustomAttributeDefinitionName string `json:"customAttributeDefinitionName"`
		Value                         string `json:"value"`
	} `json:"customAttributes"`
}

type ApplicationDomainListResponse struct {
	ApplicationDomains []ApplicationDomain `json:"data"`
	Meta               Meta                `json:"meta"`
}

type ApplicationDomainListPaginator struct {
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

func (lp *ApplicationDomainListPaginator) HasNoMorePages() bool {
	return lp.nextPage == 0 || lp.pageNumber > lp.totalPages
}

func (lp *ApplicationDomainListPaginator) NextPage() ([]ApplicationDomain, Meta, error) {
	if !lp.firstPage && lp.HasNoMorePages() {
		return nil, Meta{}, fmt.Errorf("no more pages available")
	}
	var r = &ApplicationDomainListResponse{}
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

	return r.ApplicationDomains, r.Meta, nil
}

// map[string]string{}

func (at *Client) NewApplicationDomainListPaginator(args map[string]string) *ApplicationDomainListPaginator {
	var config = NewRequestConfig("architecture/applicationDomains")
	config.Params.Add("include", "stats")
	value, exists := args["ids"]
	if exists && len(value) > 0 {
		config.Params.Add("ids", value)
	}
	value, exists = args["id"]
	if exists && len(value) > 0 {
		config.Params.Add("id", value)
	}

	config.Pagination = NewRequestPagination()

	var p = &ApplicationDomainListPaginator{
		firstPage:  true,
		pageNumber: 1,
		pageSize:   defaultPageSize,
		config:     config,
		client:     at,
	}

	return p
}

type ApplicationDomainGetResponse struct {
	ApplicationDomain ApplicationDomain `json:"data"`
	Meta              Meta              `json:"meta"`
}

func (at *Client) GetApplicationDomain(domainId string) (*ApplicationDomain, error) {
	var config = NewRequestConfig(fmt.Sprintf(`architecture/applicationDomains/%s`, domainId))

	var r = &ApplicationDomainGetResponse{}
	var _, err = at.Get(config, &r)
	if err != nil {
		return nil, at.handleKnownErrors(err)
	}

	return &r.ApplicationDomain, nil
}
