package solacesdk

import (
	"fmt"
)

type EventBrokerServiceVersion struct {
	Version                 string   `json:"version"`
	SupportedServiceClasses []string `json:"supportedServiceClasses"`
	Capabilities            struct {
	} `json:"capabilities"`
}

type EventBrokerServiceVersionResponse struct {
	EventBrokerServiceVersions []EventBrokerServiceVersion `json:"data"`
}

func (at *Client) GetEventBrokerServiceVersions(id string) ([]EventBrokerServiceVersion, error) {
	var config = NewRequestConfig(fmt.Sprintf(`missionControl/datacenters/%s/eventBrokerServiceVersions`, id))
	var result = &EventBrokerServiceVersionResponse{}
	var _, err = at.Get(config, &result)
	if err != nil {
		return nil, at.handleKnownErrors(err)
	}

	// eventBrokerServiceVersions := []EventBrokerServiceVersion{}
	// for index := range result.BrokerServiceVersions {
	// 	eventBrokerServiceVersions = append(eventBrokerServiceVersions, EventBrokerServiceVersion{
	// 		Id:                      id,
	// 		Version:                 result.BrokerServiceVersions[index].Version,
	// 		SupportedServiceClasses: result.BrokerServiceVersions[index].SupportedServiceClasses,
	// 		Capabilities:            result.BrokerServiceVersions[index].Capabilities,
	// 	})
	// }
	// log.Println(fmt.Sprintf("DEBUGGING GET ORIG", fmt.Sprintf("%+v", result.BrokerServiceVersions)))
	// log.Println(fmt.Sprintf("DEBUGGING GET MOD", fmt.Sprintf("%+v", eventBrokerServiceVersions)))

	return result.EventBrokerServiceVersions, nil
}
