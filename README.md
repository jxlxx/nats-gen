:warning: work in progress...

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

