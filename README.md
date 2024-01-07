# NATS microservice generation

:warning: work in progress...


## Install 

```sh
go install github.com/jxlxx/nats-gen/cmd/nats-gen@latest
```

## Getting started

1. Define a microservice in yaml:

```yaml
microservices:
  -  
    package: cats
    config:
      name: CatService
      description: This is a cat service
      version: 0.0.1
    targetFile: cats.gen.go
    groups:
      - 
          name: cats
          description: This is the cat microservice
          subject: 
            name: catSubject
            tokens: [cats]
    endpoints:
      - 
          name: new
          operationId: NewCat
          group: cats
          subject: 
            tokens: [new]
          payload:
            name: cat
            schema: CatIntake
      - 
          name: get
          operationId: GetCat
          group: cats
          subject: 
            tokens: [catID]
            parameters:
              - name: catID
                type: string
                required: true
    schemas:
      - name: CatIntake
        fields:
          - 
            name: Name
            type: string
          - 
            name: BirthYear
            type: int
```

2. Generate the service

```sh
nats-gen -c cats.yaml  
```

## :test_tube: Generating test files

Add a testing configuration to the microservice definition.

```yaml
  -  
    package: cats
    config:
      name: CatService
      description: This is a cat service
      version: 0.0.1
    targetFile: cats.gen.go
    testing:
      name: CatService
      file: gen/cats/cats.gen_test.go
      package: cats
      enable: true
      options: 
        Name: CatsService
        Version: 0.0.1
    groups:
      - 
          name: cats
          description: This is the cat microservice
          subject: 
            name: catSubject
            tokens: [cats]
    # ....
```

## :pencil2: TODO:

:foot: features

- ~generate types~
- ~process subjects~
- ~ser/de~
- ~generate tests~
- ~use test containers~
- generate documentation
  - md
  - html
- kv crud
- subscribers
- publishers
- optional writer interface
- optionally generate HTTP wrapper
- vendor option
- generate open api spec
- TDD unit tests generated from yaml
- cue

:art: polishing

- add schema info to metadata
- metadata to yaml
- handle different data types nicely
  - enum types
  - array types  
  - https://github.com/deepmap/oapi-codegen/blob/46d269400f4bd1f4da2f5e6b84cdf7f3f2d753dd/pkg/codegen/schema.go#L519
- add the descriptions to comments in the code
- informative preamble
- informative readme
- i think paths are not working for files, need to check
- i maybe can just get rid of Argument type and just use Parameter who cares
- making tests should be optional
- should have to specifify a test package
- adding config/enabling js to test containers
- should be able to split ip services and types etc into seperate files
- add yaml/json tags to the struct fields

