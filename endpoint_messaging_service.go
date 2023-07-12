package solace

import (
	"fmt"
	"time"
)

type MessagingService struct {
	CreatedTime                   time.Time `json:"createdTime"`
	UpdatedTime                   time.Time `json:"updatedTime"`
	CreatedBy                     string    `json:"createdBy"`
	ChangedBy                     string    `json:"changedBy"`
	Id                            string    `json:"id"`
	EventMeshId                   string    `json:"eventMeshId"`
	RuntimeAgentId                string    `json:"runtimeAgentId"`
	SolaceCloudMessagingServiceId string    `json:"solaceCloudMessagingServiceId"`
	MessagingServiceType          string    `json:"messagingServiceType"`
	Name                          string    `json:"name"`
	MessagingServiceConnections   []struct {
		CreatedTime        time.Time `json:"createdTime"`
		UpdatedTime        time.Time `json:"updatedTime"`
		CreatedBy          string    `json:"createdBy"`
		ChangedBy          string    `json:"changedBy"`
		Id                 string    `json:"id"`
		MessagingServiceId string    `json:"messagingServiceId"`
		Name               string    `json:"name"`
		URL                string    `json:"url"`
		Protocol           string    `json:"protocol"`
		ProtocolVersion    string    `json:"protocolVersion"`
		Bindings           struct {
			MsgVpn string `json:"msgVpn"`
		} `json:"bindings"`
		MessagingServiceAuthentications []struct {
			CreatedTime                  time.Time `json:"createdTime"`
			UpdatedTime                  time.Time `json:"updatedTime"`
			CreatedBy                    string    `json:"createdBy"`
			ChangedBy                    string    `json:"changedBy"`
			Id                           string    `json:"id"`
			MessagingServiceConnectionId string    `json:"messagingServiceConnectionId"`
			Name                         string    `json:"name"`
			AuthenticationType           string    `json:"authenticationType"`
			AuthenticationDetails        struct {
				BrokerOwner string `json:"broker owner"`
			} `json:"authenticationDetails"`
			MessagingServiceCredentials []struct {
				CreatedTime                      time.Time `json:"createdTime"`
				UpdatedTime                      time.Time `json:"updatedTime"`
				CreatedBy                        string    `json:"createdBy"`
				ChangedBy                        string    `json:"changedBy"`
				Id                               string    `json:"id"`
				MessagingServiceAuthenticationId string    `json:"messagingServiceAuthenticationId"`
				Name                             string    `json:"name"`
				Credentials                      struct {
					Username string `json:"username"`
					Password string `json:"password"`
				} `json:"credentials"`
				Type string `json:"type"`
			} `json:"messagingServiceCredentials"`
			Type string `json:"type"`
		} `json:"messagingServiceAuthentications"`
		Type string `json:"type"`
	} `json:"messagingServiceConnections"`
	EventManagementAgentId string `json:"eventManagementAgentId"`
	Type                   string `json:"type"`
}

type MessagingServiceListResponse struct {
	MessagingServices []MessagingService `json:"data"`
	Meta              Meta               `json:"meta"`
}

type MessagingServiceListPaginator struct {
	firstPage  bool
	pageNumber int
	count      int
	pageSize   int
	nextPage   int
	totalPages int
	config     *RequestConfig
	client     *Client
}

func (lp *MessagingServiceListPaginator) HasNoMorePages() bool {
	return lp.nextPage == 0 || lp.pageNumber > lp.totalPages
}

func (lp *MessagingServiceListPaginator) NextPage() ([]MessagingService, Meta, error) {
	if !lp.firstPage && lp.HasNoMorePages() {
		return nil, Meta{}, fmt.Errorf("no more pages available")
	}
	var r = &MessagingServiceListResponse{}
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

	return r.MessagingServices, r.Meta, nil
}

func (at *Client) NewMessagingServiceListPaginator() *MessagingServiceListPaginator {
	var config = NewRequestConfig("architecture/messagingServices")
	config.Pagination = NewRequestPagination()

	var p = &MessagingServiceListPaginator{
		firstPage:  true,
		pageNumber: 1,
		pageSize:   defaultPageSize,
		config:     config,
		client:     at,
	}

	return p
}

type MessagingServiceGetResponse struct {
	MessagingService MessagingService `json:"data"`
	Meta             Meta             `json:"meta"`
}

func (at *Client) GetMessagingService(messagingServiceId string) (*MessagingService, error) {
	var config = NewRequestConfig(fmt.Sprintf(`architecture/messagingServices/%s`, messagingServiceId))

	var r = &MessagingServiceGetResponse{}
	var _, err = at.Get(config, &r)
	if err != nil {
		return nil, at.handleKnownErrors(err)
	}

	return &r.MessagingService, nil
}
