package solacesdk

import (
	"fmt"
	"strings"
	"time"
)

type ApplicationVersion struct {
	CreatedTime                               time.Time `json:"createdTime"`
	UpdatedTime                               time.Time `json:"updatedTime"`
	CreatedBy                                 string    `json:"createdBy"`
	ChangedBy                                 string    `json:"changedBy"`
	Id                                        string    `json:"id"`
	ApplicationId                             string    `json:"applicationId"`
	Description                               string    `json:"description"`
	Version                                   string    `json:"version"`
	DisplayName                               string    `json:"displayName"`
	DeclaredProducedEventVersionIds           []string  `json:"declaredProducedEventVersionIds"`
	DeclaredProducedEventVersionIdsAsString   string    `json:"declaredProducedEventVersionIdsAsString"`
	DeclaredConsumedEventVersionIds           []string  `json:"declaredConsumedEventVersionIds"`
	DeclaredConsumedEventVersionIdsAsString   string    `json:"DeclaredConsumedEventVersionIdsAsString"`
	DeclaredEventAPIProductVersionIds         []string  `json:"declaredEventApiProductVersionIds"`
	DeclaredEventAPIProductVersionIdsAsString string    `json:"DeclaredEventAPIProductVersionIdsAsString"`
	StateId                                   string    `json:"stateId"`
	Consumers                                 []struct {
		CreatedTime          time.Time `json:"createdTime"`
		UpdatedTime          time.Time `json:"updatedTime"`
		CreatedBy            string    `json:"createdBy"`
		ChangedBy            string    `json:"changedBy"`
		Id                   string    `json:"id"`
		Name                 string    `json:"name"`
		ApplicationVersionId string    `json:"applicationVersionId"`
		BrokerType           string    `json:"brokerType"`
		ConsumerType         string    `json:"consumerType"`
		Subscriptions        []struct {
			Id                       string `json:"id"`
			SubscriptionType         string `json:"subscriptionType"`
			Value                    string `json:"value"`
			AttractedEventVersionIds []struct {
				EventVersionId string   `json:"eventVersionId"`
				EventMeshIds   []string `json:"eventMeshIds"`
			} `json:"attractedEventVersionIds"`
		} `json:"subscriptions"`
		Type string `json:"type"`
	} `json:"consumers"`
	CustomAttributes []struct {
		CustomAttributeDefinitionId   string `json:"customAttributeDefinitionId"`
		CustomAttributeDefinitionName string `json:"customAttributeDefinitionName"`
		Value                         string `json:"value"`
	} `json:"customAttributes"`
	MessagingServiceIds []string `json:"messagingServiceIds"`
	Type                string   `json:"type"`
}

type ApplicationVersionListResponse struct {
	ApplicationVersions []ApplicationVersion `json:"data"`
	Meta                Meta                 `json:"meta"`
}

type ApplicationVersionListPaginator struct {
	firstPage  bool
	pageNumber int
	count      int
	pageSize   int
	nextPage   int
	totalPages int
	config     *RequestConfig
	client     *Client
}

func buildApplicationVersionCSV(ids []string) string {
	// produced := r.ApplicationVersion.DeclaredProducedEventVersionIds
	// st := strings.Trim(strings.Join(strings.Fields(fmt.Sprint(produced)), ", "), "[]")
	// r.ApplicationVersion.DeclaredProducedEventVersionIdsAsString = st
	for index := range ids {
		ids[index] = fmt.Sprintf("'%s'", ids[index])
	}
	st := strings.Trim(strings.Join(strings.Fields(fmt.Sprint(ids)), ", "), "[]")
	return st
}

func (lp *ApplicationVersionListPaginator) HasNoMorePages() bool {
	return lp.nextPage == 0 || lp.pageNumber > lp.totalPages
}

func (lp *ApplicationVersionListPaginator) NextPage() ([]ApplicationVersion, Meta, error) {
	if !lp.firstPage && lp.HasNoMorePages() {
		return nil, Meta{}, fmt.Errorf("no more pages available")
	}
	var r = &ApplicationVersionListResponse{}
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

	for index := range r.ApplicationVersions {
		r.ApplicationVersions[index].DeclaredProducedEventVersionIdsAsString = buildApplicationVersionCSV(r.ApplicationVersions[index].DeclaredProducedEventVersionIds)
		r.ApplicationVersions[index].DeclaredConsumedEventVersionIdsAsString = buildApplicationVersionCSV(r.ApplicationVersions[index].DeclaredConsumedEventVersionIds)
		r.ApplicationVersions[index].DeclaredEventAPIProductVersionIdsAsString = buildApplicationVersionCSV(r.ApplicationVersions[index].DeclaredEventAPIProductVersionIds)
	}

	return r.ApplicationVersions, r.Meta, nil
}

func (at *Client) NewApplicationVersionListPaginator(args map[string]string) *ApplicationVersionListPaginator {
	var config = NewRequestConfig("architecture/applicationVersions")
	config.Pagination = NewRequestPagination()
	ids, exists := args["ids"]
	if exists && len(ids) > 0 {
		config.Params.Add("ids", ids)
	}
	applicationIds, exists := args["applicationIds"]
	if exists && len(applicationIds) > 0 {
		config.Params.Add("applicationIds", applicationIds)
	}

	var p = &ApplicationVersionListPaginator{
		firstPage:  true,
		pageNumber: 1,
		pageSize:   defaultPageSize,
		config:     config,
		client:     at,
	}

	return p
}

type ApplicationVersionGetResponse struct {
	ApplicationVersion ApplicationVersion `json:"data"`
	Meta               Meta               `json:"meta"`
}

func (at *Client) GetApplicationVersion(applicationVersionId string) (*ApplicationVersion, error) {
	var config = NewRequestConfig(fmt.Sprintf(`architecture/applicationVersions/%s`, applicationVersionId))

	var r = &ApplicationVersionGetResponse{}
	var _, err = at.Get(config, &r)
	if err != nil {
		return nil, at.handleKnownErrors(err)
	}

	r.ApplicationVersion.DeclaredProducedEventVersionIdsAsString = buildApplicationVersionCSV(r.ApplicationVersion.DeclaredProducedEventVersionIds)
	r.ApplicationVersion.DeclaredConsumedEventVersionIdsAsString = buildApplicationVersionCSV(r.ApplicationVersion.DeclaredConsumedEventVersionIds)
	r.ApplicationVersion.DeclaredEventAPIProductVersionIdsAsString = buildApplicationVersionCSV(r.ApplicationVersion.DeclaredEventAPIProductVersionIds)
	return &r.ApplicationVersion, nil
}
