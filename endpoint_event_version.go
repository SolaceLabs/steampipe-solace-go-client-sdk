package solace

import (
	"fmt"
	"strings"
	"time"
)

type EventVersion struct {
	CreatedTime                                    time.Time `json:"createdTime"`
	UpdatedTime                                    time.Time `json:"updatedTime"`
	CreatedBy                                      string    `json:"createdBy"`
	ChangedBy                                      string    `json:"changedBy"`
	Id                                             string    `json:"id"`
	EventId                                        string    `json:"eventId"`
	Description                                    string    `json:"description"`
	Version                                        string    `json:"version"`
	DisplayName                                    string    `json:"displayName"`
	DeclaredProducingApplicationVersionIds         []string  `json:"declaredProducingApplicationVersionIds"`
	DeclaredProducingApplicationVersionIdsAsString string    `json:"declaredProducingApplicationVersionIdsAsString"`
	DeclaredConsumingApplicationVersionIds         []string  `json:"declaredConsumingApplicationVersionIds"`
	DeclaredConsumingApplicationVersionIdsAsString string    `json:"declaredConsumingApplicationVersionIdsAsString"`
	ProducingEventApiVersionIds                    []string  `json:"producingEventApiVersionIds"`
	ProducingEventApiVersionIdsAsString            string    `json:"producingEventApiVersionIdsAsString"`
	ConsumingEventApiVersionIds                    []string  `json:"consumingEventApiVersionIds"`
	ConsumingEventApiVersionIdsAsString            string    `json:"consumingEventApiVersionIdsAsString"`
	AttractingApplicationVersionIds                []struct {
		ApplicationVersionId string   `json:"applicationVersionId"`
		EventMeshIds         []string `json:"eventMeshIds"`
	} `json:"attractingApplicationVersionIds"`
	SchemaVersionId     string `json:"schemaVersionId"`
	SchemaPrimitiveType string `json:"schemaPrimitiveType"`
	DeliveryDescriptor  struct {
		CreatedTime time.Time `json:"createdTime"`
		UpdatedTime time.Time `json:"updatedTime"`
		CreatedBy   string    `json:"createdBy"`
		ChangedBy   string    `json:"changedBy"`
		BrokerType  string    `json:"brokerType"`
		Address     struct {
			CreatedTime   time.Time `json:"createdTime"`
			UpdatedTime   time.Time `json:"updatedTime"`
			CreatedBy     string    `json:"createdBy"`
			ChangedBy     string    `json:"changedBy"`
			AddressLevels []struct {
				Name             string `json:"name"`
				AddressLevelType string `json:"addressLevelType"`
				EnumVersionId    string `json:"enumVersionId"`
			} `json:"addressLevels"`
			AddressType string `json:"addressType"`
			Id          string `json:"id"`
			Type        string `json:"type"`
		} `json:"address"`
		KeySchemaVersionId     string `json:"keySchemaVersionId"`
		KeySchemaPrimitiveType string `json:"keySchemaPrimitiveType"`
		Id                     string `json:"id"`
		Type                   string `json:"type"`
	} `json:"deliveryDescriptor"`
	StateId          string `json:"stateId"`
	CustomAttributes []struct {
		CustomAttributeDefinitionId   string `json:"customAttributeDefinitionId"`
		CustomAttributeDefinitionName string `json:"customAttributeDefinitionName"`
		Value                         string `json:"value"`
	} `json:"customAttributes"`
	MessagingServiceIds []string `json:"messagingServiceIds"`
	Type                string   `json:"type"`
}

type EventVersionListResponse struct {
	EventVersions []EventVersion `json:"data"`
	Meta          Meta           `json:"meta"`
}

type EventVersionListPaginator struct {
	firstPage  bool
	pageNumber int
	count      int
	pageSize   int
	nextPage   int
	totalPages int
	config     *RequestConfig
	client     *Client
}

func buildEventVersionCSV(ids []string) string {
	// produced := r.ApplicationVersion.DeclaredProducedEventVersionIds
	// st := strings.Trim(strings.Join(strings.Fields(fmt.Sprint(produced)), ", "), "[]")
	// r.ApplicationVersion.DeclaredProducedEventVersionIdsAsString = st
	for index := range ids {
		ids[index] = fmt.Sprintf("'%s'", ids[index])
	}
	st := strings.Trim(strings.Join(strings.Fields(fmt.Sprint(ids)), ", "), "[]")
	return st
}

func (lp *EventVersionListPaginator) HasNoMorePages() bool {
	return lp.nextPage == 0 || lp.pageNumber > lp.totalPages
}

func (lp *EventVersionListPaginator) NextPage() ([]EventVersion, Meta, error) {
	if !lp.firstPage && lp.HasNoMorePages() {
		return nil, Meta{}, fmt.Errorf("no more pages available")
	}
	var r = &EventVersionListResponse{}
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

	for index := range r.EventVersions {
		r.EventVersions[index].DeclaredProducingApplicationVersionIdsAsString = buildEventVersionCSV(r.EventVersions[index].DeclaredProducingApplicationVersionIds)
		r.EventVersions[index].DeclaredConsumingApplicationVersionIdsAsString = buildEventVersionCSV(r.EventVersions[index].DeclaredConsumingApplicationVersionIds)
		r.EventVersions[index].ProducingEventApiVersionIdsAsString = buildEventVersionCSV(r.EventVersions[index].ProducingEventApiVersionIds)
		r.EventVersions[index].ConsumingEventApiVersionIdsAsString = buildEventVersionCSV(r.EventVersions[index].ConsumingEventApiVersionIds)
	}

	return r.EventVersions, r.Meta, nil
}

func (at *Client) NewEventVersionListPaginator(args map[string]string) *EventVersionListPaginator {
	var config = NewRequestConfig("architecture/eventVersions")
	config.Pagination = NewRequestPagination()
	ids, exists := args["ids"]
	if exists && len(ids) > 0 {
		config.Params.Add("ids", ids)
	}
	eventIds, exists := args["eventIds"]
	if exists && len(eventIds) > 0 {
		config.Params.Add("eventIds", eventIds)
	}

	var p = &EventVersionListPaginator{
		firstPage:  true,
		pageNumber: 1,
		pageSize:   defaultPageSize,
		config:     config,
		client:     at,
	}

	return p
}

type EventVersionGetResponse struct {
	EventVersion EventVersion `json:"data"`
	Meta         Meta         `json:"meta"`
}

func (at *Client) GetEventVersion(eventVersionId string) (*EventVersion, error) {
	var config = NewRequestConfig(fmt.Sprintf(`architecture/eventVersions/%s`, eventVersionId))

	var r = &EventVersionGetResponse{}
	var _, err = at.Get(config, &r)
	if err != nil {
		return nil, at.handleKnownErrors(err)
	}

	r.EventVersion.DeclaredProducingApplicationVersionIdsAsString = buildEventVersionCSV(r.EventVersion.DeclaredProducingApplicationVersionIds)
	r.EventVersion.DeclaredConsumingApplicationVersionIdsAsString = buildEventVersionCSV(r.EventVersion.DeclaredConsumingApplicationVersionIds)
	r.EventVersion.ProducingEventApiVersionIdsAsString = buildEventVersionCSV(r.EventVersion.ProducingEventApiVersionIds)
	r.EventVersion.ConsumingEventApiVersionIdsAsString = buildEventVersionCSV(r.EventVersion.ConsumingEventApiVersionIds)
	return &r.EventVersion, nil
}
