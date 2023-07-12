package solacesdk

import (
	"fmt"
	"time"
)

type EventApiProductVersion struct {
	CreatedTime       time.Time `json:"createdTime"`
	UpdatedTime       time.Time `json:"updatedTime"`
	CreatedBy         string    `json:"createdBy"`
	ChangedBy         string    `json:"changedBy"`
	Id                string    `json:"id"`
	EventApiProductId string    `json:"eventApiProductId"`
	Description       string    `json:"description"`
	Version           string    `json:"version"`
	Summary           string    `json:"summary"`
	CustomAttributes  []struct {
		CustomAttributeDefinitionId   string `json:"customAttributeDefinitionId"`
		CustomAttributeDefinitionName string `json:"customAttributeDefinitionName"`
		Value                         string `json:"value"`
	} `json:"customAttributes"`
	DisplayName                  string   `json:"displayName"`
	EventApiVersionIds           []string `json:"eventApiVersionIds"`
	StateId                      string   `json:"stateId"`
	EventApiProductRegistrations []struct {
		CreatedTime              time.Time `json:"createdTime"`
		UpdatedTime              time.Time `json:"updatedTime"`
		CreatedBy                string    `json:"createdBy"`
		ChangedBy                string    `json:"changedBy"`
		Id                       string    `json:"id"`
		ApplicationDomainId      string    `json:"applicationDomainId"`
		RegistrationId           string    `json:"registrationId"`
		AccessRequestId          string    `json:"accessRequestId"`
		EventApiProductVersionId string    `json:"eventApiProductVersionId"`
		PlanId                   string    `json:"planId"`
		State                    string    `json:"state"`
		Type                     string    `json:"type"`
		CustomAttributes         struct {
			AdditionalProp string `json:"additionalProp"`
		} `json:"customAttributes"`
	} `json:"eventApiProductRegistrations"`
	Plans []struct {
		Id                         string `json:"id"`
		Name                       string `json:"name"`
		SolaceClassOfServicePolicy struct {
			Id                  string `json:"id"`
			Type                string `json:"type"`
			MessageDeliveryMode string `json:"messageDeliveryMode"`
			AccessType          string `json:"accessType"`
			MaximumTimeToLive   int    `json:"maximumTimeToLive"`
			QueueType           string `json:"queueType"`
			MaxMsgSpoolUsage    int    `json:"maxMsgSpoolUsage"`
		} `json:"solaceClassOfServicePolicy"`
		Type string `json:"type"`
	} `json:"plans"`
	SolaceMessagingServices []struct {
		Id                            string   `json:"id"`
		MessagingServiceId            string   `json:"messagingServiceId"`
		MessagingServiceName          string   `json:"messagingServiceName"`
		SupportedProtocols            []string `json:"supportedProtocols"`
		EnvironmentId                 string   `json:"environmentId"`
		EnvironmentName               string   `json:"environmentName"`
		EventMeshId                   string   `json:"eventMeshId"`
		EventMeshName                 string   `json:"eventMeshName"`
		Type                          string   `json:"type"`
		SolaceCloudMessagingServiceId string   `json:"solaceCloudMessagingServiceId"`
	} `json:"solaceMessagingServices"`
	Filters []struct {
		EventVersionId string `json:"eventVersionId"`
		TopicFilters   []struct {
			CreatedTime     time.Time `json:"createdTime"`
			UpdatedTime     time.Time `json:"updatedTime"`
			CreatedBy       string    `json:"createdBy"`
			ChangedBy       string    `json:"changedBy"`
			Name            string    `json:"name"`
			FilterValue     string    `json:"filterValue"`
			EventVersionIds []string  `json:"eventVersionIds"`
			Type            string    `json:"type"`
		} `json:"topicFilters"`
		Id   string `json:"id"`
		Type string `json:"type"`
	} `json:"filters"`
	ApprovalType  string    `json:"approvalType"`
	PublishState  string    `json:"publishState"`
	PublishedTime time.Time `json:"publishedTime"`
	Type          string    `json:"type"`
}

type EventApiProductVersionListResponse struct {
	EventApiProductVersions []EventApiProductVersion `json:"data"`
	Meta                    Meta                     `json:"meta"`
}

type EventApiProductVersionListPaginator struct {
	firstPage  bool
	pageNumber int
	count      int
	pageSize   int
	nextPage   int
	totalPages int
	config     *RequestConfig
	client     *Client
}

func (lp *EventApiProductVersionListPaginator) HasNoMorePages() bool {
	return lp.nextPage == 0 || lp.pageNumber > lp.totalPages
}

func (lp *EventApiProductVersionListPaginator) NextPage() ([]EventApiProductVersion, Meta, error) {
	if !lp.firstPage && lp.HasNoMorePages() {
		return nil, Meta{}, fmt.Errorf("no more pages available")
	}
	var r = &EventApiProductVersionListResponse{}
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

	return r.EventApiProductVersions, r.Meta, nil
}

func (at *Client) NewEventApiProductVersionListPaginator(args map[string]string) *EventApiProductVersionListPaginator {
	var config = NewRequestConfig("architecture/eventApiProductVersions")
	config.Pagination = NewRequestPagination()
	ids, exists := args["ids"]
	if exists && len(ids) > 0 {
		config.Params.Add("ids", ids)
	}
	eventApiProductIds, exists := args["eventApiProductIds"]
	if exists && len(eventApiProductIds) > 0 {
		config.Params.Add("eventApiProductIds", eventApiProductIds)
	}

	var p = &EventApiProductVersionListPaginator{
		firstPage:  true,
		pageNumber: 1,
		pageSize:   defaultPageSize,
		config:     config,
		client:     at,
	}

	return p
}

type EventApiProductVersionGetResponse struct {
	EventApiProductVersion EventApiProductVersion `json:"data"`
	Meta                   Meta                   `json:"meta"`
}

func (at *Client) GetEventApiProductVersion(eventApiProductVersionId string) (*EventApiProductVersion, error) {
	var config = NewRequestConfig(fmt.Sprintf(`architecture/eventApiProductVersions/%s`, eventApiProductVersionId))

	var r = &EventApiProductVersionGetResponse{}
	var _, err = at.Get(config, &r)
	if err != nil {
		return nil, at.handleKnownErrors(err)
	}

	return &r.EventApiProductVersion, nil
}
