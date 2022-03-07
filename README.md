# Calyptia Cloud API - OpenAPI specification, Golang reference client and types

[![Go Reference](https://pkg.go.dev/badge/github.com/calyptia/api.svg)](https://pkg.go.dev/github.com/calyptia/api)
[![CI](https://github.com/calyptia/go-repo-template/actions/workflows/ci.yml/badge.svg?branch=main)](https://github.com/calyptia/go-repo-template/actions/workflows/ci.yml)
[![codecov](https://codecov.io/gh/calyptia/api/branch/main/graph/badge.svg?token=FUCFZ7JRAS)](https://codecov.io/gh/calyptia/api)

OpenAPI spec, Golang types and client for the API to [Calyptia Cloud](https://cloud.calyptia.com).

## Install

```bash
go get github.com/calyptia/api
```

## Client usage

Get an API key from [Calyptia Cloud](https://cloud.calyptia.com) under settings.

```go
package main

import "github.com/calyptia/api/client"

func main() {
    c := client.New()
    c.SetProjectToken("YOUR_API_KEY_HERE")
}
```

API keys are bound to an specific project.
With that API key you cannot list all your other projects,
or create more API keys.
But you can perform all other actions within a project:
manage agents, aggregators and pipelines, invite members, etc.

### User authentication

Machines using this client should prefer API key authorization,
but if you need user authentication, you must setup [Auth0](https://auth0.com).
After you successfully login and get an access token,
you can create an authenticated client like so:

```go
tok := &oauth2.Token{
    AccessToken:  "YOUR_AUTH0_ACCESS_TOKEN",
    TokenType:    "Bearer",
    Expiry:       auth0Expiry,
}
tokSrc := oauth2.StaticTokenSource(tok)
c.Client = oauth2.NewClient(ctx, tokSrc)
```

Refer to [client/client_test.go](https://github.com/calyptia/api/blob/eec74522b60638539bdb7f2334548d3c4cda813d/client/client_test.go#L528-L531)
for a testing example.

## Spec

The `/spec` directory contains an OpenAPI specification.
You can use it to generate clients for your favorite programming language.
Or preview it using [SwaggerUI](https://editor.swagger.io/?url=https://raw.githubusercontent.com/calyptia/api/main/spec/open-api.yml).

### Typescript codegen

Example using [openapi-typescript-codegen](https://www.npmjs.com/package/openapi-typescript-codegen).

```bash
npx openapi-typescript-codegen --input ./spec/open-api.yml --output ./ts-client --name Client
```
