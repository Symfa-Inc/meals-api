### Create .env file from .env.example
```
cp .env.example .env
```
### Install all Go dependencies
```
go get ./...
```
### Generate and update swagger docs
```
swag init
```
### [How to describe swagger routes](https://github.com/swaggo/swag/blob/master/README.md)
#### [Examples](https://github.com/swaggo/swag/blob/master/example/celler/controller/examples.go)
___
### How to 
Run migrations
```
go run db/migrate.go
```
Run seeds
```
go run db/seeds.go
```
## Run project with live reload 
```
go get -u github.com/cosmtrek/air
type "air" in your command-line
``` 
## Without live reload
```
go run main.go
```
