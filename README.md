# Sample-service
This is a golang CRUD service for DVD store that communicates by GRPC and follows the principles of [Clean Architecture](http://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html "Clean Architecture") by Robert Martin.

It has a simplified business logic in order to concentrate on architecture, code organisation and practicing GRPC.

## Project structure
Project structure (mostly) follows [Standard Go Project Layout](https://github.com/golang-standards/project-layout "Standard Go Project Layout").

- `/cmd` - entry point of the app
- `/internal/dvdstore`  - application code (interfaces, transports, implementations)
- `/internal/dvdstore/grpc` -  GRPC transport
- `/internal/dvdstore/repository` - working with repositories, currently only postgresql
- `/internal/dvdstore/usecase` - business logic
- `/internal/models` - entities, exported errors, custom validations
- `/proto` - protobuf definition and proto-generated code

To follow dependency inversion, use cases and repositories are described through interfaces.
Concrete repository implementations realise communication with needed data sources, in this project it is postgresql.
Concrete use case implementations agregate repository interface; transport (grpc) agregates use case interface.
Such code organisation simplifies unit testing and allows us to make code flexible - we can easily add/switch between data sources and transports, write different use cases.

#### Request-response logic
Transport recieves call from client and calls use case. Use case validates the request and calls repository. Repository retrieves data from data source and forms entity. Entity is mapped to response structure and returned to client with status code. If any error appears, the app returns corresponding error with error code.

###DB schema
Database schema

## Launching and usage
Instructions on how to get, build and use

### Usage
GRPC clients, postman

### API methods
GetProduct ...