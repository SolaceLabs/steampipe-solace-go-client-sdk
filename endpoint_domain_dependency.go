package solace

import (
	"fmt"
	"log"
	"strings"
)

type Cache = map[string]any
type ArrayCache = map[string][]string
type void struct{}

type DomainDependency struct {
	//TYPE
	Resource   string `json:"resource"`
	ResourceId string `json:"resourceId"`
	//DOMAINS
	Id    string `json:"id"`
	Name  string `json:"name"`
	Stats struct {
		SchemaCount          int `json:"schemaCount"`
		EventCount           int `json:"eventCount"`
		ApplicationCount     int `json:"applicationCount"`
		EnumCount            int `json:"enumCount"`
		EventApiCount        int `json:"eventApiCount"`
		EventApiProductCount int `json:"eventApiProductCount"`
	} `json:"stats"`
	//APPLICATIONS
	Application struct {
		ApplicationId    string `json:"applicationId"`
		ApplicationName  string `json:"applicationName"`
		NumberOfVersions int    `json:"numberOfVersions"`
	} `json:"application"`
}

var cache = make(Cache)

// LIST OF DOMAINS
var cacheDomains = make(ArrayCache)

// MAP OF DOMAIN -> RESOURCES
var cacheDomainApplications = make(ArrayCache)
var cacheDomainEvents = make(ArrayCache)
var cacheDomainSchemas = make(ArrayCache)
var cacheDomainEnums = make(ArrayCache)
var cacheDomainEventApis = make(ArrayCache)
var cacheDomainEventApiProducts = make(ArrayCache)

// MAP OF APPLICATION -> APPLICATION VERSIONS
var cacheDomainApplicationVersions = make(ArrayCache)

var member void
var hierarchyRecords = []DomainDependency{}

func (at *Client) getDomainEventApiProductDetails() {
	var ids = []string{}
	for key := range cacheDomains {
		ids = append(ids, key)
	}

	args := map[string]string{}
	args["applicationDomainIds"] = strings.Trim(strings.Join(strings.Fields(fmt.Sprint(ids)), ", "), "[]")

	var up = at.NewEventApiProductListPaginator(args)
	pagesLeft := true
	count := 0
	var eventApiProductIds = []string{}
	for pagesLeft {
		eventApiProducts, meta, err := up.NextPage()
		if err != nil {
			if fmt.Sprintf("%s", err.Error()) != "no more pages available" {
				log.Println(err.Error())
			}

			pagesLeft = false
		} else {
			for _, eventApiProduct := range eventApiProducts {
				cache[eventApiProduct.Id] = eventApiProduct
				cacheDomainEventApiProducts[eventApiProduct.ApplicationDomainId] = append(cacheDomainEventApiProducts[eventApiProduct.ApplicationDomainId], eventApiProduct.Id)
				eventApiProductIds = append(eventApiProductIds, eventApiProduct.Id)
			}

			count += meta.Pagination.Count
		}
	}

	log.Println(count, " - EVENT API PRODUCTS FETCHED")

	args = map[string]string{}
	args["eventApiProductIds"] = strings.Trim(strings.Join(strings.Fields(fmt.Sprint(eventApiProductIds)), ", "), "[]")

	var upS = at.NewEventApiProductVersionListPaginator(args)
	pagesLeftS := true
	countS := 0
	for pagesLeftS {
		eventApiProductVersions, meta, err := upS.NextPage()
		if err != nil {
			if fmt.Sprintf("%s", err.Error()) != "no more pages available" {
				log.Println(err.Error())
			}

			pagesLeftS = false
		} else {
			for _, eventApiProductVersion := range eventApiProductVersions {
				cache[eventApiProductVersion.Id] = eventApiProductVersions
			}
			countS += meta.Pagination.Count
		}
	}
	log.Println(countS, " - EVENT API PRODUCT VERSIONS FETCHED")
}

func (at *Client) getDomainEventApiDetails() {
	var ids = []string{}
	for key := range cacheDomains {
		ids = append(ids, key)
	}

	args := map[string]string{}
	args["applicationDomainIds"] = strings.Trim(strings.Join(strings.Fields(fmt.Sprint(ids)), ", "), "[]")

	var up = at.NewEventApiListPaginator(args)
	pagesLeft := true
	count := 0
	var eventApiIds = []string{}
	for pagesLeft {
		eventApis, meta, err := up.NextPage()
		if err != nil {
			if fmt.Sprintf("%s", err.Error()) != "no more pages available" {
				log.Println(err.Error())
			}

			pagesLeft = false
		} else {
			for _, eventApi := range eventApis {
				cache[eventApi.Id] = eventApi
				cacheDomainEventApis[eventApi.ApplicationDomainId] = append(cacheDomainEventApis[eventApi.ApplicationDomainId], eventApi.Id)
				eventApiIds = append(eventApiIds, eventApi.Id)
			}

			count += meta.Pagination.Count
		}
	}

	log.Println(count, " - EVENT APIS FETCHED")

	args = map[string]string{}
	args["eventApiIds"] = strings.Trim(strings.Join(strings.Fields(fmt.Sprint(eventApiIds)), ", "), "[]")

	var upS = at.NewEventApiVersionListPaginator(args)
	pagesLeftS := true
	countS := 0
	for pagesLeftS {
		eventApiVersions, meta, err := upS.NextPage()
		if err != nil {
			if fmt.Sprintf("%s", err.Error()) != "no more pages available" {
				log.Println(err.Error())
			}

			pagesLeftS = false
		} else {
			for _, eventApiVersion := range eventApiVersions {
				cache[eventApiVersion.Id] = eventApiVersions
			}
			countS += meta.Pagination.Count
		}
	}
	log.Println(countS, " - EVENT API VERSIONS FETCHED")
}

func (at *Client) getDomainEnumDetails() {
	var ids = []string{}
	for key := range cacheDomains {
		ids = append(ids, key)
	}

	args := map[string]string{}
	args["applicationDomainIds"] = strings.Trim(strings.Join(strings.Fields(fmt.Sprint(ids)), ", "), "[]")

	var up = at.NewEnumListPaginator(args)
	pagesLeft := true
	count := 0
	var enumIds = []string{}
	for pagesLeft {
		enums, meta, err := up.NextPage()
		if err != nil {
			if fmt.Sprintf("%s", err.Error()) != "no more pages available" {
				log.Println(err.Error())
			}

			pagesLeft = false
		} else {
			for _, enum := range enums {
				cache[enum.Id] = enum
				cacheDomainEnums[enum.ApplicationDomainId] = append(cacheDomainEnums[enum.ApplicationDomainId], enum.Id)
				enumIds = append(enumIds, enum.Id)
			}

			count += meta.Pagination.Count
		}
	}

	log.Println(count, " - ENUMS FETCHED")

	args = map[string]string{}
	args["enumIds"] = strings.Trim(strings.Join(strings.Fields(fmt.Sprint(enumIds)), ", "), "[]")

	var upS = at.NewEnumVersionListPaginator(args)
	pagesLeftS := true
	countS := 0
	for pagesLeftS {
		enumVersions, meta, err := upS.NextPage()
		if err != nil {
			if fmt.Sprintf("%s", err.Error()) != "no more pages available" {
				log.Println(err.Error())
			}

			pagesLeftS = false
		} else {
			for _, enumVersion := range enumVersions {
				cache[enumVersion.Id] = enumVersions
			}
			countS += meta.Pagination.Count
		}
	}
	log.Println(countS, " - ENUM VERSIONS FETCHED")
}

func (at *Client) getDomainSchemaDetails() {
	var ids = []string{}
	for key := range cacheDomains {
		ids = append(ids, key)
	}

	args := map[string]string{}
	args["applicationDomainIds"] = strings.Trim(strings.Join(strings.Fields(fmt.Sprint(ids)), ", "), "[]")

	var up = at.NewSchemaListPaginator(args)
	pagesLeft := true
	count := 0
	var schemaIds = []string{}
	for pagesLeft {
		schemas, meta, err := up.NextPage()
		if err != nil {
			if fmt.Sprintf("%s", err.Error()) != "no more pages available" {
				log.Println(err.Error())
			}

			pagesLeft = false
		} else {
			for _, schema := range schemas {
				cache[schema.Id] = schema
				cacheDomainSchemas[schema.ApplicationDomainId] = append(cacheDomainSchemas[schema.ApplicationDomainId], schema.Id)
				schemaIds = append(schemaIds, schema.Id)
			}

			count += meta.Pagination.Count
		}
	}

	log.Println(count, " - SCHEMAS FETCHED")

	args = map[string]string{}
	args["schemaIds"] = strings.Trim(strings.Join(strings.Fields(fmt.Sprint(schemaIds)), ", "), "[]")

	var upS = at.NewSchemaVersionListPaginator(args)
	pagesLeftS := true
	countS := 0
	for pagesLeftS {
		schemaVersions, meta, err := upS.NextPage()
		if err != nil {
			if fmt.Sprintf("%s", err.Error()) != "no more pages available" {
				log.Println(err.Error())
			}

			pagesLeftS = false
		} else {
			for _, schemaVersion := range schemaVersions {
				cache[schemaVersion.Id] = schemaVersions
			}
			countS += meta.Pagination.Count
		}
	}
	log.Println(countS, " - SCHEMA VERSIONS FETCHED")
}

func (at *Client) getDomainEventDetails() {
	var ids = []string{}
	for key := range cacheDomains {
		ids = append(ids, key)
	}

	args := map[string]string{}
	args["applicationDomainIds"] = strings.Trim(strings.Join(strings.Fields(fmt.Sprint(ids)), ", "), "[]")

	var up = at.NewEventListPaginator(args)
	pagesLeft := true
	count := 0
	var eventIds = []string{}
	for pagesLeft {
		events, meta, err := up.NextPage()
		if err != nil {
			if fmt.Sprintf("%s", err.Error()) != "no more pages available" {
				log.Println(err.Error())
			}

			pagesLeft = false
		} else {
			for _, event := range events {
				cache[event.Id] = event
				cacheDomainEvents[event.ApplicationDomainId] = append(cacheDomainEvents[event.ApplicationDomainId], event.Id)
				eventIds = append(eventIds, event.Id)
			}

			count += meta.Pagination.Count
		}
	}

	log.Println(count, " - EVENTS FETCHED")

	args = map[string]string{}
	args["eventIds"] = strings.Trim(strings.Join(strings.Fields(fmt.Sprint(eventIds)), ", "), "[]")

	var upE = at.NewEventVersionListPaginator(args)
	pagesLeftE := true
	countE := 0
	for pagesLeftE {
		eventVersions, meta, err := upE.NextPage()
		if err != nil {
			if fmt.Sprintf("%s", err.Error()) != "no more pages available" {
				log.Println(err.Error())
			}

			pagesLeftE = false
		} else {
			for _, eventVersion := range eventVersions {
				cache[eventVersion.Id] = eventVersion
			}
			countE += meta.Pagination.Count
		}
	}
	log.Println(countE, " - EVENT VERSIONS FETCHED")
}

func (at *Client) getDomainApplicationDetails() {
	for _, domains := range cacheDomains {
		for _, id := range domains {
			args := map[string]string{}
			args["applicationDomainId"] = id

			var up = at.NewApplicationListPaginator(args)
			pagesLeft := true
			count := 0
			for pagesLeft {
				applications, meta, err := up.NextPage()
				if err != nil {
					if fmt.Sprintf("%s", err.Error()) != "no more pages available" {
						log.Println(err.Error())
					}

					pagesLeft = false
				} else {
					for _, application := range applications {
						cache[application.Id] = application
						cacheDomainApplications[id] = append(cacheDomainApplications[id], application.Id)
					}

					count += meta.Pagination.Count
					log.Println(count, " - APPLICATIONS FETCHED")

					args := map[string]string{}
					args["applicationIds"] = strings.Trim(strings.Join(strings.Fields(fmt.Sprint(cacheDomainApplications[id])), ", "), "[]")

					var upA = at.NewApplicationVersionListPaginator(args)
					pagesLeftA := true
					countA := 0
					for pagesLeftA {
						applicationVersions, meta, err := upA.NextPage()
						if err != nil {
							if fmt.Sprintf("%s", err.Error()) != "no more pages available" {
								log.Println(err.Error())
							}

							pagesLeftA = false
						} else {
							for _, applicationVersion := range applicationVersions {
								cache[applicationVersion.Id] = applicationVersion
							}
							countA += meta.Pagination.Count
							log.Println(countA, " - APPLICATION VERSIONS FETCHED")
						}
					}
				}
			}
		}
	}
}

func (at *Client) constructApplication(domainRecord DomainDependency) {
	for _, applicationId := range cacheDomainApplications[domainRecord.Id] {
		application := cache[applicationId].(Application)
		applicationRecord := domainRecord
		applicationRecord.Resource = "APPLICATION"
		applicationRecord.ResourceId = application.Id
		applicationRecord.Application.ApplicationId = application.Id
		applicationRecord.Application.ApplicationName = application.Name
		applicationRecord.Application.NumberOfVersions = application.NumberOfVersions
		hierarchyRecords = append(hierarchyRecords, applicationRecord)
	}
}

// func (at *Client) constructEvent(domainRecord DomainDependency) {
// 	for _, eventId := range cacheDomainEvents[domainRecord.Id] {
// 		event := cache[eventId].(Event)
// 		eventRecord := domainRecord
// 		eventRecord.Resource = "EVENT"
// 		eventRecord.ResourceId = event.Id
// 		eventRecord.Event.EventId = event.Id
// 		eventRecord.Event.EventName = event.Name
// 		eventRecord.Event.NumberOfVersions = event.NumberOfVersions
// 		hierarchyRecords = append(hierarchyRecords, eventRecord)
// 	}
// }

// func (at *Client) constructSchema(domainRecord DomainDependency) {
// 	for _, schemaId := range cacheDomainSchemas[domainRecord.Id] {
// 		schema := cache[schemaId].(Schema)
// 		schemaRecord := domainRecord
// 		schemaRecord.Resource = "SCHEMA"
// 		schemaRecord.ResourceId = schema.Id
// 		schemaRecord.Schema.SchemaId = schema.Id
// 		schemaRecord.Schema.SchemaName = schema.Name
// 		schemaRecord.Schema.NumberOfVersions = schema.NumberOfVersions
// 		hierarchyRecords = append(hierarchyRecords, schemaRecord)
// 	}
// }

// func (at *Client) constructEnum(domainRecord DomainDependency) {
// 	for _, enumId := range cacheDomainEnums[domainRecord.Id] {
// 		enum := cache[enumId].(Enum)
// 		enumRecord := domainRecord
// 		enumRecord.Id = enum.Id
// 		enumRecord.Resource = "ENUM"
// 		enumRecord.Id = enum.Id
// 		enumRecord.Enum.EnumId = enum.Id
// 		enumRecord.Enum.EnumName = enum.Name
// 		enumRecord.Enum.NumberOfVersions = enum.NumberOfVersions
// 		hierarchyRecords = append(hierarchyRecords, enumRecord)
// 	}
// }

// func (at *Client) constructEventApi(domainRecord DomainDependency) {
// 	for _, eventApiId := range cacheDomainEventApis[domainRecord.Id] {
// 		eventApi := cache[eventApiId].(EventApi)
// 		eventApiRecord := domainRecord
// 		eventApiRecord.Resource = "EVENTAPI"
// 		eventApiRecord.Id = eventApi.Id
// 		eventApiRecord.EventApi.EventApiId = eventApi.Id
// 		eventApiRecord.EventApi.EventApiName = eventApi.Name
// 		eventApiRecord.EventApi.NumberOfVersions = eventApi.NumberOfVersions
// 		hierarchyRecords = append(hierarchyRecords, eventApiRecord)
// 	}
// }

// func (at *Client) constructEventApiProduct(domainRecord DomainDependency) {
// 	for _, eventApiProductId := range cacheDomainEventApiProducts[domainRecord.Id] {
// 		eventApiProduct := cache[eventApiProductId].(EventApiProduct)
// 		eventApiProductRecord := domainRecord
// 		eventApiProductRecord.Resource = "EVENTAPIPRODUCT"
// 		eventApiProductRecord.Id = eventApiProduct.Id
// 		eventApiProductRecord.EventApiProduct.EventApiProductId = eventApiProduct.Id
// 		eventApiProductRecord.EventApiProduct.EventApiProductName = eventApiProduct.Name
// 		eventApiProductRecord.EventApiProduct.NumberOfVersions = eventApiProduct.NumberOfVersions
// 		hierarchyRecords = append(hierarchyRecords, eventApiProductRecord)
// 	}
// }

func (at *Client) constructDomainDependency() {
	for _, domains := range cacheDomains {
		domain := ApplicationDomain{}
		for _, id := range domains {
			domain = cache[id].(ApplicationDomain)
			domainRecord := DomainDependency{}
			domainRecord.Id = domain.Id
			domainRecord.Resource = "DOMAIN"
			domainRecord.ResourceId = domain.Id
			domainRecord.Name = domain.Name
			domainRecord.Stats = domain.Stats
			hierarchyRecords = append(hierarchyRecords, domainRecord)
			if domain.Stats.ApplicationCount > 0 {
				at.constructApplication(domainRecord)
			}
			// if domain.Stats.EventCount > 0 {
			// 	at.constructApplication(domainRecord)
			// }
			// if domain.Stats.ApplicationCount > 0 {
			// 	at.constructApplication(domainRecord)
			// }
			// if domain.Stats.ApplicationCount > 0 {
			// 	at.constructApplication(domainRecord)
			// }
			// if domain.Stats.ApplicationCount > 0 {
			// 	at.constructApplication(domainRecord)
			// }
			// if domain.Stats.ApplicationCount > 0 {
			// 	at.constructApplication(domainRecord)
			// }
		}
	}
}

type DomainDependencyListResponse struct {
	DomainDependency []DomainDependency `json:"data"`
	Meta             Meta               `json:"meta"`
}

type DomainDependencyListPaginator struct {
	firstPage  bool
	pageNumber int
	count      int
	pageSize   int
	nextPage   int
	totalPages int
	config     *RequestConfig
	client     *Client
}

func (lp *DomainDependencyListPaginator) HasNoMorePages() bool {
	return lp.nextPage == 0 || lp.pageNumber > lp.totalPages
}

func (lp *DomainDependencyListPaginator) NextPage() ([]DomainDependency, Meta, error) {
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
	lp.config.Pagination.NextPage(0)

	var d = &DomainDependencyGetResponse{}
	for _, domain := range r.ApplicationDomains {
		dependencies, err := lp.client.GetDomainDependency(domain.Id)
		if err == nil {
			d.DomainDependencies = append(d.DomainDependencies, dependencies...)
		}
	}
	log.Println("DOMAIN HIERARCHY", fmt.Sprintf("%+v", d.DomainDependencies))

	return d.DomainDependencies, r.Meta, nil
}

func (at *Client) NewDomainDependencyListPaginator(args map[string]string) *DomainDependencyListPaginator {
	var config = NewRequestConfig("architecture/applicationDomains")
	config.Params.Add("include", "stats")
	value, exists := args["ids"]
	if exists && len(value) > 0 {
		config.Params.Add("ids", value)
	}
	config.Pagination = NewRequestSingleElementPagination(false)

	var p = &DomainDependencyListPaginator{
		firstPage:  true,
		pageNumber: 1,
		pageSize:   1,
		config:     config,
		client:     at,
	}

	return p
}

type DomainDependencyGetResponse struct {
	DomainDependencies []DomainDependency `json:"data"`
	Meta               Meta               `json:"meta"`
}

func (at *Client) GetDomainDependency(id string) ([]DomainDependency, error) {
	args := map[string]string{}
	args["ids"] = strings.Trim(strings.Join(strings.Fields(fmt.Sprint(id)), ", "), "[]")

	var up = at.NewApplicationDomainListPaginator(args)
	pagesLeft := true
	count := 0
	for pagesLeft {
		domains, meta, err := up.NextPage()
		if err != nil {
			if fmt.Sprintf("%s", err.Error()) != "no more pages available" {
				log.Println(err.Error())
			}

			pagesLeft = false
		} else {
			for _, domain := range domains {
				cache[domain.Id] = domain
				cacheDomains[domain.Id] = append(cacheDomains[domain.Id], domain.Id)
			}
			count += meta.Pagination.Count
			log.Println(count, " - DOMAINS FETCHED")

			at.getDomainApplicationDetails()
			// at.getDomainEventDetails()
			// at.getDomainSchemaDetails()
			// at.getDomainEnumDetails()
			// at.getDomainEventApiDetails()
			// at.getDomainEventApiProductDetails()
			at.constructDomainDependency()
			log.Println("DOMAIN HIERARCHY", fmt.Sprintf("%+v", hierarchyRecords))
		}
	}

	var r = &DomainDependencyGetResponse{}
	r.DomainDependencies = hierarchyRecords

	return r.DomainDependencies, nil
}
