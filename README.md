### Create .env file from .env.example
```
cp .env.example .env
```
### Install all Go dependencies
```
go get -v ./...
```
### Generate and update swagger docs
```
go get -u github.com/swaggo/swag/cmd/swag
swag init
```
### [How to describe swagger routes](https://github.com/swaggo/swag/blob/master/README.md)
#### [Examples](https://github.com/swaggo/swag/blob/master/example/celler/controller/examples.go)
### Dockerization
```
docker-compose up -d
```
___
### How to 
##### Run migrations and seeds
```
go run db/migrate.go
```
##### Run tests 

```
go run db/migrate.go && go test ./... -count=1
```
use flag -count=1 to clear cache
##### Run linter
```
go fmt ./...
```
___
## Run project with live reload 
```
go get -u github.com/cosmtrek/air
type "air" in your command-line
``` 
## Without live reload
```
go run main.go
```
