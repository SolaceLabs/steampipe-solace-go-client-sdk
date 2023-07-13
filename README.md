[![Contributor Covenant](https://img.shields.io/badge/Contributor%20Covenant-v2.0%20adopted-ff69b4.svg)](CODE_OF_CONDUCT.md)

# Go Client SDK (lite) for Solace Cloud REST APIs

## Overview
This project creates a lite Go Client SDK to access PubSub+ Cloud REST APIs. The reason this is a lite-SDK because it implements only GET calls to Solace Cloud resources. This SDK is utilized by the steampipe plugin for solace project to access Cloud REST APIs via SQL. 

## Getting started quickly
1. Get the SDK
```shell
go get github.com/SolaceLabs/steampipe-solace-go-client-sdk
```
2. Create a client
```go
var config, err = solace.NewConfig(apiToken, apiUrl, rateLimit)
if err != nil {
  return nil, err
}

solaceClient = solace.GetClient(config)
```

3. Fetch resources
```go
var config = NewRequestConfig(fmt.Sprintf(`architecture/applicationDomains/%s`, domainId))

var r = &ApplicationDomainGetResponse{}
var _, err = solaceClient.Get(config, &r)
if err != nil {
  return nil, at.handleKnownErrors(err)
}

return &r.ApplicationDomain, nil
```
## Resources
For more use-cases check:

[Steampipe Plugin for Solace](https://github.com/SolaceLabs/steampipe-plugin-solace)

Further reading:

[Solace Cloud REST API documentation](https://api.solace.dev/cloud/reference/using-the-v2-rest-apis-for-pubsub-cloud)

Get involved:

[Issues](https://github.com/SolaceLabs/steampipe-solace-go-client-sdk/issues)

This is not an officially supported Solace product.

For more information try these resources:
- Ask the [Solace Community](https://solace.community)
- The Solace Developer Portal website at: https://solace.dev


## Contributing
Contributions are encouraged! Please read [CONTRIBUTING.md](CONTRIBUTING.md) for details on our code of conduct, and the process for submitting pull requests to us.

## Authors
See the list of [contributors](https://github.com/solacecommunity/<github-repo>/graphs/contributors) who participated in this project.

## License
See the [LICENSE](LICENSE) file for details.
