# Authentication-Service
Go authentication service

# Help for me
### install migrate with sqlite3
    needed package: go get github.com/mattn/go-sqlite3
    go install -tags 'sqlite3' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
### create new migration:
    example: migrate create -seq -ext sql -dir ./migrations create_users_table
### apply migration
    migrate -database sqlite3://./database.db -path ./migrations up
