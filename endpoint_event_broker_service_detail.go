package solace

import (
	"fmt"
	"time"
)

type EventBrokerServiceDetail struct {
	CreatedTime                time.Time `json:"createdTime"`
	UpdatedTime                time.Time `json:"updatedTime"`
	CreatedBy                  string    `json:"createdBy"`
	ChangedBy                  string    `json:"changedBy"`
	Id                         string    `json:"id"`
	Type                       string    `json:"type"`
	Name                       string    `json:"name"`
	OwnedBy                    string    `json:"ownedBy"`
	InfrastructureId           string    `json:"infrastructureId"`
	DatacenterId               string    `json:"datacenterId"`
	ServiceClassId             string    `json:"serviceClassId"`
	EventMeshId                string    `json:"eventMeshId"`
	OngoingOperationIds        []string  `json:"ongoingOperationIds"`
	AdminState                 string    `json:"adminState"`
	CreationState              string    `json:"creationState"`
	Locked                     bool      `json:"locked"`
	DefaultManagementHostname  string    `json:"defaultManagementHostname"`
	ServiceConnectionEndpoints []struct {
		Id             string   `json:"id"`
		Type           string   `json:"type"`
		Name           string   `json:"name"`
		Description    string   `json:"description"`
		AccessType     string   `json:"accessType"`
		K8SServiceType string   `json:"k8sServiceType"`
		K8SServiceId   string   `json:"k8sServiceId"`
		HostNames      []string `json:"hostNames"`
		Ports          struct {
			ServiceWebPlainTextListenPort          int `json:"serviceWebPlainTextListenPort"`
			ServiceRestIncomingTLSListenPort       int `json:"serviceRestIncomingTlsListenPort"`
			ServiceManagementTLSListenPort         int `json:"serviceManagementTlsListenPort"`
			ServiceAmqpPlainTextListenPort         int `json:"serviceAmqpPlainTextListenPort"`
			ServiceMqttWebSocketListenPort         int `json:"serviceMqttWebSocketListenPort"`
			ServiceRestIncomingPlainTextListenPort int `json:"serviceRestIncomingPlainTextListenPort"`
			ServiceWebTLSListenPort                int `json:"serviceWebTlsListenPort"`
			ServiceSmfCompressedListenPort         int `json:"serviceSmfCompressedListenPort"`
			ServiceMqttPlainTextListenPort         int `json:"serviceMqttPlainTextListenPort"`
			ServiceSmfPlainTextListenPort          int `json:"serviceSmfPlainTextListenPort"`
			ServiceAmqpTLSListenPort               int `json:"serviceAmqpTlsListenPort"`
			ServiceMqttTLSListenPort               int `json:"serviceMqttTlsListenPort"`
			ServiceSmfTLSListenPort                int `json:"serviceSmfTlsListenPort"`
			ServiceMqttTLSWebSocketListenPort      int `json:"serviceMqttTlsWebSocketListenPort"`
			ManagementSSHTLSListenPort             int `json:"managementSshTlsListenPort"`
		} `json:"ports"`
	} `json:"serviceConnectionEndpoints"`
	Broker struct {
		Version                      string `json:"version"`
		VersionFamily                string `json:"versionFamily"`
		ServicePackageId             string `json:"servicePackageId"`
		MaxSpoolUsage                int    `json:"maxSpoolUsage"`
		DiskSize                     int    `json:"diskSize"`
		RedundancyGroupSslEnabled    bool   `json:"redundancyGroupSslEnabled"`
		ConfigSyncSslEnabled         bool   `json:"configSyncSslEnabled"`
		MonitoringMode               string `json:"monitoringMode"`
		ClientCertificateAuthorities []struct {
			Name string `json:"name"`
		} `json:"clientCertificateAuthorities"`
		DomainCertificateAuthorities []struct {
			Name string `json:"name"`
		} `json:"domainCertificateAuthorities"`
		TLSStandardDomainCertificateAuthoritiesEnabled bool `json:"tlsStandardDomainCertificateAuthoritiesEnabled"`
		LdapProfiles                                   []struct {
			Name string `json:"name"`
		} `json:"ldapProfiles"`
		Cluster struct {
			Name                        string   `json:"name"`
			Password                    string   `json:"password"`
			RemoteAddress               string   `json:"remoteAddress"`
			PrimaryRouterName           string   `json:"primaryRouterName"`
			BackupRouterName            string   `json:"backupRouterName"`
			MonitoringRouterName        string   `json:"monitoringRouterName"`
			SupportedAuthenticationMode []string `json:"supportedAuthenticationMode"`
		} `json:"cluster"`
		ManagementReadOnlyLoginCredential struct {
			Username string `json:"username"`
			Password string `json:"password"`
			Token    string `json:"token"`
		} `json:"managementReadOnlyLoginCredential"`
		MsgVpns []struct {
			MsgVpnName                                  string `json:"msgVpnName"`
			AuthenticationBasicEnabled                  bool   `json:"authenticationBasicEnabled"`
			AuthenticationBasicType                     string `json:"authenticationBasicType"`
			AuthenticationClientCertEnabled             bool   `json:"authenticationClientCertEnabled"`
			AuthenticationClientCertValidateDateEnabled bool   `json:"authenticationClientCertValidateDateEnabled"`
			ClientProfiles                              []struct {
				Name string `json:"name"`
			} `json:"clientProfiles"`
			Enabled                        bool `json:"enabled"`
			EventLargeMsgThreshold         int  `json:"eventLargeMsgThreshold"`
			ManagementAdminLoginCredential struct {
				Username string `json:"username"`
				Password string `json:"password"`
				Token    string `json:"token"`
			} `json:"managementAdminLoginCredential"`
			ServiceLoginCredential struct {
				Username string `json:"username"`
				Password string `json:"password"`
			} `json:"serviceLoginCredential"`
			MaxConnectionCount        int `json:"maxConnectionCount"`
			MaxEgressFlowCount        int `json:"maxEgressFlowCount"`
			MaxEndpointCount          int `json:"maxEndpointCount"`
			MaxIngressFlowCount       int `json:"maxIngressFlowCount"`
			MaxMsgSpoolUsage          int `json:"maxMsgSpoolUsage"`
			MaxSubscriptionCount      int `json:"maxSubscriptionCount"`
			MaxTransactedSessionCount int `json:"maxTransactedSessionCount"`
			MaxTransactionCount       int `json:"maxTransactionCount"`
			SempOverMessageBus        struct {
				SempOverMsgBusEnabled              bool `json:"sempOverMsgBusEnabled"`
				SempAccessToShowCmdsEnabled        bool `json:"sempAccessToShowCmdsEnabled"`
				SempAccessToAdminCmdsEnabled       bool `json:"sempAccessToAdminCmdsEnabled"`
				SempAccessToClientAdminCmdsEnabled bool `json:"sempAccessToClientAdminCmdsEnabled"`
				SempAccessToCacheCmdsEnabled       bool `json:"sempAccessToCacheCmdsEnabled"`
			} `json:"sempOverMessageBus"`
			SubDomainName string `json:"subDomainName"`
			TruststoreURI string `json:"truststoreUri"`
		} `json:"msgVpns"`
	} `json:"broker"`
}

type EventBrokerServiceDetailGetResponse struct {
	EventBrokerServiceDetail EventBrokerServiceDetail `json:"data"`
	Pagination               Pagination               `json:"meta"`
}

func (at *Client) GetEventBrokerServiceDetail(eventBrokerServiceId string) (*EventBrokerServiceDetail, error) {
	var config = NewRequestConfig(fmt.Sprintf(`missionControl/eventBrokerServices/%s`, eventBrokerServiceId))
	config.Params.Add("expand", "broker")
	config.Params.Add("expand", "serviceConnectionEndpoints")

	var result = &EventBrokerServiceDetailGetResponse{}
	var _, err = at.Get(config, &result)
	if err != nil {
		return nil, at.handleKnownErrors(err)
	}

	return &result.EventBrokerServiceDetail, nil
}
