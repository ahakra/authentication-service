# Authentication-Service
Go authentication service
# Project To-Do List

## To-Do List

### Completed Tasks:
- [x] Add `registerUserHandler` route (`r.Post("/registerUser", app.registerUserHandler)`)
- [x] Add `updateUserHandler` route (`r.Put("/registerUser", app.updateUserHandler)`)
- [x] Create database table for tokens 
- [x] Add insert token to database
  - [x] Implement `Insert` method in `TokenRepository`
  - [x] Handle token insertion logic in `TokenService`
- [x] Add Token service routes
  - [x] Define the route for token creation (`r.Post("/tokens/email", app.ReGenerateEmailTokenHandler)`)
  - [x] Define the route for token validation (`"/tokens/validate", app.ValidateTokenHandler`)
- [x] Add Permission insert
  - [x] Define route for permission insertion
  - [x] Implement logic to insert new permissions into the database
- [x] Get and manage permissions for users
  - [x] Create logic for fetching user-specific permissions
  - [x] Implement updating permissions for a user (Add/Remove)
- [x] Add Insert Token to database when creating new email

### ALL
- [ ] Modify all operation error
- [ ] Refactoring
- [ ] Send email for email verification (needs an email microservice)
### User Service:
- [ ] Add reset email
    - [ ] Implement reset password functionality
    - [ ] Send reset email with a unique token
    
- [ ] Add reset user password
    - [ ] Define route for resetting user password
    - [ ] Implement logic to verify token and reset password in the database

- [ ] Add regenerate password
    - [ ] Create an endpoint to regenerate the password (e.g., `/regeneratePassword`)

### Testing and Improvements:
- [ ] Write Unit tests
    - [ ] Write tests for token insertion
    - [ ] Write tests for token validation
    - [ ] Write tests for permission management
    - [ ] Write tests for user registration and password reset

### Cloud Native and DevOps:
- [ ] Add Grpc support
    - [ ] Implement gRPC services for token and user management
    - [ ] Define gRPC protocols (Protobufs)

- [ ] Add some cloud native patterns
    - [ ] Implement patterns such as circuit breakers or retries
    - [ ] Integrate with cloud services (e.g., S3 for file storage)

- [ ] Dockerize
    - [ ] Create Dockerfile for the application
    - [ ] Configure docker-compose for development and testing environments

### Miscellaneous:
- [ ] Add routes for token service
    - [ ] Define `/token` endpoint for managing user tokens
    - [ ] Handle token expiry and refresh logic

# Help for me
### install migrate with sqlite3
    needed package: go get github.com/mattn/go-sqlite3
    go install -tags 'sqlite3' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
### create new migration:
    example: migrate create -seq -ext sql -dir ./migrations create_users_table
### apply migration
    migrate -database sqlite3://./database.db -path ./migrations up
