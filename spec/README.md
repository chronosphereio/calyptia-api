# Calyptia Cloud API - OpenAPI specification

OpenAPI spec for the API to [Calyptia Core](https://core.calyptia.com).

## Spec

This directory contains an OpenAPI specification.
You can use it to generate clients for your favorite programming language or preview it
using [SwaggerUI](https://editor.swagger.io/?url=https://raw.githubusercontent.com/chronosphereio/calyptia-api/main/spec/open-api.yml).

### Typescript codegen

Example using [openapi-typescript-codegen](https://www.npmjs.com/package/openapi-typescript-codegen).

```bash
npx openapi-typescript-codegen \
    --input ./spec/open-api.yml \
    --output ./ts-client \
    --name Client
```

### Go codegen

Example using [openapi-codegen](https://github.com/oapi-codegen/oapi-codegen).

```bash
oapi-codegen spec/open-api.yml
```
